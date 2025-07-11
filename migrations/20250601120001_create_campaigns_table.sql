-- +goose Up
-- Criação da tabela de campanhas
CREATE TABLE campaigns (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    system VARCHAR(50) NOT NULL, -- D&D 5e, Pathfinder, etc.
    max_players INTEGER DEFAULT 6,
    current_players INTEGER DEFAULT 0,
    master_id INTEGER NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    is_public BOOLEAN DEFAULT FALSE,
    session_frequency VARCHAR(50), -- Semanal, Quinzenal, etc.
    next_session DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (master_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Índices
CREATE INDEX idx_campaigns_master ON campaigns(master_id);
CREATE INDEX idx_campaigns_active ON campaigns(is_active);
CREATE INDEX idx_campaigns_public ON campaigns(is_public);
CREATE INDEX idx_campaigns_system ON campaigns(system);

-- +goose Down
-- Remoção da tabela de campanhas
DROP INDEX IF EXISTS idx_campaigns_system;
DROP INDEX IF EXISTS idx_campaigns_public;
DROP INDEX IF EXISTS idx_campaigns_active;
DROP INDEX IF EXISTS idx_campaigns_master;
DROP TABLE campaigns;
