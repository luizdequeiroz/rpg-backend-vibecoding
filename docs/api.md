# API Reference - RPG Backend

## 🌐 Base URL
```
http://localhost:8080
```

## 🔐 Autenticação

A API usa autenticação JWT (JSON Web Tokens). Para acessar endpoints protegidos, inclua o header:

```
Authorization: Bearer <seu_jwt_token>
```

### Obter Token JWT
Faça login através do endpoint `/api/v1/auth/login` para receber um token.

---

## 📋 Endpoints

### 🏥 Health Check

#### `GET /health`
Verifica a saúde da aplicação e seus serviços.

**Resposta:**
```json
{
  "status": "ok",
  "timestamp": "2025-07-11T09:26:29Z",
  "version": "1.0.0",
  "services": {
    "database": "ok"
  }
}
```

---

### 🔐 Autenticação

#### `POST /api/v1/auth/signup`
Registra um novo usuário no sistema.

**Request Body:**
```json
{
  "email": "usuario@exemplo.com",
  "password": "senha123"
}
```

**Resposta (201):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "usuario@exemplo.com",
    "created_at": "2025-07-11T09:26:29Z",
    "updated_at": "2025-07-11T09:26:29Z"
  }
}
```

**Erros:**
- `400` - Email já existe ou dados inválidos
- `500` - Erro interno do servidor

---

#### `POST /api/v1/auth/login`
Autentica um usuário e retorna um token JWT.

**Request Body:**
```json
{
  "email": "usuario@exemplo.com",
  "password": "senha123"
}
```

**Resposta (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "usuario@exemplo.com",
    "created_at": "2025-07-11T09:26:29Z",
    "updated_at": "2025-07-11T09:26:29Z"
  }
}
```

**Erros:**
- `400` - Dados inválidos
- `401` - Credenciais incorretas
- `500` - Erro interno do servidor

---

#### `GET /api/v1/auth/me` 🔒
Retorna informações do usuário autenticado.

**Headers:**
```
Authorization: Bearer <token>
```

**Resposta (200):**
```json
{
  "id": 1,
  "email": "usuario@exemplo.com",
  "created_at": "2025-07-11T09:26:29Z",
  "updated_at": "2025-07-11T09:26:29Z"
}
```

**Erros:**
- `401` - Token inválido ou expirado
- `500` - Erro interno do servidor

---

### 👥 Usuários

#### `GET /api/v1/users`
Lista todos os usuários (versão pública).

**Query Parameters:**
- `page` (opcional) - Número da página (padrão: 1)
- `limit` (opcional) - Itens por página (padrão: 20)

**Resposta (200):**
```json
{
  "users": [
    {
      "id": 1,
      "email": "usuario1@exemplo.com",
      "created_at": "2025-07-11T09:26:29Z",
      "updated_at": "2025-07-11T09:26:29Z"
    },
    {
      "id": 2,
      "email": "usuario2@exemplo.com",
      "created_at": "2025-07-11T09:26:29Z",
      "updated_at": "2025-07-11T09:26:29Z"
    }
  ],
  "total": 2
}
```

---

#### `GET /api/v1/users/protected` 🔒
Lista usuários com informações adicionais (apenas para usuários autenticados).

**Headers:**
```
Authorization: Bearer <token>
```

**Resposta (200):**
```json
{
  "users": [
    {
      "id": 1,
      "email": "usuario1@exemplo.com",
      "created_at": "2025-07-11T09:26:29Z",
      "updated_at": "2025-07-11T09:26:29Z"
    }
  ],
  "total": 2,
  "requested_by_user": "admin@exemplo.com",
  "requested_by_id": 1,
  "additional_details": "Informações extras para usuários autenticados"
}
```

**Erros:**
- `401` - Token inválido ou expirado

---

### 📄 Templates de Ficha

#### `GET /api/v1/templates`
Lista todos os templates de ficha ativos.

**Query Parameters:**
- `page` (opcional) - Número da página (padrão: 1)
- `limit` (opcional) - Itens por página (padrão: 20)

