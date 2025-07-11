package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/luizdequeiroz/rpg-backend/internal/app/models"
	"github.com/luizdequeiroz/rpg-backend/pkg/db"
)

// SheetTemplateRepository gerencia operações de banco para templates de ficha
type SheetTemplateRepository struct {
	db *db.DB
}

// NewSheetTemplateRepository cria uma nova instância do repositório
func NewSheetTemplateRepository(database *db.DB) *SheetTemplateRepository {
	return &SheetTemplateRepository{
		db: database,
	}
}

// GetAll busca todos os templates ativos
func (r *SheetTemplateRepository) GetAll(page, limit int) ([]models.SheetTemplate, int, error) {
	var templates []models.SheetTemplate
	var total int

	// Contar total de registros
	err := r.db.Get(&total, "SELECT COUNT(*) FROM sheet_templates WHERE is_active = true")
	if err != nil {
		return nil, 0, fmt.Errorf("erro ao contar templates: %w", err)
	}

	// Buscar templates com paginação
	offset := (page - 1) * limit
	query := `
		SELECT id, name, definition, description, is_active, created_at, updated_at 
		FROM sheet_templates 
		WHERE is_active = true 
		ORDER BY created_at DESC 
		LIMIT ? OFFSET ?
	`

	err = r.db.Select(&templates, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("erro ao buscar templates: %w", err)
	}

	return templates, total, nil
}

// GetByID busca um template por ID
func (r *SheetTemplateRepository) GetByID(id int) (*models.SheetTemplate, error) {
	var template models.SheetTemplate
	query := `
		SELECT id, name, definition, description, is_active, created_at, updated_at 
		FROM sheet_templates 
		WHERE id = ? AND is_active = true
	`

	err := r.db.Get(&template, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Não encontrado
		}
		return nil, fmt.Errorf("erro ao buscar template: %w", err)
	}

	return &template, nil
}

// Create cria um novo template
func (r *SheetTemplateRepository) Create(req models.CreateSheetTemplateRequest) (*models.SheetTemplate, error) {
	now := time.Now()
	query := `
		INSERT INTO sheet_templates (name, definition, description, is_active, created_at, updated_at)
		VALUES (?, ?, ?, true, ?, ?)
	`

	// Converter definition para string JSON
	definitionBytes, err := json.Marshal(req.Definition)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter definition para JSON: %w", err)
	}

	result, err := r.db.Exec(query, req.Name, string(definitionBytes), req.Description, now, now)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar template: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter ID do template criado: %w", err)
	}

	return r.GetByID(int(id))
}

// Update atualiza um template existente
func (r *SheetTemplateRepository) Update(id int, req models.UpdateSheetTemplateRequest) (*models.SheetTemplate, error) {
	// Verificar se o template existe
	existing, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, nil // Não encontrado
	}

	now := time.Now()

	// Preparar valores para atualização, usando valores existentes se não fornecidos
	name := existing.Name
	if req.Name != nil {
		name = *req.Name
	}

	description := existing.Description
	if req.Description != nil {
		description = *req.Description
	}

	definitionJSON := existing.Definition
	if req.Definition != nil {
		// Converter definition para string JSON
		definitionBytes, err := json.Marshal(req.Definition)
		if err != nil {
			return nil, fmt.Errorf("erro ao converter definition para JSON: %w", err)
		}
		definitionJSON = string(definitionBytes)
	}

	query := `
		UPDATE sheet_templates 
		SET name = ?, definition = ?, description = ?, updated_at = ?
		WHERE id = ? AND is_active = true
	`

	_, err = r.db.Exec(query, name, definitionJSON, description, now, id)
	if err != nil {
		return nil, fmt.Errorf("erro ao atualizar template: %w", err)
	}

	return r.GetByID(id)
}

// Delete remove um template (soft delete)
func (r *SheetTemplateRepository) Delete(id int) error {
	// Verificar se o template existe
	existing, err := r.GetByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return nil // Não encontrado, mas não é erro
	}

	query := "UPDATE sheet_templates SET is_active = false, updated_at = ? WHERE id = ?"
	_, err = r.db.Exec(query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("erro ao deletar template: %w", err)
	}

	return nil
}
