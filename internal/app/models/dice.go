package models

import (
	"time"
)

// DiceRollRequest representa uma solicitação de rolagem de dados
type DiceRollRequest struct {
	Expression string `json:"expression" binding:"required" example:"1d20+3"`
	Comment    string `json:"comment,omitempty" example:"Teste de Força"`
}

// DiceRollWithSheetRequest representa uma rolagem usando dados da ficha
type DiceRollWithSheetRequest struct {
	SheetID        string `json:"sheet_id" binding:"required" example:"1"`
	Expression     string `json:"expression" binding:"required" example:"1d20+{strength}"`
	AttributeField string `json:"attribute_field,omitempty" example:"strength"`
	Comment        string `json:"comment,omitempty" example:"Teste de Força com modificador da ficha"`
}

// DiceRollResponse representa o resultado de uma rolagem
type DiceRollResponse struct {
	ID         string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Expression string    `json:"expression" example:"1d20+3"`
	Result     int       `json:"result" example:"18"`
	Details    string    `json:"details" example:"[15] + 3 = 18"`
	IsCritical bool      `json:"is_critical" example:"false"`
	IsFumble   bool      `json:"is_fumble" example:"false"`
	SheetID    *string   `json:"sheet_id,omitempty" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID     int       `json:"user_id" example:"1"`
	CreatedAt  time.Time `json:"created_at" example:"2024-01-01T10:00:00Z"`
}

// DiceHistoryResponse representa o histórico de rolagens
type DiceHistoryResponse struct {
	Rolls      []DiceRollResponse `json:"rolls"`
	Total      int                `json:"total" example:"50"`
	Page       int                `json:"page" example:"1"`
	Limit      int                `json:"limit" example:"10"`
	TotalPages int                `json:"total_pages" example:"5"`
}
