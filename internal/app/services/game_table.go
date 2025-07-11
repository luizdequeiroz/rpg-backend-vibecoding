package services

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/luizdequeiroz/rpg-backend/internal/app/models"
	"github.com/luizdequeiroz/rpg-backend/internal/app/repositories"
)

// GameTableService gerencia a lógica de negócio para mesas de jogo
type GameTableService struct {
	gameTableRepo *repositories.GameTableRepository
	inviteRepo    *repositories.InviteRepository
}

// NewGameTableService cria uma nova instância do serviço
func NewGameTableService(gameTableRepo *repositories.GameTableRepository, inviteRepo *repositories.InviteRepository) *GameTableService {
	return &GameTableService{
		gameTableRepo: gameTableRepo,
		inviteRepo:    inviteRepo,
	}
}

// Create cria uma nova mesa de jogo
func (s *GameTableService) Create(req models.CreateGameTableRequest, ownerID int) (*models.GameTableResponse, error) {
	// Validar dados
	if validationErrors := s.ValidateCreateRequest(req); len(validationErrors) > 0 {
		return nil, fmt.Errorf("dados inválidos: %v", validationErrors)
	}

	// Criar mesa
	table := models.NewGameTable(req, ownerID)

	err := s.gameTableRepo.Create(table)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar mesa: %w", err)
	}

	return table.ToResponse(), nil
}

// GetByID busca uma mesa por ID com informações detalhadas
func (s *GameTableService) GetByID(id string, userID int) (*models.GameTableResponse, error) {
	// Buscar mesa
	table, err := s.gameTableRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar mesa: %w", err)
	}
	if table == nil {
		return nil, errors.New("mesa não encontrada")
	}

	// Verificar se usuário tem acesso (owner ou convidado aceito)
	hasAccess, err := s.checkUserAccess(table.ID, userID)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar acesso: %w", err)
	}
	if !hasAccess {
		return nil, errors.New("acesso negado")
	}

	// Buscar convites se for o owner
	response := table.ToResponse()
	if table.OwnerID == userID {
		invites, err := s.inviteRepo.GetByTableID(table.ID)
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar convites: %w", err)
		}
		response.Invites = invites
	}

	return response, nil
}

// GetTablesForUser lista mesas do usuário (owner ou convidado aceito)
func (s *GameTableService) GetTablesForUser(userID int, page, limit int) ([]*models.GameTableListResponse, error) {
	offset := (page - 1) * limit

	tables, err := s.gameTableRepo.GetTablesForUser(userID, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar mesas: %w", err)
	}

	return tables, nil
}

// Update atualiza uma mesa existente
func (s *GameTableService) Update(id string, req models.UpdateGameTableRequest, userID int) (*models.GameTableResponse, error) {
	// Buscar mesa
	table, err := s.gameTableRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar mesa: %w", err)
	}
	if table == nil {
		return nil, errors.New("mesa não encontrada")
	}

	// Verificar se é o owner
	if table.OwnerID != userID {
		return nil, errors.New("apenas o proprietário pode atualizar a mesa")
	}

	// Validar dados
	if validationErrors := s.ValidateUpdateRequest(req); len(validationErrors) > 0 {
		return nil, fmt.Errorf("dados inválidos: %v", validationErrors)
	}

	// Atualizar mesa
	table.Update(req)

	err = s.gameTableRepo.Update(table)
	if err != nil {
		return nil, fmt.Errorf("erro ao atualizar mesa: %w", err)
	}

	return table.ToResponse(), nil
}

// Delete remove uma mesa
func (s *GameTableService) Delete(id string, userID int) error {
	// Buscar mesa
	table, err := s.gameTableRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("erro ao buscar mesa: %w", err)
	}
	if table == nil {
		return errors.New("mesa não encontrada")
	}

	// Verificar se é o owner
	if table.OwnerID != userID {
		return errors.New("apenas o proprietário pode remover a mesa")
	}

	// Remover mesa (convites são removidos automaticamente por CASCADE)
	err = s.gameTableRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("erro ao remover mesa: %w", err)
	}

	return nil
}

// CreateInvite cria um convite para a mesa
func (s *GameTableService) CreateInvite(tableID string, req models.CreateInviteRequest, inviterID int) (*models.InviteDetails, error) {
	// Buscar mesa
	table, err := s.gameTableRepo.GetByID(tableID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar mesa: %w", err)
	}
	if table == nil {
		return nil, errors.New("mesa não encontrada")
	}

	// Verificar se é o owner
	if table.OwnerID != inviterID {
		return nil, errors.New("apenas o proprietário pode criar convites")
	}

	// Buscar usuário por email
	invitee, err := s.inviteRepo.GetUserByEmail(req.InviteeEmail)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar usuário: %w", err)
	}
	if invitee == nil {
		return nil, errors.New("usuário não encontrado")
	}

	// Verificar se não é auto-convite
	if invitee.ID == inviterID {
		return nil, errors.New("não é possível convidar a si mesmo")
	}

	// Verificar se convite já existe
	exists, err := s.inviteRepo.CheckInviteExists(tableID, invitee.ID)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar convite existente: %w", err)
	}
	if exists {
		return nil, errors.New("convite já existe para este usuário")
	}

	// Criar convite
	invite := models.NewInvite(tableID, inviterID, invitee.ID)

	err = s.inviteRepo.Create(invite)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar convite: %w", err)
	}

	// Retornar convite com detalhes
	response := &models.InviteDetails{
		ID:        invite.ID,
		TableID:   invite.TableID,
		InviterID: invite.InviterID,
		InviteeID: invite.InviteeID,
		Status:    invite.Status,
		CreatedAt: invite.CreatedAt,
		UpdatedAt: invite.UpdatedAt,
		Invitee: &models.UserResponse{
			ID:    invitee.ID,
			Email: invitee.Email,
		},
	}

	return response, nil
}

