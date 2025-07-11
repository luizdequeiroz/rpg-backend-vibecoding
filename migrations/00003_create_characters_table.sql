-- +goose Up
-- Criação da tabela de personagens
CREATE TABLE characters (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(100) NOT NULL,
    class VARCHAR(50) NOT NULL,
    race VARCHAR(50) NOT NULL,
    level INTEGER DEFAULT 1,
    experience_points INTEGER DEFAULT 0,
    max_health INTEGER NOT NULL,
    current_health INTEGER NOT NULL,
    armor_class INTEGER DEFAULT 10,
    initiative_bonus INTEGER DEFAULT 0,
    
    -- Atributos básicos
    strength INTEGER DEFAULT 10,
    dexterity INTEGER DEFAULT 10,
    constitution INTEGER DEFAULT 10,
    intelligence INTEGER DEFAULT 10,
    wisdom INTEGER DEFAULT 10,
    charisma INTEGER DEFAULT 10,
    
    -- Relacionamentos
    user_id INTEGER NOT NULL,
    campaign_id INTEGER NOT NULL,
    
    -- Dados adicionais como JSON (flexibilidade para diferentes sistemas)
    additional_data TEXT, -- JSON com skills, equipment, spells, etc.
    
    -- Metadados
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (campaign_id) REFERENCES campaigns(id) ON DELETE CASCADE
);

-- Índices
CREATE INDEX idx_characters_user ON characters(user_id);
CREATE INDEX idx_characters_campaign ON characters(campaign_id);
CREATE INDEX idx_characters_level ON characters(level);
CREATE INDEX idx_characters_active ON characters(is_active);

-- +goose Down
-- Remoção da tabela de personagens
DROP INDEX IF EXISTS idx_characters_active;
DROP INDEX IF EXISTS idx_characters_level;
DROP INDEX IF EXISTS idx_characters_campaign;
DROP INDEX IF EXISTS idx_characters_user;
DROP TABLE characters;
