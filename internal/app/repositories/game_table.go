package repositories

import (
	"database/sql"

	"github.com/luizdequeiroz/rpg-backend/internal/app/models"
	"github.com/luizdequeiroz/rpg-backend/pkg/db"
)

// GameTableRepository gerencia operações de dados para mesas de jogo
type GameTableRepository struct {
	db *db.DB
}

// NewGameTableRepository cria uma nova instância do repositório
func NewGameTableRepository(database *db.DB) *GameTableRepository {
	return &GameTableRepository{
		db: database,
	}
}

// Create cria uma nova mesa de jogo
func (r *GameTableRepository) Create(table *models.GameTable) error {
	query := `
		INSERT INTO game_tables (id, name, system, owner_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query,
		table.ID, table.Name, table.System, table.OwnerID,
		table.CreatedAt, table.UpdatedAt,
	)

	return err
}

// GetByID busca uma mesa por ID
func (r *GameTableRepository) GetByID(id string) (*models.GameTable, error) {
	var table models.GameTable

	query := `
		SELECT id, name, system, owner_id, created_at, updated_at
		FROM game_tables 
		WHERE id = ?
	`

	err := r.db.Get(&table, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &table, nil
}

// GetByOwnerID busca mesas por proprietário
func (r *GameTableRepository) GetByOwnerID(ownerID int, offset, limit int) ([]*models.GameTable, error) {
	var tables []*models.GameTable

	query := `
		SELECT id, name, system, owner_id, created_at, updated_at
		FROM game_tables 
		WHERE owner_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	err := r.db.Select(&tables, query, ownerID, limit, offset)
	return tables, err
}

