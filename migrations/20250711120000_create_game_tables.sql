-- +goose Up
-- Criação da tabela de mesas de jogo
CREATE TABLE game_tables (
    id TEXT PRIMARY KEY, -- UUID como string
    name VARCHAR(100) NOT NULL,
    system VARCHAR(50) NOT NULL, -- D&D 5e, Pathfinder, Call of Cthulhu, etc.
    owner_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Índices para performance
CREATE INDEX idx_game_tables_owner ON game_tables(owner_id);
CREATE INDEX idx_game_tables_system ON game_tables(system);
CREATE INDEX idx_game_tables_created ON game_tables(created_at);

-- +goose Down
-- Remoção da tabela de mesas de jogo
DROP INDEX IF EXISTS idx_game_tables_created;
DROP INDEX IF EXISTS idx_game_tables_system;
DROP INDEX IF EXISTS idx_game_tables_owner;
DROP TABLE IF EXISTS game_tables;
