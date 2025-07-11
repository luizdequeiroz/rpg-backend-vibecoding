-- +goose Up
-- +goose StatementBegin
CREATE TABLE rolls (
    id VARCHAR(36) PRIMARY KEY DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-' || '4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('ab89',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
    sheet_id VARCHAR(36) NOT NULL,
    table_id VARCHAR(36) NOT NULL,
    user_id INTEGER NOT NULL,
    expression VARCHAR(200) NOT NULL,
    field_name VARCHAR(100),
    result_value INTEGER NOT NULL,
    result_details TEXT, -- JSON com detalhes da rolagem
    success BOOLEAN,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (sheet_id) REFERENCES player_sheets(id) ON DELETE CASCADE,
    FOREIGN KEY (table_id) REFERENCES game_tables(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Índices para performance
CREATE INDEX idx_rolls_sheet_id ON rolls(sheet_id);
CREATE INDEX idx_rolls_table_id ON rolls(table_id);
CREATE INDEX idx_rolls_user_id ON rolls(user_id);
CREATE INDEX idx_rolls_created_at ON rolls(created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_rolls_created_at;
DROP INDEX IF EXISTS idx_rolls_user_id;
DROP INDEX IF EXISTS idx_rolls_table_id;
DROP INDEX IF EXISTS idx_rolls_sheet_id;
DROP TABLE IF EXISTS rolls;
-- +goose StatementEnd