// GetTablesForUser busca mesas onde o usuário é owner ou convidado aceito
func (r *GameTableRepository) GetTablesForUser(userID int, offset, limit int) ([]*models.GameTableListResponse, error) {
	var results []*models.GameTableListResponse

	query := `
		SELECT DISTINCT 
			gt.id, gt.name, gt.system, gt.owner_id, gt.created_at, gt.updated_at,
			u.id as owner_user_id, u.email as owner_email,
			CASE 
				WHEN gt.owner_id = ? THEN 'owner'
				ELSE 'player'
			END as role
		FROM game_tables gt
		JOIN users u ON gt.owner_id = u.id
		LEFT JOIN invites i ON gt.id = i.table_id
		WHERE gt.owner_id = ? 
		   OR (i.invitee_id = ? AND i.status = 'accepted')
		ORDER BY gt.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, userID, userID, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var result models.GameTableListResponse
		var owner models.UserResponse

		err := rows.Scan(
			&result.ID, &result.Name, &result.System, &result.OwnerID,
			&result.CreatedAt, &result.UpdatedAt,
			&owner.ID, &owner.Email,
			&result.Role,
		)
		if err != nil {
			return nil, err
		}

		result.Owner = &owner
		results = append(results, &result)
	}

	return results, rows.Err()
}

// Update atualiza uma mesa existente
func (r *GameTableRepository) Update(table *models.GameTable) error {
	query := `
		UPDATE game_tables 
		SET name = ?, system = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.Exec(query, table.Name, table.System, table.UpdatedAt, table.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Delete remove uma mesa
func (r *GameTableRepository) Delete(id string) error {
	query := `DELETE FROM game_tables WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetOwnerByTableID busca o ID do proprietário de uma mesa
func (r *GameTableRepository) GetOwnerByTableID(tableID string) (int, error) {
	var ownerID int

	query := `SELECT owner_id FROM game_tables WHERE id = ?`

	err := r.db.Get(&ownerID, query, tableID)
	if err != nil {
		return 0, err
	}

	return ownerID, nil
}

// InviteRepository gerencia operações de dados para convites
type InviteRepository struct {
	db *db.DB
}

// NewInviteRepository cria uma nova instância do repositório de convites
func NewInviteRepository(database *db.DB) *InviteRepository {
	return &InviteRepository{
		db: database,
	}
}

// Create cria um novo convite
func (r *InviteRepository) Create(invite *models.Invite) error {
	query := `
		INSERT INTO invites (id, table_id, inviter_id, invitee_id, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query,
		invite.ID, invite.TableID, invite.InviterID, invite.InviteeID,
		invite.Status, invite.CreatedAt, invite.UpdatedAt,
	)

	return err
}

// GetByID busca um convite por ID
func (r *InviteRepository) GetByID(id string) (*models.Invite, error) {
	var invite models.Invite

	query := `
		SELECT id, table_id, inviter_id, invitee_id, status, created_at, updated_at
		FROM invites 
		WHERE id = ?
	`

	err := r.db.Get(&invite, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &invite, nil
}

// GetByTableID busca convites por mesa
func (r *InviteRepository) GetByTableID(tableID string) ([]*models.InviteDetails, error) {
	var invites []*models.InviteDetails

	query := `
		SELECT 
			i.id, i.table_id, i.inviter_id, i.invitee_id, i.status, i.created_at, i.updated_at,
			inviter.email as inviter_email,
			invitee.email as invitee_email
		FROM invites i
		JOIN users inviter ON i.inviter_id = inviter.id
		JOIN users invitee ON i.invitee_id = invitee.id
		WHERE i.table_id = ?
		ORDER BY i.created_at DESC
	`

	rows, err := r.db.Query(query, tableID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var invite models.InviteDetails
		var inviterEmail, inviteeEmail string

		err := rows.Scan(
			&invite.ID, &invite.TableID, &invite.InviterID, &invite.InviteeID,
			&invite.Status, &invite.CreatedAt, &invite.UpdatedAt,
			&inviterEmail, &inviteeEmail,
		)
		if err != nil {
			return nil, err
		}

		invite.Inviter = &models.UserResponse{
			ID:    invite.InviterID,
			Email: inviterEmail,
		}
		invite.Invitee = &models.UserResponse{
			ID:    invite.InviteeID,
			Email: inviteeEmail,
		}

		invites = append(invites, &invite)
	}

	return invites, nil
}

// GetByUserID busca convites para um usuário (como invitee)
func (r *InviteRepository) GetByUserID(userID int, offset, limit int) ([]*models.InviteDetails, error) {
	var invites []*models.InviteDetails

	query := `
		SELECT 
			i.id, i.table_id, i.inviter_id, i.invitee_id, i.status, i.created_at, i.updated_at,
			gt.name as table_name, gt.system as table_system,
			inviter.email as inviter_email
		FROM invites i
		JOIN game_tables gt ON i.table_id = gt.id
		JOIN users inviter ON i.inviter_id = inviter.id
		WHERE i.invitee_id = ?
		ORDER BY i.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var invite models.InviteDetails
		var tableName, tableSystem, inviterEmail string

		err := rows.Scan(
			&invite.ID, &invite.TableID, &invite.InviterID, &invite.InviteeID,
			&invite.Status, &invite.CreatedAt, &invite.UpdatedAt,
			&tableName, &tableSystem, &inviterEmail,
		)
		if err != nil {
			return nil, err
		}

		invite.Table = &models.GameTable{
			ID:     invite.TableID,
			Name:   tableName,
			System: tableSystem,
		}
		invite.Inviter = &models.UserResponse{
			ID:    invite.InviterID,
			Email: inviterEmail,
		}

		invites = append(invites, &invite)
	}

	return invites, nil
}

// UpdateStatus atualiza o status de um convite
func (r *InviteRepository) UpdateStatus(id, status string) error {
	query := `
		UPDATE invites 
		SET status = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	result, err := r.db.Exec(query, status, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetUserByEmail busca usuário por email (usado para convites)
func (r *InviteRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	query := `SELECT id, email, created_at, updated_at FROM users WHERE email = ?`

	err := r.db.Get(&user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// CheckInviteExists verifica se já existe convite para essa mesa/usuário
func (r *InviteRepository) CheckInviteExists(tableID string, inviteeID int) (bool, error) {
	var count int

	query := `SELECT COUNT(*) FROM invites WHERE table_id = ? AND invitee_id = ?`

	err := r.db.Get(&count, query, tableID, inviteeID)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
