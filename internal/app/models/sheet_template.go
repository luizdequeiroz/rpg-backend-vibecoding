package models

import (
	"encoding/json"
	"time"
)

// BaseTemplateFields contém os campos comuns para todos os tipos de template
type BaseTemplateFields struct {
	Name        string `json:"name" db:"name" validate:"required,min=3,max=100" example:"Ficha D&D 5e"`
	Description string `json:"description,omitempty" db:"description" validate:"max=500" example:"Template completo para personagens de D&D 5ª edição"`
}

// BaseTemplateFieldsOptional contém os campos comuns para atualizações (com ponteiros)
type BaseTemplateFieldsOptional struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=3,max=100" example:"Ficha D&D 5e Atualizada"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=500" example:"Template atualizado para D&D 5ª edição"`
}

// DatabaseFields contém os campos específicos do banco de dados
type DatabaseFields struct {
	ID        int       `json:"id" db:"id" example:"1"`
	IsActive  bool      `json:"is_active" db:"is_active" example:"true"`
	CreatedAt time.Time `json:"created_at" db:"created_at" example:"2025-07-11T09:26:29Z"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" example:"2025-07-11T09:26:29Z"`
}

// DefinitionField contém o campo de definição para diferentes contextos
type DefinitionField struct {
	Definition interface{} `json:"definition" validate:"required" swaggertype:"object"`
}

// DefinitionFieldOptional contém o campo de definição opcional para atualizações
type DefinitionFieldOptional struct {
	Definition interface{} `json:"definition,omitempty" swaggertype:"object"`
}

// DefinitionFieldDB contém o campo de definição para o banco de dados
type DefinitionFieldDB struct {
	Definition string `json:"-" db:"definition" validate:"required"` // JSON armazenado como string no banco (não exposto na API)
}

// SheetTemplate representa um template de ficha de personagem armazenado no banco de dados
type SheetTemplate struct {
	DatabaseFields
	BaseTemplateFields
	DefinitionFieldDB
}

// CreateSheetTemplateRequest representa os dados necessários para criar um novo template de ficha
type CreateSheetTemplateRequest struct {
	BaseTemplateFields
	DefinitionField
}

// UpdateSheetTemplateRequest representa os dados para atualizar um template existente
type UpdateSheetTemplateRequest struct {
	BaseTemplateFieldsOptional
	DefinitionFieldOptional
}

// SheetTemplateResponse representa a resposta completa de um template de ficha para a API
type SheetTemplateResponse struct {
	DatabaseFields
	BaseTemplateFields
	DefinitionField
}

// SheetTemplateListResponse representa a resposta paginada da listagem de templates
type SheetTemplateListResponse struct {
	Templates []SheetTemplateResponse `json:"templates"`
	Total     int                     `json:"total" example:"25"`
	Page      int                     `json:"page,omitempty" example:"1"`
	Limit     int                     `json:"limit,omitempty" example:"20"`
	HasNext   bool                    `json:"has_next,omitempty" example:"true"`
}

// ToResponse converte um SheetTemplate do banco de dados para SheetTemplateResponse da API
func (st *SheetTemplate) ToResponse() SheetTemplateResponse {
	var definition interface{}
	// Converte a string JSON do banco para interface{} para a API
	if st.Definition != "" {
		json.Unmarshal([]byte(st.Definition), &definition)
	}

	return SheetTemplateResponse{
		DatabaseFields: DatabaseFields{
			ID:        st.ID,
			IsActive:  st.IsActive,
			CreatedAt: st.CreatedAt,
			UpdatedAt: st.UpdatedAt,
		},
		BaseTemplateFields: BaseTemplateFields{
			Name:        st.Name,
			Description: st.Description,
		},
		DefinitionField: DefinitionField{
			Definition: definition,
		},
	}
}

// BaseErrorResponse contém campos comuns para respostas de erro
type BaseErrorResponse struct {
	Error     string    `json:"error" example:"Erro no processamento"`
	Message   string    `json:"message,omitempty" example:"Detalhes adicionais do erro"`
	Timestamp time.Time `json:"timestamp" example:"2025-07-11T09:26:29Z"`
}

// SheetTemplateValidationError representa um erro específico de validação de template
type SheetTemplateValidationError struct {
	Field   string `json:"field" example:"name"`
	Message string `json:"message" example:"Nome deve ter pelo menos 3 caracteres"`
	Value   string `json:"value,omitempty" example:"ab"`
}

// SheetTemplateErrorResponse representa uma resposta de erro específica para templates
type SheetTemplateErrorResponse struct {
	BaseErrorResponse
	Fields []SheetTemplateValidationError `json:"fields,omitempty"`
}

// SheetTemplateNotFoundError representa erro quando template não é encontrado
type SheetTemplateNotFoundError struct {
	BaseErrorResponse
	ID int `json:"id,omitempty" example:"999"`
}

// ConvertDefinitionToString converte a interface{} definition para string JSON
func ConvertDefinitionToString(definition interface{}) (string, error) {
	if definition == nil {
		return "", nil
	}

	jsonBytes, err := json.Marshal(definition)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

// IsValidDefinition verifica se a definition tem estrutura válida básica
func IsValidDefinition(definition interface{}) bool {
	if definition == nil {
		return false
	}

	// Converte para map para verificar estrutura básica
	definitionMap, ok := definition.(map[string]interface{})
	if !ok {
		return false
	}

	// Verifica se tem pelo menos a chave "sections"
	sections, exists := definitionMap["sections"]
	if !exists {
		return false
	}

	// Verifica se sections é um array
	_, ok = sections.([]interface{})
	return ok
}

// NewSheetTemplateErrorResponse cria uma resposta de erro padronizada para templates
func NewSheetTemplateErrorResponse(errorMsg string, message string) SheetTemplateErrorResponse {
	return SheetTemplateErrorResponse{
		BaseErrorResponse: BaseErrorResponse{
			Error:     errorMsg,
			Message:   message,
			Timestamp: time.Now(),
		},
	}
}

// NewSheetTemplateNotFoundError cria um erro padronizado para template não encontrado
func NewSheetTemplateNotFoundError(id int) SheetTemplateNotFoundError {
	return SheetTemplateNotFoundError{
		BaseErrorResponse: BaseErrorResponse{
			Error:     "Template não encontrado",
			Message:   "O template solicitado não existe ou foi removido",
			Timestamp: time.Now(),
		},
		ID: id,
	}
}