**Resposta (200):**
```json
{
  "templates": [
    {
      "id": 1,
      "name": "Ficha D&D 5e",
      "system": "D&D 5e",
      "description": "Template para personagens de D&D 5ª edição",
      "definition": {
        "sections": [
          {
            "name": "Atributos",
            "fields": [
              {
                "name": "Força",
                "type": "number",
                "min": 1,
                "max": 20
              }
            ]
          }
        ]
      },
      "is_active": true,
      "created_at": "2025-07-11T09:26:29Z",
      "updated_at": "2025-07-11T09:26:29Z"
    }
  ],
  "total": 1,
  "page": 1,
  "limit": 20,
  "has_next": false
}
```

---

#### `GET /api/v1/templates/{id}`
Busca um template específico por ID.

**Path Parameters:**
- `id` - ID do template

**Resposta (200):**
```json
{
  "id": 1,
  "name": "Ficha D&D 5e",
  "system": "D&D 5e",
  "description": "Template para personagens de D&D 5ª edição",
  "definition": {
    "sections": [
      {
        "name": "Atributos",
        "fields": [
          {
            "name": "Força",
            "type": "number",
            "min": 1,
            "max": 20
          }
        ]
      }
    ]
  },
  "is_active": true,
  "created_at": "2025-07-11T09:26:29Z",
  "updated_at": "2025-07-11T09:26:29Z"
}
```

**Erros:**
- `400` - ID inválido
- `404` - Template não encontrado
- `500` - Erro interno do servidor

---

#### `POST /api/v1/templates`
Cria um novo template de ficha.

**Request Body:**
```json
{
  "name": "Ficha D&D 5e",
  "system": "D&D 5e",
  "description": "Template completo para personagens de D&D 5ª edição",
  "definition": {
    "sections": [
      {
        "name": "Atributos",
        "fields": [
          {
            "name": "Força",
            "type": "number",
            "min": 1,
            "max": 20,
            "default": 10
          },
          {
            "name": "Destreza",
            "type": "number",
            "min": 1,
            "max": 20,
            "default": 10
          }
        ]
      },
      {
        "name": "Informações Básicas",
        "fields": [
          {
            "name": "Nome",
            "type": "text",
            "required": true
          },
          {
            "name": "Classe",
            "type": "select",
            "options": ["Guerreiro", "Mago", "Ladino"],
            "required": true
          }
        ]
      }
    ],
    "rules": {
      "calculations": [
        {
          "field": "modificador_forca",
          "formula": "floor((forca - 10) / 2)"
        }
      ]
    }
  }
}
```

**Resposta (201):**
```json
{
  "id": 2,
  "name": "Ficha D&D 5e",
  "system": "D&D 5e",
  "description": "Template completo para personagens de D&D 5ª edição",
  "definition": { ... },
  "is_active": true,
  "created_at": "2025-07-11T09:26:29Z",
  "updated_at": "2025-07-11T09:26:29Z"
}
```

**Erros:**
- `400` - Dados inválidos ou definition JSON malformado
- `500` - Erro interno do servidor

---

#### `PUT /api/v1/templates/{id}`
Atualiza um template existente.

**Path Parameters:**
- `id` - ID do template

**Request Body (todos os campos opcionais):**
```json
{
  "name": "Ficha D&D 5e Atualizada",
  "system": "D&D 5e",
  "description": "Template atualizado com novos campos",
  "definition": {
    "sections": [
      {
        "name": "Atributos Básicos",
        "fields": [
          {
            "name": "Força",
            "type": "number",
            "min": 1,
            "max": 20
          }
        ]
      }
    ]
  }
}
```

**Resposta (200):**
```json
{
  "id": 1,
  "name": "Ficha D&D 5e Atualizada",
  "system": "D&D 5e",
  "description": "Template atualizado com novos campos",
  "definition": { ... },
  "is_active": true,
  "created_at": "2025-07-11T09:26:29Z",
  "updated_at": "2025-07-11T09:32:15Z"
}
```

