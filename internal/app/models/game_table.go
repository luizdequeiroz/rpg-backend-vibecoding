package models

import (
	"time"

	"github.com/google/uuid"
)

// GameTable representa uma mesa de jogo de RPG
type GameTable struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name" validate:"required,min=3,max=100"`
	System    string    `json:"system" db:"system" validate:"required,min=2,max=50"`
	OwnerID   int       `json:"owner_id" db:"owner_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Invite representa um convite para uma mesa de jogo
type Invite struct {
	ID        string    `json:"id" db:"id"`
	TableID   string    `json:"table_id" db:"table_id"`
	InviterID int       `json:"inviter_id" db:"inviter_id"`
	InviteeID int       `json:"invitee_id" db:"invitee_id"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Request/Response models para GameTable

// CreateGameTableRequest para criação de mesa
type CreateGameTableRequest struct {
	Name   string `json:"name" validate:"required,min=3,max=100"`
	System string `json:"system" validate:"required,min=2,max=50"`
}

// UpdateGameTableRequest para atualização de mesa
type UpdateGameTableRequest struct {
	Name   string `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
	System string `json:"system,omitempty" validate:"omitempty,min=2,max=50"`
}

// GameTableResponse resposta detalhada da mesa com invites
type GameTableResponse struct {
	ID        string           `json:"id"`
	Name      string           `json:"name"`
	System    string           `json:"system"`
	OwnerID   int              `json:"owner_id"`
	Owner     *UserResponse    `json:"owner,omitempty"`
	Invites   []*InviteDetails `json:"invites,omitempty"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

// GameTableListResponse para listagem de mesas
type GameTableListResponse struct {
	ID        string        `json:"id"`
	Name      string        `json:"name"`
	System    string        `json:"system"`
	OwnerID   int           `json:"owner_id"`
	Owner     *UserResponse `json:"owner,omitempty"`
	Role      string        `json:"role"` // "owner" ou "player"
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// Request/Response models para Invite

// CreateInviteRequest para criação de convite
type CreateInviteRequest struct {
	InviteeEmail string `json:"invitee_email" validate:"required,email"`
}

// InviteDetails resposta detalhada do convite
type InviteDetails struct {
	ID        string        `json:"id"`
	TableID   string        `json:"table_id"`
	Table     *GameTable    `json:"table,omitempty"`
	InviterID int           `json:"inviter_id"`
	Inviter   *UserResponse `json:"inviter,omitempty"`
	InviteeID int           `json:"invitee_id"`
	Invitee   *UserResponse `json:"invitee,omitempty"`
	Status    string        `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// UserResponse resposta simplificada do usuário para relacionamentos
type UserResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

// Constantes para status de convite
const (
	InviteStatusPending  = "pending"
	InviteStatusAccepted = "accepted"
	InviteStatusDeclined = "declined"
)

// Constantes para roles de usuário na mesa
const (
	GameTableRoleOwner  = "owner"
	GameTableRolePlayer = "player"
)

// GameTableValidationError representa erros de validação
type GameTableValidationError struct {
	Field   string      `json:"field"`
	Message string      `json:"message"`
	Value   interface{} `json:"value"`
}

// Métodos utilitários

// NewGameTable cria uma nova mesa com UUID
func NewGameTable(req CreateGameTableRequest, ownerID int) *GameTable {
	now := time.Now()
	return &GameTable{
		ID:        uuid.New().String(),
		Name:      req.Name,
		System:    req.System,
		OwnerID:   ownerID,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewInvite cria um novo convite com UUID
func NewInvite(tableID string, inviterID, inviteeID int) *Invite {
	now := time.Now()
	return &Invite{
		ID:        uuid.New().String(),
		TableID:   tableID,
		InviterID: inviterID,
		InviteeID: inviteeID,
		Status:    InviteStatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// IsValidInviteStatus verifica se o status é válido
func IsValidInviteStatus(status string) bool {
	return status == InviteStatusPending ||
		status == InviteStatusAccepted ||
		status == InviteStatusDeclined
}

// ToResponse converte GameTable para GameTableResponse
func (gt *GameTable) ToResponse() *GameTableResponse {
	return &GameTableResponse{
		ID:        gt.ID,
		Name:      gt.Name,
		System:    gt.System,
		OwnerID:   gt.OwnerID,
		CreatedAt: gt.CreatedAt,
		UpdatedAt: gt.UpdatedAt,
	}
}

// ToListResponse converte GameTable para GameTableListResponse
func (gt *GameTable) ToListResponse(role string) *GameTableListResponse {
	return &GameTableListResponse{
		ID:        gt.ID,
		Name:      gt.Name,
		System:    gt.System,
		OwnerID:   gt.OwnerID,
		Role:      role,
		CreatedAt: gt.CreatedAt,
		UpdatedAt: gt.UpdatedAt,
	}
}

// Update atualiza os campos da mesa
func (gt *GameTable) Update(req UpdateGameTableRequest) {
	if req.Name != "" {
		gt.Name = req.Name
	}
	if req.System != "" {
		gt.System = req.System
	}
	gt.UpdatedAt = time.Now()
}
