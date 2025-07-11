-- +goose Up
-- +goose StatementBegin
-- Remover campo system da tabela sheet_templates
-- SQLite não suporta DROP COLUMN, então recriaremos a tabela
CREATE TABLE sheet_templates_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(100) NOT NULL,
    definition TEXT NOT NULL, -- JSON com campos e regras de rolagem
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Copiar dados existentes (sem o campo system)
INSERT INTO sheet_templates_new (id, name, definition, description, is_active, created_at, updated_at)
SELECT id, name, definition, description, is_active, created_at, updated_at 
FROM sheet_templates;

-- Remover tabela antiga
DROP TABLE sheet_templates;

-- Renomear nova tabela
ALTER TABLE sheet_templates_new RENAME TO sheet_templates;

-- Recriar índices (sem o sistema)
CREATE INDEX idx_sheet_templates_active ON sheet_templates(is_active);
CREATE INDEX idx_sheet_templates_name ON sheet_templates(name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Adicionar campo system de volta
CREATE TABLE sheet_templates_old (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(100) NOT NULL,
    system VARCHAR(50) NOT NULL DEFAULT 'Generic',
    definition TEXT NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Copiar dados existentes (adicionando system padrão)
INSERT INTO sheet_templates_old (id, name, system, definition, description, is_active, created_at, updated_at)
SELECT id, name, 'Generic', definition, description, is_active, created_at, updated_at 
FROM sheet_templates;

-- Remover tabela atual
DROP TABLE sheet_templates;

-- Renomear tabela
ALTER TABLE sheet_templates_old RENAME TO sheet_templates;

-- Recriar índices
CREATE INDEX idx_sheet_templates_system ON sheet_templates(system);
CREATE INDEX idx_sheet_templates_active ON sheet_templates(is_active);
CREATE INDEX idx_sheet_templates_name ON sheet_templates(name);
-- +goose StatementEnd
