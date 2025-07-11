-- +goose Up
-- Tabela para armazenar rolagens de dados
CREATE TABLE IF NOT EXISTS rolls (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sheet_id INTEGER,
    table_id TEXT,
    user_id INTEGER NOT NULL,
    expression TEXT NOT NULL,
    field_name TEXT,
    result_value INTEGER NOT NULL,
    result_details TEXT, -- JSON com detalhes da rolagem
    comment TEXT,
    success BOOLEAN,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (sheet_id) REFERENCES player_sheets(rowid) ON DELETE SET NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Índices para melhorar performance
CREATE INDEX IF NOT EXISTS idx_rolls_user_id ON rolls(user_id);
CREATE INDEX IF NOT EXISTS idx_rolls_sheet_id ON rolls(sheet_id);
CREATE INDEX IF NOT EXISTS idx_rolls_table_id ON rolls(table_id);
CREATE INDEX IF NOT EXISTS idx_rolls_created_at ON rolls(created_at);

-- +goose Down
DROP INDEX IF EXISTS idx_rolls_created_at;
DROP INDEX IF EXISTS idx_rolls_table_id;
DROP INDEX IF EXISTS idx_rolls_sheet_id;
DROP INDEX IF EXISTS idx_rolls_user_id;
DROP TABLE IF EXISTS rolls;
