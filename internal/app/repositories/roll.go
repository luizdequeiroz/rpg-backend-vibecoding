package repositories

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/luizdequeiroz/rpg-backend/internal/app/models"
)

type RollRepository struct {
	db *sqlx.DB
}

func NewRollRepository(db *sqlx.DB) *RollRepository {
	return &RollRepository{db: db}
}

// Create cria uma nova rolagem
func (r *RollRepository) Create(roll *models.Roll) error {
	query := `
		INSERT INTO rolls (id, sheet_id, table_id, user_id, expression, field_name, 
		                  result_value, result_details, success, created_at)
		VALUES (:id, :sheet_id, :table_id, :user_id, :expression, :field_name, 
		        :result_value, :result_details, :success, :created_at)
	`

	// Preparar detalhes como JSON
	details := models.RollDetails{
		Total: roll.ResultValue,
	}
	detailsJSON, _ := json.Marshal(details)
	roll.ResultDetails = string(detailsJSON)
	roll.CreatedAt = time.Now()

	_, err := r.db.NamedExec(query, roll)
	return err
}

// GetByUserID recupera rolagens por usuário com paginação
func (r *RollRepository) GetByUserID(userID, limit, offset int) ([]models.Roll, error) {
	query := `
		SELECT id, sheet_id, table_id, user_id, expression, field_name, 
		       result_value, result_details, success, created_at
		FROM rolls 
		WHERE user_id = ? 
		ORDER BY created_at DESC 
		LIMIT ? OFFSET ?
	`

	var rolls []models.Roll
	err := r.db.Select(&rolls, query, userID, limit, offset)
	return rolls, err
}

// CountByUserID conta rolagens do usuário
func (r *RollRepository) CountByUserID(userID int) (int, error) {
	query := `SELECT COUNT(*) FROM rolls WHERE user_id = ?`

	var count int
	err := r.db.Get(&count, query, userID)
	return count, err
}

// UpdateSheetID atualiza o sheet_id de uma rolagem
func (r *RollRepository) UpdateSheetID(rollID, sheetID string) error {
	query := `UPDATE rolls SET sheet_id = ? WHERE id = ?`
	_, err := r.db.Exec(query, sheetID, rollID)
	return err
}

// GetBySheetID recupera rolagens por ficha
func (r *RollRepository) GetBySheetID(sheetID string) ([]models.Roll, error) {
	query := `
		SELECT id, sheet_id, table_id, user_id, expression, field_name, 
		       result_value, result_details, success, created_at
		FROM rolls 
		WHERE sheet_id = ? 
		ORDER BY created_at DESC
	`

	var rolls []models.Roll
	err := r.db.Select(&rolls, query, sheetID)
	return rolls, err
}

// GetByTableID lista rolagens de uma mesa
func (r *RollRepository) GetByTableID(tableID string, offset, limit int) ([]*models.RollResponse, error) {
	query := `
		SELECT 
			r.id, r.sheet_id, r.table_id, r.user_id, r.expression, r.field_name,
			r.result_value, r.result_details, r.success, r.created_at,
			u.id as "user.id", u.email as "user.email"
		FROM rolls r
		LEFT JOIN users u ON r.user_id = u.id
		WHERE r.table_id = ?
		ORDER BY r.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, tableID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rolls []*models.RollResponse

	for rows.Next() {
		var roll models.RollResponse
		var user models.UserResponse
		var detailsJSON sql.NullString

		err := rows.Scan(
			&roll.ID, &roll.SheetID, &roll.TableID, &roll.UserID, &roll.Expression, &roll.FieldName,
			&roll.ResultValue, &detailsJSON, &roll.Success, &roll.CreatedAt,
			&user.ID, &user.Email,
		)
		if err != nil {
			return nil, err
		}

		// Parse details JSON
		if detailsJSON.Valid {
			var details models.RollDetails
			json.Unmarshal([]byte(detailsJSON.String), &details)
			roll.ResultDetails = &details
		}

		roll.User = &user
		rolls = append(rolls, &roll)
	}

	return rolls, nil
}

// GetBySheetIDWithPagination lista rolagens de uma ficha com paginação
func (r *RollRepository) GetBySheetIDWithPagination(sheetID string, offset, limit int) ([]*models.RollResponse, error) {
	query := `
		SELECT 
			r.id, r.sheet_id, r.table_id, r.user_id, r.expression, r.field_name,
			r.result_value, r.result_details, r.success, r.created_at,
			u.id as "user.id", u.email as "user.email"
		FROM rolls r
		LEFT JOIN users u ON r.user_id = u.id
		WHERE r.sheet_id = ?
		ORDER BY r.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, sheetID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rolls []*models.RollResponse

	for rows.Next() {
		var roll models.RollResponse
		var user models.UserResponse
		var detailsJSON sql.NullString

		err := rows.Scan(
			&roll.ID, &roll.SheetID, &roll.TableID, &roll.UserID, &roll.Expression, &roll.FieldName,
			&roll.ResultValue, &detailsJSON, &roll.Success, &roll.CreatedAt,
			&user.ID, &user.Email,
		)
		if err != nil {
			return nil, err
		}

		// Parse details JSON
		if detailsJSON.Valid {
			var details models.RollDetails
			json.Unmarshal([]byte(detailsJSON.String), &details)
			roll.ResultDetails = &details
		}

		roll.User = &user
		rolls = append(rolls, &roll)
	}

	return rolls, nil
}
