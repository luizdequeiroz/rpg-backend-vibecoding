# Manual de Uso do Sistema RPG Backend

## 1. Iniciando o Sistema

### 1.1 Pré-requisitos
- Go 1.21 ou superior
- Git

### 1.2 Instalação
```bash
git clone https://github.com/luizdequeiroz/rpg-backend.git
cd rpg-backend
go mod download
```

### 1.3 Configuração do Banco
```bash
# Executar migrações
go run cmd/migrate/main.go -action=up

# Verificar status
go run cmd/migrate/main.go -action=status
```

### 1.4 Iniciando o Servidor
```bash
# Compilar
go build -o bin/api cmd/api/main.go

# Executar
./bin/api
```

**Servidor iniciará em:** `http://localhost:8080`

## 2. Autenticação

### 2.1 Cadastro de Usuário
```bash
curl -X POST "http://localhost:8080/api/v1/auth/signup" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "usuario@exemplo.com",
    "password": "minhasenha123"
  }'
```

**Resposta:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "usuario@exemplo.com",
    "created_at": "2025-07-11T14:45:53Z",
    "updated_at": "2025-07-11T14:45:53Z"
  }
}
```

### 2.2 Login
```bash
curl -X POST "http://localhost:8080/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "usuario@exemplo.com",
    "password": "minhasenha123"
  }'
```

### 2.3 Verificar Perfil
```bash
curl -X GET "http://localhost:8080/api/v1/auth/me" \
  -H "Authorization: Bearer SEU_TOKEN_AQUI"
```

## 3. Templates de Fichas

### 3.1 Listar Templates
```bash
curl -X GET "http://localhost:8080/api/v1/templates"
```

### 3.2 Criar Template
```bash
curl -X POST "http://localhost:8080/api/v1/templates" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "D&D 5e Custom",
    "system": "D&D 5e",
    "description": "Template personalizado para D&D 5e",
    "definition": {
      "sections": [
        {
          "name": "Atributos",
          "fields": [
            {"name": "strength", "type": "number", "min": 1, "max": 20},
            {"name": "dexterity", "type": "number", "min": 1, "max": 20},
            {"name": "constitution", "type": "number", "min": 1, "max": 20}
          ]
        },
        {
          "name": "Skills",
          "fields": [
            {"name": "athletics", "type": "number"},
            {"name": "acrobatics", "type": "number"}
          ]
        }
      ]
    }
  }'
```

## 4. Mesas de Jogo

### 4.1 Criar Mesa
```bash
curl -X POST "http://localhost:8080/api/v1/tables/" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer SEU_TOKEN" \
  -d '{
    "name": "Mesa D&D: A Busca pelo Artefato",
    "system": "D&D 5e"
  }'
```

### 4.2 Listar Minhas Mesas
```bash
curl -X GET "http://localhost:8080/api/v1/tables/" \
  -H "Authorization: Bearer SEU_TOKEN"
```

### 4.3 Detalhes da Mesa
```bash
curl -X GET "http://localhost:8080/api/v1/tables/UUID_DA_MESA" \
  -H "Authorization: Bearer SEU_TOKEN"
```

### 4.4 Criar Convite
```bash
curl -X POST "http://localhost:8080/api/v1/tables/UUID_DA_MESA/invites" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer SEU_TOKEN" \
  -d '{
    "invitee_email": "jogador@exemplo.com"
  }'
```

### 4.5 Aceitar Convite
```bash
curl -X POST "http://localhost:8080/api/v1/tables/UUID_DA_MESA/invites/UUID_DO_CONVITE/accept" \
  -H "Authorization: Bearer SEU_TOKEN"
```

## 5. Fichas de Personagem

### 5.1 Criar Ficha
```bash
curl -X POST "http://localhost:8080/api/v1/sheets/" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer SEU_TOKEN" \
  -d '{
    "table_id": "UUID_DA_MESA",
    "template_id": 1,
    "name": "Aragorn, o Montanhês",
    "data": {
      "attributes": {
        "strength": 16,
        "dexterity": 14,
        "constitution": 15,
        "intelligence": 12,
        "wisdom": 13,
        "charisma": 14
      },
      "skills": {
        "athletics": 8,
        "survival": 6,
        "investigation": 3
      },
      "combat": {
        "armor_class": 16,
        "hit_points": 42,
        "speed": 30
      }
    }
  }'
```

### 5.2 Listar Fichas da Mesa
```bash
curl -X GET "http://localhost:8080/api/v1/sheets/?table_id=UUID_DA_MESA" \
  -H "Authorization: Bearer SEU_TOKEN"
```

### 5.3 Detalhes da Ficha
```bash
curl -X GET "http://localhost:8080/api/v1/sheets/UUID_DA_FICHA" \
  -H "Authorization: Bearer SEU_TOKEN"
```

### 5.4 Atualizar Ficha
```bash
curl -X PUT "http://localhost:8080/api/v1/sheets/UUID_DA_FICHA" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer SEU_TOKEN" \
  -d '{
    "name": "Aragorn, o Rei",
    "data": {
      "attributes": {
        "strength": 18,
        "dexterity": 14,
        "constitution": 16
      }
    }
  }'
