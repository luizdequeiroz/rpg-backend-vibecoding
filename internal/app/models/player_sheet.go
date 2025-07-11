package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// PlayerSheet representa uma ficha de personagem em uma mesa
type PlayerSheet struct {
	ID         string    `json:"id" db:"id"`
	TableID    string    `json:"table_id" db:"table_id"`
	TemplateID int       `json:"template_id" db:"template_id"`
	OwnerID    int       `json:"owner_id" db:"owner_id"`
	Name       string    `json:"name" db:"name"`
	Data       string    `json:"-" db:"data"` // JSON como string no banco
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// PlayerSheetData representa os dados JSON da ficha
type PlayerSheetData map[string]interface{}

// PlayerSheetResponse representa a resposta da API
type PlayerSheetResponse struct {
	ID         string          `json:"id"`
	TableID    string          `json:"table_id"`
	TemplateID int             `json:"template_id"`
	OwnerID    int             `json:"owner_id"`
	Name       string          `json:"name"`
	Data       PlayerSheetData `json:"data"`
	Owner      *UserResponse   `json:"owner,omitempty"`
	Template   *TemplateInfo   `json:"template,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

// PlayerSheetListResponse representa item na listagem
type PlayerSheetListResponse struct {
	ID         string        `json:"id"`
	TableID    string        `json:"table_id"`
	TemplateID int           `json:"template_id"`
	OwnerID    int           `json:"owner_id"`
	Name       string        `json:"name"`
	Owner      *UserResponse `json:"owner,omitempty"`
	Template   *TemplateInfo `json:"template,omitempty"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
}

// TemplateInfo representa informações básicas do template
type TemplateInfo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CreatePlayerSheetRequest representa dados para criação
type CreatePlayerSheetRequest struct {
	TableID    string          `json:"table_id" validate:"required,uuid"`
	TemplateID int             `json:"template_id" validate:"required,min=1"`
	Name       string          `json:"name" validate:"required,min=3,max=100"`
	Data       PlayerSheetData `json:"data"`
}

// UpdatePlayerSheetRequest representa dados para atualização
type UpdatePlayerSheetRequest struct {
	Name string          `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
	Data PlayerSheetData `json:"data,omitempty"`
}

// Roll representa uma rolagem de dados
type Roll struct {
	ID            string    `json:"id" db:"id"`
	SheetID       *string   `json:"sheet_id" db:"sheet_id"`
	TableID       *string   `json:"table_id" db:"table_id"`
	UserID        int       `json:"user_id" db:"user_id"`
	Expression    string    `json:"expression" db:"expression"`
	FieldName     *string   `json:"field_name" db:"field_name"`
	ResultValue   int       `json:"result_value" db:"result_value"`
	ResultDetails string    `json:"-" db:"result_details"` // JSON como string
	Success       *bool     `json:"success" db:"success"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

// RollDetails representa detalhes da rolagem
type RollDetails struct {
	Dice     []int `json:"dice"`     // Valores individuais dos dados
	Modifier int   `json:"modifier"` // Modificador aplicado
	Total    int   `json:"total"`    // Resultado final
	Critical bool  `json:"critical"` // Se foi crítico
	Fumble   bool  `json:"fumble"`   // Se foi fumble
}

// RollResponse representa resposta da rolagem
type RollResponse struct {
	ID            string        `json:"id"`
	SheetID       *string       `json:"sheet_id"`
	TableID       *string       `json:"table_id"`
	UserID        int           `json:"user_id"`
	Expression    string        `json:"expression"`
	FieldName     *string       `json:"field_name"`
	ResultValue   int           `json:"result_value"`
	ResultDetails *RollDetails  `json:"result_details"`
	Success       *bool         `json:"success"`
	User          *UserResponse `json:"user,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
}

// CreateRollRequest representa dados para rolagem
type CreateRollRequest struct {
	SheetID    string `json:"sheet_id" validate:"required,uuid"`
	Expression string `json:"expression,omitempty" validate:"omitempty,max=200"`
	FieldName  string `json:"field_name,omitempty" validate:"omitempty,max=100"`
}

// PlayerSheetValidationError representa erro de validação
type PlayerSheetValidationError struct {
	Field   string      `json:"field"`
	Message string      `json:"message"`
	Value   interface{} `json:"value"`
}

// NewPlayerSheet cria nova ficha
func NewPlayerSheet(req CreatePlayerSheetRequest, tableID string, ownerID int) *PlayerSheet {
	dataJSON, _ := json.Marshal(req.Data)
	if req.Data == nil {
		dataJSON = []byte("{}")
	}

	return &PlayerSheet{
		ID:         uuid.New().String(),
		TableID:    tableID,
		TemplateID: req.TemplateID,
		OwnerID:    ownerID,
		Name:       req.Name,
		Data:       string(dataJSON),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

// NewRoll cria nova rolagem
func NewRoll(sheetID, tableID string, userID int, expression string, fieldName *string) *Roll {
	// Usar sheetID como string
	var sheetStrID *string
	if sheetID != "" {
		sheetStrID = &sheetID
	}

	// Converter tableID para string pointer
	var tableStrID *string
	if tableID != "" {
		tableStrID = &tableID
	}

	return &Roll{
		ID:         uuid.New().String(),
		SheetID:    sheetStrID,
		TableID:    tableStrID,
		UserID:     userID,
		Expression: expression,
		FieldName:  fieldName,
		CreatedAt:  time.Now(),
	}
}

// ToResponse converte PlayerSheet para resposta
func (ps *PlayerSheet) ToResponse() *PlayerSheetResponse {
	var data PlayerSheetData
	json.Unmarshal([]byte(ps.Data), &data)
	if data == nil {
		data = make(PlayerSheetData)
	}

	return &PlayerSheetResponse{
		ID:         ps.ID,
		TableID:    ps.TableID,
		TemplateID: ps.TemplateID,
		OwnerID:    ps.OwnerID,
		Name:       ps.Name,
		Data:       data,
		CreatedAt:  ps.CreatedAt,
		UpdatedAt:  ps.UpdatedAt,
	}
}

// ToListResponse converte para item de listagem
func (ps *PlayerSheet) ToListResponse() *PlayerSheetListResponse {
	return &PlayerSheetListResponse{
		ID:         ps.ID,
		TableID:    ps.TableID,
		TemplateID: ps.TemplateID,
		OwnerID:    ps.OwnerID,
		Name:       ps.Name,
		CreatedAt:  ps.CreatedAt,
		UpdatedAt:  ps.UpdatedAt,
	}
}

// Update atualiza dados da ficha
func (ps *PlayerSheet) Update(req UpdatePlayerSheetRequest) {
	if req.Name != "" {
		ps.Name = req.Name
	}

	if req.Data != nil {
		dataJSON, _ := json.Marshal(req.Data)
		ps.Data = string(dataJSON)
	}

	ps.UpdatedAt = time.Now()
}

// ToResponse converte Roll para resposta
func (r *Roll) ToResponse() *RollResponse {
	var details *RollDetails
	if r.ResultDetails != "" {
		json.Unmarshal([]byte(r.ResultDetails), &details)
	}

	return &RollResponse{
		ID:            r.ID,
		SheetID:       r.SheetID,
		TableID:       r.TableID,
		UserID:        r.UserID,
		Expression:    r.Expression,
		FieldName:     r.FieldName,
		ResultValue:   r.ResultValue,
		ResultDetails: details,
		Success:       r.Success,
		CreatedAt:     r.CreatedAt,
	}
}
