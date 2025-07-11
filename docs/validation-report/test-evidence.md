# Evidências dos Testes Realizados

## 1. Health Check
**Comando:** `curl -X GET "http://localhost:8080/health"`  
**Status:** ✅ SUCESSO  
**Resposta:**
```json
{
  "status": "ok",
  "timestamp": "2025-07-11T14:45:40.0323354-03:00",
  "version": "1.0.0",
  "services": {
    "database": "ok"
  }
}
```

## 2. Cadastro de Usuário
**Comando:** `curl -X POST "http://localhost:8080/api/v1/auth/signup"`  
**Body:**
```json
{
  "email": "test@validation.com",
  "password": "test123"
}
```
**Status:** ✅ SUCESSO  
**Resposta:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3RAdmFsaWRhdGlvbi5jb20iLCJleHAiOjE3NTIzNDIzNTMsImlhdCI6MTc1MjI1NTk1MywidXNlcl9pZCI6OH0.pYlbINpyEFa3DiBTHR_220wFclEvxT8eAFEum3fMUYU",
  "user": {
    "id": 8,
    "email": "test@validation.com",
    "created_at": "2025-07-11T14:45:53.5435888-03:00",
    "updated_at": "2025-07-11T14:45:53.5435888-03:00"
  }
}
```

## 3. Validação de Token JWT
**Comando:** `curl -X GET "http://localhost:8080/api/v1/auth/me"`  
**Header:** `Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3RAdmFsaWRhdGlvbi5jb20iLCJleHAiOjE3NTIzNDIzNTMsImlhdCI6MTc1MjI1NTk1MywidXNlcl9pZCI6OH0.pYlbINpyEFa3DiBTHR_220wFclEvxT8eAFEum3fMUYU`  
**Status:** ✅ SUCESSO  
**Resposta:**
```json
{
  "id": 8,
  "email": "test@validation.com",
  "created_at": "2025-07-11T14:45:53.5435888-03:00",
  "updated_at": "2025-07-11T14:45:53.5435888-03:00"
}
```

## 4. Listagem de Templates
**Comando:** `curl -X GET "http://localhost:8080/api/v1/templates"`  
**Status:** ✅ SUCESSO  
**Resposta:** 3 templates encontrados
```json
{
  "templates": [
    {
      "id": 4,
      "is_active": true,
      "created_at": "2025-07-11T14:17:43.4214683-03:00",
      "updated_at": "2025-07-11T14:17:43.4214683-03:00",
      "name": "Ficha D&D 5e",
      "system": "D&D 5e",
      "description": "Template para D&D 5e",
      "definition": {
        "sections": [
          {
            "fields": [
              {"name": "strength", "type": "number"},
              {"name": "dexterity", "type": "number"}
            ],
            "name": "Atributos"
          }
        ]
      }
    }
  ],
  "total": 3,
  "page": 1,
  "limit": 20
}
```

## 5. Criação de Mesa de Jogo
**Comando:** `curl -X POST "http://localhost:8080/api/v1/tables/"`  
**Body:**
```json
{
  "name": "Mesa de Validação",
  "system": "D&D 5e"
}
```
**Status:** ✅ SUCESSO  
**Resposta:**
```json
{
  "id": "ba66eb05-326b-4480-bda2-5e36e5f5d317",
  "name": "Mesa de Validação",
  "system": "D&D 5e",
  "owner_id": 8,
  "created_at": "2025-07-11T14:46:52.8932614-03:00",
  "updated_at": "2025-07-11T14:46:52.8932614-03:00"
}
```

## 6. Listagem de Mesas
**Comando:** `curl -X GET "http://localhost:8080/api/v1/tables/"`  
**Status:** ✅ SUCESSO  
**Resposta:**
```json
{
  "limit": 20,
  "page": 1,
  "tables": [
    {
      "id": "ba66eb05-326b-4480-bda2-5e36e5f5d317",
      "name": "Mesa de Validação",
      "system": "D&D 5e",
      "owner_id": 8,
      "owner": {
        "id": 8,
        "email": "test@validation.com"
      },
      "role": "owner",
      "created_at": "2025-07-11T14:46:52.8932614-03:00",
      "updated_at": "2025-07-11T14:46:52.8932614-03:00"
    }
  ],
  "total": 1
}
```

## 7. Status das Migrações
**Comando:** `go run cmd/migrate/main.go -action=status`  
**Status:** ✅ TODAS APLICADAS  
**Resultado:**
```
2025/07/11 14:47:36     Applied At                  Migration
2025/07/11 14:47:36     =======================================
2025/07/11 14:47:36     Thu Jul 10 18:29:38 2025 -- 00001_create_users_table.sql
2025/07/11 14:47:36     Thu Jul 10 18:29:38 2025 -- 00002_create_campaigns_table.sql
2025/07/11 14:47:36     Thu Jul 10 18:29:38 2025 -- 00003_create_characters_table.sql
2025/07/11 14:47:36     Thu Jul 10 18:33:34 2025 -- 20250710183236_add_sessions_table.sql
2025/07/11 14:47:36     Thu Jul 10 19:03:45 2025 -- 20250710190310_update_users_for_auth.sql
2025/07/11 14:47:36     Thu Jul 10 19:06:45 2025 -- 20250710190617_create_sheet_templates_table.sql
2025/07/11 14:47:36     Fri Jul 11 14:54:52 2025 -- 20250711120000_create_game_tables.sql
2025/07/11 14:47:36     Fri Jul 11 14:54:52 2025 -- 20250711120100_create_invites.sql
2025/07/11 14:47:36     Fri Jul 11 17:08:24 2025 -- 20250711160000_create_player_sheets.sql
2025/07/11 14:47:36     Fri Jul 11 17:08:24 2025 -- 20250711160100_create_rolls.sql
2025/07/11 14:47:36     Fri Jul 11 17:38:14 2025 -- 20250711160101_create_rolls_table.sql
```

## 8. Endpoints Disponíveis
**Comando:** `curl -X GET "http://localhost:8080/"`  
**Status:** ✅ SUCESSO  
**Endpoints Registrados:**
```json
{
  "description": "Backend para gerenciamento de sessões de RPG",
  "endpoints": {
    "api_v1": "/api/v1",
    "auth": {
      "login": "/api/v1/auth/login",
      "me": "/api/v1/auth/me",
      "signup": "/api/v1/auth/signup"
    },
    "campaigns": "/api/v1/campaigns",
    "characters": "/api/v1/characters",
    "docs": "/docs/index.html",
    "health": "/health",
    "sessions": "/api/v1/sessions",
    "templates": "/api/v1/templates",
    "users": "/api/v1/users"
  },
  "name": "RPG Session Backend",
  "version": "1.0.0"
}
```

## Logs do Servidor
```
2025/07/11 14:41:39 Configurações carregadas: Host=localhost, Port=8080
2025/07/11 14:41:39 goose: no migrations to run. current version: 20250711160101
2025/07/11 14:41:39 Servidor RPG Backend iniciado com sucesso!
2025/07/11 14:41:39 Banco de dados: file:./data/rpg.db?cache=shared&mode=rwc
2025/07/11 14:41:39 Servidor rodando em: http://localhost:8080
2025/07/11 14:41:39 Healthcheck: http://localhost:8080/health
2025/07/11 14:41:39 API v1: http://localhost:8080/api/v1
[GIN] 2025/07/11 - 14:45:40 | 200 |       505.2µs |       127.0.0.1 | GET      "/health"
[GIN] 2025/07/11 - 14:45:42 | 200 |            0s |       127.0.0.1 | GET      "/"
[GIN] 2025/07/11 - 14:45:53 | 201 |     57.657ms |       127.0.0.1 | POST     "/api/v1/auth/signup"
[GIN] 2025/07/11 - 14:45:58 | 200 |       507.8µs |       127.0.0.1 | GET      "/api/v1/auth/me"
[GIN] 2025/07/11 - 14:46:05 | 200 |     15.8054ms |       127.0.0.1 | GET      "/api/v1/templates"
[GIN] 2025/07/11 - 14:46:52 | 201 |     25.9045ms |       127.0.0.1 | POST     "/api/v1/tables/"
[GIN] 2025/07/11 - 14:47:05 | 200 |      7.0053ms |       127.0.0.1 | GET      "/api/v1/tables/"
```
