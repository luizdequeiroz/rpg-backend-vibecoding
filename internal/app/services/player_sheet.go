package services

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/luizdequeiroz/rpg-backend/internal/app/models"
	"github.com/luizdequeiroz/rpg-backend/internal/app/repositories"
	"github.com/luizdequeiroz/rpg-backend/pkg/roll"
)

// PlayerSheetService gerencia lógica de negócio para fichas
type PlayerSheetService struct {
	sheetRepo     *repositories.PlayerSheetRepository
	rollRepo      *repositories.RollRepository
	gameTableRepo *repositories.GameTableRepository
	rollEngine    *roll.RollEngine
}

// NewPlayerSheetService cria nova instância do serviço
func NewPlayerSheetService(
	sheetRepo *repositories.PlayerSheetRepository,
	rollRepo *repositories.RollRepository,
	gameTableRepo *repositories.GameTableRepository,
) *PlayerSheetService {
	return &PlayerSheetService{
		sheetRepo:     sheetRepo,
		rollRepo:      rollRepo,
		gameTableRepo: gameTableRepo,
		rollEngine:    roll.NewRollEngine(),
	}
}

// Create cria nova ficha
func (s *PlayerSheetService) Create(req models.CreatePlayerSheetRequest, tableID string, ownerID int) (*models.PlayerSheetResponse, error) {
	// Validar dados
	if validationErrors := s.ValidateCreateRequest(req); len(validationErrors) > 0 {
		return nil, fmt.Errorf("dados inválidos: %v", validationErrors)
	}

	// Verificar se usuário tem acesso à mesa
	hasAccess, err := s.checkTableAccess(tableID, ownerID)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar acesso à mesa: %w", err)
	}
	if !hasAccess {
		return nil, errors.New("acesso negado à mesa")
	}

	// Criar ficha
	sheet := models.NewPlayerSheet(req, tableID, ownerID)

	err = s.sheetRepo.Create(sheet)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar ficha: %w", err)
	}

	// Buscar ficha criada com detalhes
	response, err := s.sheetRepo.GetByIDWithDetails(sheet.ID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar ficha criada: %w", err)
	}

	return response, nil
}

// GetByID busca ficha por ID
func (s *PlayerSheetService) GetByID(id string, userID int) (*models.PlayerSheetResponse, error) {
	// Buscar ficha
	sheet, err := s.sheetRepo.GetByIDWithDetails(id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar ficha: %w", err)
	}
	if sheet == nil {
		return nil, errors.New("ficha não encontrada")
	}

	// Verificar acesso à mesa
	hasAccess, err := s.checkTableAccess(sheet.TableID, userID)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar acesso: %w", err)
	}
	if !hasAccess {
		return nil, errors.New("acesso negado")
	}

	return sheet, nil
}

// GetByTableID lista fichas da mesa
func (s *PlayerSheetService) GetByTableID(tableID string, userID int, page, limit int) ([]*models.PlayerSheetListResponse, error) {
	// Verificar acesso à mesa
	hasAccess, err := s.checkTableAccess(tableID, userID)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar acesso: %w", err)
	}
	if !hasAccess {
		return nil, errors.New("acesso negado à mesa")
	}

	offset := (page - 1) * limit
	sheets, err := s.sheetRepo.GetByTableID(tableID, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar fichas: %w", err)
	}

	return sheets, nil
}

// Update atualiza ficha
func (s *PlayerSheetService) Update(id string, req models.UpdatePlayerSheetRequest, userID int) (*models.PlayerSheetResponse, error) {
	// Buscar ficha
	sheetData, err := s.sheetRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar ficha: %w", err)
	}
	if sheetData == nil {
		return nil, errors.New("ficha não encontrada")
	}

	// Verificar se é o owner
	if sheetData.OwnerID != userID {
		return nil, errors.New("apenas o proprietário pode atualizar a ficha")
	}

	// Validar dados
	if validationErrors := s.ValidateUpdateRequest(req); len(validationErrors) > 0 {
		return nil, fmt.Errorf("dados inválidos: %v", validationErrors)
	}

	// Atualizar ficha
	sheetData.Update(req)

	err = s.sheetRepo.Update(sheetData)
	if err != nil {
		return nil, fmt.Errorf("erro ao atualizar ficha: %w", err)
	}

	// Retornar ficha atualizada
	response, err := s.sheetRepo.GetByIDWithDetails(id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar ficha atualizada: %w", err)
	}

	return response, nil
}

