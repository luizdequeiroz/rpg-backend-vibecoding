-- +goose Up
-- Criar nova tabela users simplificada para autenticação
CREATE TABLE users_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    username VARCHAR(50),
    display_name VARCHAR(100),
    avatar_url VARCHAR(255),
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Copiar dados existentes (se houver)
INSERT INTO users_new (id, email, password_hash, username, display_name, avatar_url, is_active, created_at, updated_at)
SELECT id, email, password_hash, username, display_name, avatar_url, is_active, created_at, updated_at
FROM users;

-- Remover tabela antiga
DROP TABLE users;

-- Renomear nova tabela
ALTER TABLE users_new RENAME TO users;

-- Recriar índices
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_active ON users(is_active);

-- +goose Down
-- Reverter para estrutura original
CREATE TABLE users_old (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    display_name VARCHAR(100) NOT NULL,
    avatar_url VARCHAR(255),
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Copiar dados de volta (apenas registros válidos)
INSERT INTO users_old (id, email, password_hash, username, display_name, avatar_url, is_active, created_at, updated_at)
SELECT id, email, password_hash, 
       COALESCE(username, email), 
       COALESCE(display_name, email), 
       avatar_url, is_active, created_at, updated_at
FROM users
WHERE username IS NOT NULL AND display_name IS NOT NULL;

DROP TABLE users;
ALTER TABLE users_old RENAME TO users;

-- Recriar índices originais
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_active ON users(is_active);
