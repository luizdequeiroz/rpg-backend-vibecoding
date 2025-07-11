-- +goose Up
-- +goose StatementBegin
CREATE TABLE player_sheets (
    id VARCHAR(36) PRIMARY KEY DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-' || '4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('ab89',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
    table_id VARCHAR(36) NOT NULL,
    template_id INTEGER NOT NULL,
    owner_id INTEGER NOT NULL,
    name VARCHAR(100) NOT NULL,
    data TEXT NOT NULL DEFAULT '{}', -- JSON como string
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (table_id) REFERENCES game_tables(id) ON DELETE CASCADE,
    FOREIGN KEY (template_id) REFERENCES sheet_templates(id) ON DELETE RESTRICT,
    FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Índices para performance
CREATE INDEX idx_player_sheets_table_id ON player_sheets(table_id);
CREATE INDEX idx_player_sheets_owner_id ON player_sheets(owner_id);
CREATE INDEX idx_player_sheets_template_id ON player_sheets(template_id);

-- Trigger para atualizar updated_at
CREATE TRIGGER update_player_sheets_updated_at 
    AFTER UPDATE ON player_sheets
    FOR EACH ROW
    BEGIN
        UPDATE player_sheets SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_player_sheets_updated_at;
DROP INDEX IF EXISTS idx_player_sheets_template_id;
DROP INDEX IF EXISTS idx_player_sheets_owner_id;
DROP INDEX IF EXISTS idx_player_sheets_table_id;
DROP TABLE IF EXISTS player_sheets;
-- +goose StatementEnd