// Delete remove ficha
func (s *PlayerSheetService) Delete(id string, userID int) error {
	// Verificar ownership
	isOwner, err := s.sheetRepo.CheckOwnership(id, userID)
	if err != nil {
		return fmt.Errorf("erro ao verificar propriedade: %w", err)
	}

	if !isOwner {
		// Verificar se é owner da mesa
		tableID, err := s.sheetRepo.GetTableIDBySheetID(id)
		if err != nil {
			return fmt.Errorf("erro ao buscar mesa da ficha: %w", err)
		}

		tableOwnerID, err := s.gameTableRepo.GetOwnerByTableID(tableID)
		if err != nil {
			return fmt.Errorf("erro ao verificar proprietário da mesa: %w", err)
		}

		if tableOwnerID != userID {
			return errors.New("apenas o proprietário da ficha ou da mesa pode removê-la")
		}
	}

	err = s.sheetRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("erro ao remover ficha: %w", err)
	}

	return nil
}

// CreateRoll executa rolagem de dados
func (s *PlayerSheetService) CreateRoll(sheetID string, req models.CreateRollRequest, userID int) (*models.RollResponse, error) {
	// Buscar ficha
	sheet, err := s.sheetRepo.GetByIDWithDetails(sheetID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar ficha: %w", err)
	}
	if sheet == nil {
		return nil, errors.New("ficha não encontrada")
	}

	// Verificar acesso à mesa
	hasAccess, err := s.checkTableAccess(sheet.TableID, userID)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar acesso: %w", err)
	}
	if !hasAccess {
		return nil, errors.New("acesso negado à mesa")
	}

	// Validar request
	if req.Expression == "" && req.FieldName == "" {
		return nil, errors.New("expression ou field_name é obrigatório")
	}

	var rollDetails *models.RollDetails
	var expression string
	var fieldName *string

	// Executar rolagem
	if req.Expression != "" {
		// Rolagem direta por expressão
		rollDetails, err = s.rollEngine.Roll(req.Expression)
		expression = req.Expression
	} else {
		// Rolagem baseada em campo da ficha
		rollDetails, expression, err = s.rollEngine.RollFromField(sheet.Data, req.FieldName)
		fieldName = &req.FieldName
	}

	if err != nil {
		return nil, fmt.Errorf("erro na rolagem: %w", err)
	}

	// Criar record da rolagem
	rollRecord := models.NewRoll(sheetID, sheet.TableID, userID, expression, fieldName)
	rollRecord.ResultValue = rollDetails.Total

	// Serializar detalhes
	detailsJSON, _ := json.Marshal(rollDetails)
	rollRecord.ResultDetails = string(detailsJSON)

	// Salvar no banco
	err = s.rollRepo.Create(rollRecord)
	if err != nil {
		return nil, fmt.Errorf("erro ao salvar rolagem: %w", err)
	}

	// Retornar resposta
	response := rollRecord.ToResponse()
	response.ResultDetails = rollDetails
	response.User = &models.UserResponse{ID: userID}

	return response, nil
}

// GetRollsByTableID lista rolagens da mesa
func (s *PlayerSheetService) GetRollsByTableID(tableID string, userID int, page, limit int) ([]*models.RollResponse, error) {
	// Verificar acesso à mesa
	hasAccess, err := s.checkTableAccess(tableID, userID)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar acesso: %w", err)
	}
	if !hasAccess {
		return nil, errors.New("acesso negado à mesa")
	}

	offset := (page - 1) * limit
	rolls, err := s.rollRepo.GetByTableID(tableID, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar rolagens: %w", err)
	}

	return rolls, nil
}

// GetRollsBySheetID lista rolagens da ficha
func (s *PlayerSheetService) GetRollsBySheetID(sheetID string, userID int, page, limit int) ([]*models.RollResponse, error) {
	// Buscar mesa da ficha
	tableID, err := s.sheetRepo.GetTableIDBySheetID(sheetID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar mesa da ficha: %w", err)
	}

	// Verificar acesso à mesa
	hasAccess, err := s.checkTableAccess(tableID, userID)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar acesso: %w", err)
	}
	if !hasAccess {
		return nil, errors.New("acesso negado à mesa")
	}

	offset := (page - 1) * limit
	rolls, err := s.rollRepo.GetBySheetIDWithPagination(sheetID, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar rolagens: %w", err)
	}

	return rolls, nil
}

