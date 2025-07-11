package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/luizdequeiroz/rpg-backend/internal/app/models"
)

// PlayerSheetRepository gerencia acesso aos dados de fichas
type PlayerSheetRepository struct {
	db *sqlx.DB
}

// NewPlayerSheetRepository cria nova instância do repositório
func NewPlayerSheetRepository(db *sqlx.DB) *PlayerSheetRepository {
	return &PlayerSheetRepository{db: db}
}

// Create cria nova ficha no banco
func (r *PlayerSheetRepository) Create(sheet *models.PlayerSheet) error {
	query := `
		INSERT INTO player_sheets (id, table_id, template_id, owner_id, name, data, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query,
		sheet.ID, sheet.TableID, sheet.TemplateID, sheet.OwnerID,
		sheet.Name, sheet.Data, sheet.CreatedAt, sheet.UpdatedAt)

	return err
}

// GetByID busca ficha por ID
func (r *PlayerSheetRepository) GetByID(id string) (*models.PlayerSheet, error) {
	var sheet models.PlayerSheet

	query := `
		SELECT id, table_id, template_id, owner_id, name, data, created_at, updated_at
		FROM player_sheets 
		WHERE id = ?
	`

	err := r.db.Get(&sheet, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &sheet, err
}

// GetByTableID lista fichas por mesa
func (r *PlayerSheetRepository) GetByTableID(tableID string, offset, limit int) ([]*models.PlayerSheetListResponse, error) {
	// Primeiro buscar apenas as fichas
	querySheets := `
		SELECT id, table_id, template_id, owner_id, name, created_at, updated_at
		FROM player_sheets 
		WHERE table_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(querySheets, tableID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sheets []*models.PlayerSheetListResponse

	for rows.Next() {
		var sheet models.PlayerSheetListResponse

		err := rows.Scan(
			&sheet.ID, &sheet.TableID, &sheet.TemplateID, &sheet.OwnerID, &sheet.Name,
			&sheet.CreatedAt, &sheet.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Buscar dados do owner
		var owner models.UserResponse
		ownerQuery := `SELECT id, email FROM users WHERE id = ?`
		err = r.db.Get(&owner, ownerQuery, sheet.OwnerID)
		if err == nil {
			sheet.Owner = &owner
		} else {
			// Owner não encontrado, criar placeholder
			sheet.Owner = &models.UserResponse{
				ID:    sheet.OwnerID,
				Email: fmt.Sprintf("user-%d@deleted.com", sheet.OwnerID),
			}
		}

		// Buscar dados do template
		var template models.TemplateInfo
		templateQuery := `SELECT id, name, COALESCE(description, '') as description FROM sheet_templates WHERE id = ?`
		err = r.db.Get(&template, templateQuery, sheet.TemplateID)
		if err == nil {
			sheet.Template = &template
		}

		sheets = append(sheets, &sheet)
	}

	return sheets, nil
}

// GetByIDWithDetails busca ficha com detalhes completos
func (r *PlayerSheetRepository) GetByIDWithDetails(id string) (*models.PlayerSheetResponse, error) {
	// Buscar ficha primeiro
	var sheet models.PlayerSheet
	sheetQuery := `
		SELECT id, table_id, template_id, owner_id, name, data, created_at, updated_at
		FROM player_sheets 
		WHERE id = ?
	`

	err := r.db.Get(&sheet, sheetQuery, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Parse JSON data
	var data models.PlayerSheetData
	err = json.Unmarshal([]byte(sheet.Data), &data)
	if err != nil {
		return nil, fmt.Errorf("erro ao decodificar dados da ficha: %w", err)
	}

	// Buscar dados do owner
	var owner models.UserResponse
	ownerQuery := `SELECT id, email FROM users WHERE id = ?`
	err = r.db.Get(&owner, ownerQuery, sheet.OwnerID)
	if err != nil {
		// Se não encontrar o owner, criar um placeholder
		owner = models.UserResponse{
			ID:    sheet.OwnerID,
			Email: fmt.Sprintf("user-%d@deleted.com", sheet.OwnerID),
		}
	}

	// Buscar dados do template
	var template models.TemplateInfo
	templateQuery := `SELECT id, name, COALESCE(description, '') as description FROM sheet_templates WHERE id = ?`
	err = r.db.Get(&template, templateQuery, sheet.TemplateID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar template: %w", err)
	}

	// Construir response
	response := &models.PlayerSheetResponse{
		ID:         sheet.ID,
		TableID:    sheet.TableID,
		TemplateID: sheet.TemplateID,
		OwnerID:    sheet.OwnerID,
		Name:       sheet.Name,
		Data:       data,
		CreatedAt:  sheet.CreatedAt,
		UpdatedAt:  sheet.UpdatedAt,
		Owner:      &owner,
		Template:   &template,
	}

	return response, nil
}

// Update atualiza ficha existente
func (r *PlayerSheetRepository) Update(sheet *models.PlayerSheet) error {
	query := `
		UPDATE player_sheets 
		SET name = ?, data = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query, sheet.Name, sheet.Data, sheet.UpdatedAt, sheet.ID)
	return err
}

// Delete remove ficha
func (r *PlayerSheetRepository) Delete(id string) error {
	query := `DELETE FROM player_sheets WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

// CheckOwnership verifica se usuário é dono da ficha
func (r *PlayerSheetRepository) CheckOwnership(sheetID string, userID int) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM player_sheets WHERE id = ? AND owner_id = ?`
	err := r.db.Get(&count, query, sheetID, userID)
	return count > 0, err
}

// GetTableIDBySheetID retorna o ID da mesa de uma ficha
func (r *PlayerSheetRepository) GetTableIDBySheetID(sheetID string) (string, error) {
	var tableID string
	query := `SELECT table_id FROM player_sheets WHERE id = ?`
	err := r.db.Get(&tableID, query, sheetID)
	return tableID, err
}