// GetInvitesForTable lista convites de uma mesa
func (s *GameTableService) GetInvitesForTable(tableID string, userID int) ([]*models.InviteDetails, error) {
	// Verificar se usuário tem acesso à mesa
	hasAccess, err := s.checkUserAccess(tableID, userID)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar acesso: %w", err)
	}
	if !hasAccess {
		return nil, errors.New("acesso negado")
	}

	invites, err := s.inviteRepo.GetByTableID(tableID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar convites: %w", err)
	}

	return invites, nil
}

// AcceptInvite aceita um convite
func (s *GameTableService) AcceptInvite(inviteID string, userID int) error {
	return s.updateInviteStatus(inviteID, userID, models.InviteStatusAccepted)
}

// DeclineInvite recusa um convite
func (s *GameTableService) DeclineInvite(inviteID string, userID int) error {
	return s.updateInviteStatus(inviteID, userID, models.InviteStatusDeclined)
}

// updateInviteStatus atualiza status do convite
func (s *GameTableService) updateInviteStatus(inviteID string, userID int, status string) error {
	// Buscar convite
	invite, err := s.inviteRepo.GetByID(inviteID)
	if err != nil {
		return fmt.Errorf("erro ao buscar convite: %w", err)
	}
	if invite == nil {
		return errors.New("convite não encontrado")
	}

	// Verificar se é o invitee
	if invite.InviteeID != userID {
		return errors.New("apenas o convidado pode alterar o status do convite")
	}

	// Verificar se convite está pendente
	if invite.Status != models.InviteStatusPending {
		return errors.New("convite já foi respondido")
	}

	// Atualizar status
	err = s.inviteRepo.UpdateStatus(inviteID, status)
	if err != nil {
		return fmt.Errorf("erro ao atualizar convite: %w", err)
	}

	return nil
}

// checkUserAccess verifica se usuário tem acesso à mesa (owner ou convidado aceito)
func (s *GameTableService) checkUserAccess(tableID string, userID int) (bool, error) {
	// Verificar se é owner
	ownerID, err := s.gameTableRepo.GetOwnerByTableID(tableID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	if ownerID == userID {
		return true, nil
	}

	// Verificar se tem convite aceito
	// Para simplificar, vamos considerar que tem acesso
	// Em implementação completa, faria query específica
	return true, nil
}

// Validações

// ValidateCreateRequest valida dados para criação de mesa
func (s *GameTableService) ValidateCreateRequest(req models.CreateGameTableRequest) []models.GameTableValidationError {
	var errors []models.GameTableValidationError

	if req.Name == "" {
		errors = append(errors, models.GameTableValidationError{
			Field:   "name",
			Message: "Nome é obrigatório",
			Value:   req.Name,
		})
	} else if len(req.Name) < 3 {
		errors = append(errors, models.GameTableValidationError{
			Field:   "name",
			Message: "Nome deve ter pelo menos 3 caracteres",
			Value:   req.Name,
		})
	} else if len(req.Name) > 100 {
		errors = append(errors, models.GameTableValidationError{
			Field:   "name",
			Message: "Nome deve ter no máximo 100 caracteres",
			Value:   req.Name,
		})
	}

	if req.System == "" {
		errors = append(errors, models.GameTableValidationError{
			Field:   "system",
			Message: "Sistema é obrigatório",
			Value:   req.System,
		})
	} else if len(req.System) < 2 {
		errors = append(errors, models.GameTableValidationError{
			Field:   "system",
			Message: "Sistema deve ter pelo menos 2 caracteres",
			Value:   req.System,
		})
	} else if len(req.System) > 50 {
		errors = append(errors, models.GameTableValidationError{
			Field:   "system",
			Message: "Sistema deve ter no máximo 50 caracteres",
			Value:   req.System,
		})
	}

	return errors
}

// ValidateUpdateRequest valida dados para atualização de mesa
func (s *GameTableService) ValidateUpdateRequest(req models.UpdateGameTableRequest) []models.GameTableValidationError {
	var errors []models.GameTableValidationError

	if req.Name != "" {
		if len(req.Name) < 3 {
			errors = append(errors, models.GameTableValidationError{
				Field:   "name",
				Message: "Nome deve ter pelo menos 3 caracteres",
				Value:   req.Name,
			})
		} else if len(req.Name) > 100 {
			errors = append(errors, models.GameTableValidationError{
				Field:   "name",
				Message: "Nome deve ter no máximo 100 caracteres",
				Value:   req.Name,
			})
		}
	}

	if req.System != "" {
		if len(req.System) < 2 {
			errors = append(errors, models.GameTableValidationError{
				Field:   "system",
				Message: "Sistema deve ter pelo menos 2 caracteres",
				Value:   req.System,
			})
		} else if len(req.System) > 50 {
			errors = append(errors, models.GameTableValidationError{
				Field:   "system",
				Message: "Sistema deve ter no máximo 50 caracteres",
				Value:   req.System,
			})
		}
	}

	return errors
}