// checkTableAccess verifica se usuário tem acesso à mesa
func (s *PlayerSheetService) checkTableAccess(tableID string, userID int) (bool, error) {
	// Verificar se é owner da mesa
	ownerID, err := s.gameTableRepo.GetOwnerByTableID(tableID)
	if err != nil {
		return false, err
	}

	if ownerID == userID {
		return true, nil
	}

	// TODO: Verificar se tem convite aceito
	// Por enquanto, liberar acesso para simplificar
	return true, nil
}

// Validações

// ValidateCreateRequest valida dados para criação
func (s *PlayerSheetService) ValidateCreateRequest(req models.CreatePlayerSheetRequest) []models.PlayerSheetValidationError {
	var errors []models.PlayerSheetValidationError

	if req.TemplateID <= 0 {
		errors = append(errors, models.PlayerSheetValidationError{
			Field:   "template_id",
			Message: "Template ID é obrigatório",
			Value:   req.TemplateID,
		})
	}

	if req.Name == "" {
		errors = append(errors, models.PlayerSheetValidationError{
			Field:   "name",
			Message: "Nome é obrigatório",
			Value:   req.Name,
		})
	} else if len(req.Name) < 3 {
		errors = append(errors, models.PlayerSheetValidationError{
			Field:   "name",
			Message: "Nome deve ter pelo menos 3 caracteres",
			Value:   req.Name,
		})
	} else if len(req.Name) > 100 {
		errors = append(errors, models.PlayerSheetValidationError{
			Field:   "name",
			Message: "Nome deve ter no máximo 100 caracteres",
			Value:   req.Name,
		})
	}

	return errors
}

// ValidateUpdateRequest valida dados para atualização
func (s *PlayerSheetService) ValidateUpdateRequest(req models.UpdatePlayerSheetRequest) []models.PlayerSheetValidationError {
	var errors []models.PlayerSheetValidationError

	if req.Name != "" {
		if len(req.Name) < 3 {
			errors = append(errors, models.PlayerSheetValidationError{
				Field:   "name",
				Message: "Nome deve ter pelo menos 3 caracteres",
				Value:   req.Name,
			})
		} else if len(req.Name) > 100 {
			errors = append(errors, models.PlayerSheetValidationError{
				Field:   "name",
				Message: "Nome deve ter no máximo 100 caracteres",
				Value:   req.Name,
			})
		}
	}

	return errors
}

// getNestedField busca valor de campo aninhado no JSON de dados da ficha
// Suporta notação com ponto como "skills.arcana", "attributes.strength"
func (s *PlayerSheetService) getNestedField(data models.PlayerSheetData, fieldPath string) (interface{}, error) {
	// Converter data para map[string]interface{} para navegação dinâmica
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar dados da ficha: %w", err)
	}

	var dataMap map[string]interface{}
	err = json.Unmarshal(dataBytes, &dataMap)
	if err != nil {
		return nil, fmt.Errorf("erro ao deserializar dados da ficha: %w", err)
	}

	// Separar o caminho por pontos
	parts := []string{}
	currentPart := ""
	for _, char := range fieldPath {
		if char == '.' {
			if currentPart != "" {
				parts = append(parts, currentPart)
				currentPart = ""
			}
		} else {
			currentPart += string(char)
		}
	}
	if currentPart != "" {
		parts = append(parts, currentPart)
	}

	// Navegar pela estrutura
	current := dataMap
	for i, part := range parts {
		value, exists := current[part]
		if !exists {
			return nil, fmt.Errorf("campo '%s' não encontrado na ficha", fieldPath)
		}

		// Se for o último elemento, retornar o valor
		if i == len(parts)-1 {
			return value, nil
		}

		// Caso contrário, continuar navegando
		nextMap, ok := value.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("campo '%s' não é um objeto, não é possível navegar mais", part)
		}
		current = nextMap
	}

	return nil, fmt.Errorf("campo '%s' não encontrado", fieldPath)
}
