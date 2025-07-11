-- +goose Up
-- Criação da tabela de templates de ficha de personagem
CREATE TABLE sheet_templates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(100) NOT NULL,
    system VARCHAR(50) NOT NULL,
    definition TEXT NOT NULL, -- JSON com campos e regras de rolagem
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Índices para melhor performance
CREATE INDEX idx_sheet_templates_system ON sheet_templates(system);
CREATE INDEX idx_sheet_templates_active ON sheet_templates(is_active);
CREATE INDEX idx_sheet_templates_name ON sheet_templates(name);

-- +goose Down
-- Remover tabela de templates
DROP INDEX IF EXISTS idx_sheet_templates_name;
DROP INDEX IF EXISTS idx_sheet_templates_active;
DROP INDEX IF EXISTS idx_sheet_templates_system;
DROP TABLE sheet_templates;