```

## 6. Sistema de Rolagem de Dados

### 6.1 Rolagem Simples
```bash
curl -X POST "http://localhost:8080/api/v1/dice/roll" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer SEU_TOKEN" \
  -d '{
    "expression": "1d20+3",
    "comment": "Teste de Força"
  }'
```

**Resposta:**
```json
{
  "id": 1,
  "expression": "1d20+3",
  "result": 18,
  "details": "[15] + 3 = 18",
  "is_critical": false,
  "is_fumble": false,
  "comment": "Teste de Força",
  "user_id": 1,
  "created_at": "2025-07-11T14:30:00Z"
}
```

### 6.2 Rolagem com Ficha
```bash
curl -X POST "http://localhost:8080/api/v1/dice/roll-with-sheet" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer SEU_TOKEN" \
  -d '{
    "sheet_id": "UUID_DA_FICHA",
    "expression": "1d20+{strength}",
    "attribute_field": "strength",
    "comment": "Teste de Força usando atributo da ficha"
  }'
```

### 6.3 Histórico de Rolagens
```bash
curl -X GET "http://localhost:8080/api/v1/dice/history?page=1&limit=10" \
  -H "Authorization: Bearer SEU_TOKEN"
```

### 6.4 Rolagens da Ficha
```bash
curl -X GET "http://localhost:8080/api/v1/rolls/sheet/UUID_DA_FICHA" \
  -H "Authorization: Bearer SEU_TOKEN"
```

### 6.5 Rolagens da Mesa
```bash
curl -X GET "http://localhost:8080/api/v1/rolls/table/UUID_DA_MESA" \
  -H "Authorization: Bearer SEU_TOKEN"
```

## 7. Expressões de Dados Suportadas

### 7.1 Formatos Básicos
- `1d20` - Um dado de 20 lados
- `2d6` - Dois dados de 6 lados
- `3d8+5` - Três dados de 8 faces mais 5
- `1d100-10` - Um dado de 100 faces menos 10

### 7.2 Limites
- **Número de dados:** 1 a 100
- **Número de lados:** 2 a 1000
- **Modificador:** -999 a +999

### 7.3 Expressões com Atributos
- `1d20+{strength}` - Usa valor do atributo strength
- `2d6+{dexterity}` - Usa valor do atributo dexterity
- `1d100+{skills.athletics}` - Usa valor aninhado

### 7.4 Detecção Especial
- **Crítico:** d20 = 20 natural
- **Fumble:** d20 = 1 natural

## 8. Códigos de Status

### 8.1 Sucesso
- `200 OK` - Operação bem-sucedida
- `201 Created` - Recurso criado

### 8.2 Erro do Cliente
- `400 Bad Request` - Dados inválidos
- `401 Unauthorized` - Token inválido/ausente
- `403 Forbidden` - Sem permissão
- `404 Not Found` - Recurso não encontrado

### 8.3 Erro do Servidor
- `500 Internal Server Error` - Erro interno

## 9. Exemplos de Uso Completo

### 9.1 Fluxo Completo de Jogo
```bash
# 1. Criar usuário
curl -X POST "http://localhost:8080/api/v1/auth/signup" \
  -H "Content-Type: application/json" \
  -d '{"email": "mestre@rpg.com", "password": "senha123"}'

# 2. Fazer login (salvar o token)
TOKEN=$(curl -s -X POST "http://localhost:8080/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email": "mestre@rpg.com", "password": "senha123"}' | \
  jq -r '.token')

# 3. Criar mesa
MESA_ID=$(curl -s -X POST "http://localhost:8080/api/v1/tables/" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name": "Mesa D&D Aventura", "system": "D&D 5e"}' | \
  jq -r '.id')

# 4. Criar ficha
FICHA_ID=$(curl -s -X POST "http://localhost:8080/api/v1/sheets/" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"table_id\": \"$MESA_ID\",
    \"template_id\": 1,
    \"name\": \"Gandalf, o Cinzento\",
    \"data\": {
      \"attributes\": {\"strength\": 12, \"dexterity\": 14, \"wisdom\": 18},
      \"skills\": {\"arcana\": 10, \"investigation\": 8}
    }
  }" | jq -r '.id')

# 5. Fazer rolagem
curl -X POST "http://localhost:8080/api/v1/dice/roll-with-sheet" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"sheet_id\": \"$FICHA_ID\",
    \"expression\": \"1d20+{wisdom}\",
    \"attribute_field\": \"wisdom\",
    \"comment\": \"Teste de Sabedoria\"
  }"
```

## 10. Troubleshooting

### 10.1 Problemas Comuns
- **Token expirado:** Fazer login novamente
- **404 Not Found:** Verificar se o endpoint tem `/` no final
- **401 Unauthorized:** Verificar se o header Authorization está correto
- **400 Bad Request:** Verificar se o JSON está bem formado

### 10.2 Health Check
```bash
curl -X GET "http://localhost:8080/health"
```

### 10.3 Logs do Servidor
- Verificar logs no terminal onde o servidor está rodando
- Códigos de status HTTP são registrados automaticamente

## 11. Documentação Interativa

### 11.1 Swagger UI
Acesse: `http://localhost:8080/docs/index.html`

### 11.2 Lista de Endpoints
Acesse: `http://localhost:8080/`
