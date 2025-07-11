package services

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/luizdequeiroz/rpg-backend/internal/app/models"
	"github.com/luizdequeiroz/rpg-backend/internal/app/repositories"
	"github.com/luizdequeiroz/rpg-backend/pkg/db"
)

var (
	ErrTemplateNotFound  = errors.New("template não encontrado")
	ErrInvalidDefinition = errors.New("definição JSON inválida")
	ErrInvalidPagination = errors.New("parâmetros de paginação inválidos")
)

// SheetTemplateService gerencia a lógica de negócio para templates
type SheetTemplateService struct {
	repo *repositories.SheetTemplateRepository
}

// NewSheetTemplateService cria uma nova instância do serviço
func NewSheetTemplateService(database *db.DB) *SheetTemplateService {
	return &SheetTemplateService{
		repo: repositories.NewSheetTemplateRepository(database),
	}
}

// GetAll retorna todos os templates com paginação
func (s *SheetTemplateService) GetAll(pageStr, limitStr string) (*models.SheetTemplateListResponse, error) {
	// Valores padrão
	page := 1
	limit := 20

	// Parse dos parâmetros
	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil || p < 1 {
			return nil, ErrInvalidPagination
		}
		page = p
	}

	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil || l < 1 || l > 100 {
			return nil, ErrInvalidPagination
		}
		limit = l
	}

	templates, total, err := s.repo.GetAll(page, limit)
	if err != nil {
		return nil, err
	}

	// Converter para response
	responses := make([]models.SheetTemplateResponse, len(templates))
	for i, template := range templates {
		responses[i] = template.ToResponse()
	}

	return &models.SheetTemplateListResponse{
		Templates: responses,
		Total:     total,
		Page:      page,
		Limit:     limit,
	}, nil
}

// GetByID retorna um template por ID
func (s *SheetTemplateService) GetByID(id int) (*models.SheetTemplateResponse, error) {
	template, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if template == nil {
		return nil, ErrTemplateNotFound
	}

	response := template.ToResponse()
	return &response, nil
}

// Create cria um novo template
func (s *SheetTemplateService) Create(req models.CreateSheetTemplateRequest) (*models.SheetTemplateResponse, error) {
	template, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	response := template.ToResponse()
	return &response, nil
}

// Update atualiza um template existente
func (s *SheetTemplateService) Update(id int, req models.UpdateSheetTemplateRequest) (*models.SheetTemplateResponse, error) {
	template, err := s.repo.Update(id, req)
	if err != nil {
		return nil, err
	}
	if template == nil {
		return nil, ErrTemplateNotFound
	}

	response := template.ToResponse()
	return &response, nil
}

// Delete remove um template
func (s *SheetTemplateService) Delete(id int) error {
	// Verificar se existe
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrTemplateNotFound
	}

	return s.repo.Delete(id)
}

// GetBySystem retorna templates por sistema
func (s *SheetTemplateService) GetBySystem(system string) ([]models.SheetTemplateResponse, error) {
	templates, err := s.repo.GetBySystem(system)
	if err != nil {
		return nil, err
	}

	responses := make([]models.SheetTemplateResponse, len(templates))
	for i, template := range templates {
		responses[i] = template.ToResponse()
	}

	return responses, nil
}

// validateDefinition valida se a definição JSON é válida
func (s *SheetTemplateService) validateDefinition(definition json.RawMessage) error {
	var temp interface{}
	if err := json.Unmarshal(definition, &temp); err != nil {
		return ErrInvalidDefinition
	}
	return nil
}

// ValidateCreateRequest valida os dados para criação de template
func (s *SheetTemplateService) ValidateCreateRequest(req models.CreateSheetTemplateRequest) []models.SheetTemplateValidationError {
	var errors []models.SheetTemplateValidationError

	if req.Name == "" {
		errors = append(errors, models.SheetTemplateValidationError{
			Field:   "name",
			Message: "Nome é obrigatório",
		})
	} else if len(req.Name) < 3 {
		errors = append(errors, models.SheetTemplateValidationError{
			Field:   "name",
			Message: "Nome deve ter pelo menos 3 caracteres",
			Value:   req.Name,
		})
	} else if len(req.Name) > 100 {
		errors = append(errors, models.SheetTemplateValidationError{
			Field:   "name",
			Message: "Nome deve ter no máximo 100 caracteres",
			Value:   req.Name,
		})
	}

	if req.System == "" {
		errors = append(errors, models.SheetTemplateValidationError{
			Field:   "system",
			Message: "Sistema é obrigatório",
		})
	} else if len(req.System) < 2 {
		errors = append(errors, models.SheetTemplateValidationError{
			Field:   "system",
			Message: "Sistema deve ter pelo menos 2 caracteres",
			Value:   req.System,
		})
	} else if len(req.System) > 50 {
		errors = append(errors, models.SheetTemplateValidationError{
			Field:   "system",
			Message: "Sistema deve ter no máximo 50 caracteres",
			Value:   req.System,
		})
	}

	if req.Definition == nil {
		errors = append(errors, models.SheetTemplateValidationError{
			Field:   "definition",
			Message: "Definition é obrigatória",
		})
	} else if !models.IsValidDefinition(req.Definition) {
		errors = append(errors, models.SheetTemplateValidationError{
			Field:   "definition",
			Message: "Definition deve ter estrutura válida com sections",
		})
	}

	if len(req.Description) > 500 {
		errors = append(errors, models.SheetTemplateValidationError{
			Field:   "description",
			Message: "Descrição deve ter no máximo 500 caracteres",
			Value:   req.Description,
		})
	}

	return errors
}

// ValidateUpdateRequest valida os dados para atualização de template
func (s *SheetTemplateService) ValidateUpdateRequest(req models.UpdateSheetTemplateRequest) []models.SheetTemplateValidationError {
	var errors []models.SheetTemplateValidationError

	if req.Name != nil {
		if *req.Name == "" {
			errors = append(errors, models.SheetTemplateValidationError{
				Field:   "name",
				Message: "Nome não pode ser vazio",
			})
		} else if len(*req.Name) < 3 {
			errors = append(errors, models.SheetTemplateValidationError{
				Field:   "name",
				Message: "Nome deve ter pelo menos 3 caracteres",
				Value:   *req.Name,
			})
		} else if len(*req.Name) > 100 {
			errors = append(errors, models.SheetTemplateValidationError{
				Field:   "name",
				Message: "Nome deve ter no máximo 100 caracteres",
				Value:   *req.Name,
			})
		}
	}

	if req.System != nil {
		if *req.System == "" {
			errors = append(errors, models.SheetTemplateValidationError{
				Field:   "system",
				Message: "Sistema não pode ser vazio",
			})
		} else if len(*req.System) < 2 {
			errors = append(errors, models.SheetTemplateValidationError{
				Field:   "system",
				Message: "Sistema deve ter pelo menos 2 caracteres",
				Value:   *req.System,
			})
		} else if len(*req.System) > 50 {
			errors = append(errors, models.SheetTemplateValidationError{
				Field:   "system",
				Message: "Sistema deve ter no máximo 50 caracteres",
				Value:   *req.System,
			})
		}
	}

	if req.Definition != nil && !models.IsValidDefinition(req.Definition) {
		errors = append(errors, models.SheetTemplateValidationError{
			Field:   "definition",
			Message: "Definition deve ter estrutura válida com sections",
		})
	}

	if req.Description != nil && len(*req.Description) > 500 {
		errors = append(errors, models.SheetTemplateValidationError{
			Field:   "description",
			Message: "Descrição deve ter no máximo 500 caracteres",
			Value:   *req.Description,
		})
	}

	return errors
}
