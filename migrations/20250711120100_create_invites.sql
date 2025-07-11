-- +goose Up
-- Criação da tabela de convites para mesas
CREATE TABLE invites (
    id TEXT PRIMARY KEY, -- UUID como string
    table_id TEXT NOT NULL,
    inviter_id INTEGER NOT NULL, -- Quem convidou (owner da mesa)
    invitee_id INTEGER NOT NULL, -- Quem foi convidado
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'accepted', 'declined')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (table_id) REFERENCES game_tables(id) ON DELETE CASCADE,
    FOREIGN KEY (inviter_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (invitee_id) REFERENCES users(id) ON DELETE CASCADE,
    -- Previne convites duplicados para a mesma mesa/usuário
    UNIQUE(table_id, invitee_id)
);

-- Índices para performance
CREATE INDEX idx_invites_table ON invites(table_id);
CREATE INDEX idx_invites_inviter ON invites(inviter_id);
CREATE INDEX idx_invites_invitee ON invites(invitee_id);
CREATE INDEX idx_invites_status ON invites(status);
CREATE INDEX idx_invites_created ON invites(created_at);

-- +goose Down
-- Remoção da tabela de convites
DROP INDEX IF EXISTS idx_invites_created;
DROP INDEX IF EXISTS idx_invites_status;
DROP INDEX IF EXISTS idx_invites_invitee;
DROP INDEX IF EXISTS idx_invites_inviter;
DROP INDEX IF EXISTS idx_invites_table;
DROP TABLE IF EXISTS invites;