**Erros:**
- `400` - ID ou dados inválidos
- `404` - Template não encontrado
- `500` - Erro interno do servidor

---

#### `DELETE /api/v1/templates/{id}`
Remove um template (soft delete).

**Path Parameters:**
- `id` - ID do template

**Resposta (204):**
```
No Content
```

**Erros:**
- `400` - ID inválido
- `404` - Template não encontrado
- `500` - Erro interno do servidor

---

## 📊 Códigos de Status HTTP

| Código | Significado | Uso |
|--------|-------------|-----|
| `200` | OK | Requisição bem-sucedida |
| `201` | Created | Recurso criado com sucesso |
| `204` | No Content | Operação bem-sucedida sem conteúdo |
| `400` | Bad Request | Dados inválidos na requisição |
| `401` | Unauthorized | Token inválido ou ausente |
| `403` | Forbidden | Acesso negado |
| `404` | Not Found | Recurso não encontrado |
| `409` | Conflict | Conflito (ex: email já existe) |
| `422` | Unprocessable Entity | Dados válidos mas regra de negócio violada |
| `500` | Internal Server Error | Erro interno do servidor |

---

## 📝 Estruturas de Dados

### Definition Schema (Templates)
O campo `definition` em templates suporta a seguinte estrutura:

```json
{
  "sections": [
    {
      "name": "Nome da Seção",
      "description": "Descrição opcional",
      "fields": [
        {
          "name": "Nome do Campo",
          "type": "text|number|boolean|select|textarea",
          "required": true,
          "default": "Valor padrão",
          "min": 1,
          "max": 20,
          "options": ["Opção 1", "Opção 2"],
          "placeholder": "Digite aqui...",
          "help": "Texto de ajuda"
        }
      ]
    }
  ],
  "rules": {
    "calculations": [
      {
        "field": "campo_calculado",
        "formula": "expressão matemática",
        "dependencies": ["campo1", "campo2"]
      }
    ],
    "validations": [
      {
        "field": "campo",
        "rule": "regra de validação",
        "message": "Mensagem de erro"
      }
    ]
  },
  "metadata": {
    "version": "1.0",
    "author": "Nome do Autor",
    "tags": ["tag1", "tag2"]
  }
}
```

### Field Types
- **text**: Campo de texto simples
- **textarea**: Campo de texto multilinha
- **number**: Campo numérico com min/max
- **boolean**: Checkbox verdadeiro/falso
- **select**: Lista de opções predefinidas

---

## 🧪 Exemplos de Uso

### Fluxo Completo de Autenticação
```bash
# 1. Registrar usuário
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "jogador@exemplo.com",
    "password": "minhasenha123"
  }'

# 2. Fazer login
TOKEN=$(curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "jogador@exemplo.com",
    "password": "minhasenha123"
  }' | jq -r '.token')

# 3. Acessar endpoint protegido
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/auth/me
```

### Criar Template Completo
```bash
curl -X POST http://localhost:8080/api/v1/templates \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Ficha GURPS",
    "system": "GURPS 4e",
    "description": "Template para personagens GURPS 4ª edição",
    "definition": {
      "sections": [
        {
          "name": "Atributos Primários",
          "fields": [
            {
              "name": "ST",
              "type": "number",
              "min": 1,
              "max": 30,
              "default": 10,
              "help": "Força do personagem"
            },
            {
              "name": "DX",
              "type": "number", 
              "min": 1,
              "max": 30,
              "default": 10,
              "help": "Destreza do personagem"
            }
          ]
        }
      ]
    }
  }'
```

---

## 📚 Swagger Documentation

Para uma experiência interativa completa, acesse a documentação Swagger:

**URL**: [http://localhost:8080/docs/index.html](http://localhost:8080/docs/index.html)

A documentação Swagger permite:
- Testar endpoints diretamente no navegador
- Ver exemplos de request/response
- Autenticar com JWT token
- Explorar todos os modelos de dados

---

**Para dúvidas ou suporte, consulte o [README.md](README.md) ou entre em contato com a equipe de desenvolvimento.**
