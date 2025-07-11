-- +goose Up
-- Criação da tabela de sessões de jogo
CREATE TABLE game_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    campaign_id INTEGER NOT NULL,
    session_number INTEGER NOT NULL,
    title VARCHAR(200) NOT NULL,
    summary TEXT,
    session_date DATETIME NOT NULL,
    duration_minutes INTEGER,
    experience_awarded INTEGER DEFAULT 0,
    is_completed BOOLEAN DEFAULT FALSE,
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (campaign_id) REFERENCES campaigns(id) ON DELETE CASCADE
);

-- Índices
CREATE INDEX idx_sessions_campaign ON game_sessions(campaign_id);
CREATE INDEX idx_sessions_date ON game_sessions(session_date);
CREATE INDEX idx_sessions_completed ON game_sessions(is_completed);

-- +goose Down
-- Remoção da tabela de sessões
DROP INDEX IF EXISTS idx_sessions_completed;
DROP INDEX IF EXISTS idx_sessions_date;
DROP INDEX IF EXISTS idx_sessions_campaign;
DROP TABLE game_sessions;
