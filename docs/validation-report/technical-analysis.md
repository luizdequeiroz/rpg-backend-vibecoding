# Análise Técnica da Arquitetura

## 1. Estrutura de Arquitetura

### 1.1 Clean Architecture Implementation
O sistema segue os princípios de Clean Architecture com separação clara de responsabilidades:

```
Presentation Layer (BFF)
├── HTTP Handlers
├── Request/Response DTOs
└── Route Configuration

Application Layer (Services)
├── Business Logic
├── Validation Rules
└── Orchestration

Domain Layer (Models)
├── Entity Definitions
├── Value Objects
└── Business Rules

Infrastructure Layer (Repositories)
├── Database Access
├── External Services
└── Data Persistence
```

### 1.2 Dependency Injection
- Configuração centralizada no main.go
- Inversão de dependências implementada
- Facilita testes unitários e mocking

## 2. Modelos de Dados

### 2.1 Entidades Principais
```go
// Usuários e Autenticação
type User struct {
    ID        int       `json:"id" db:"id"`
    Email     string    `json:"email" db:"email"`
    Password  string    `json:"-" db:"password_hash"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Mesas de Jogo
type GameTable struct {
    ID        string    `json:"id" db:"id"`
    Name      string    `json:"name" db:"name"`
    System    string    `json:"system" db:"system"`
    OwnerID   int       `json:"owner_id" db:"owner_id"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Fichas de Personagem
type PlayerSheet struct {
    ID         string    `json:"id" db:"id"`
    TableID    string    `json:"table_id" db:"table_id"`
    TemplateID int       `json:"template_id" db:"template_id"`
    OwnerID    int       `json:"owner_id" db:"owner_id"`
    Name       string    `json:"name" db:"name"`
    Data       string    `json:"-" db:"data"` // JSON como string
    CreatedAt  time.Time `json:"created_at" db:"created_at"`
    UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// Rolagens de Dados
type Roll struct {
    ID            int        `json:"id" db:"id"`
    SheetID       *int       `json:"sheet_id" db:"sheet_id"`
    TableID       *string    `json:"table_id" db:"table_id"`
    UserID        int        `json:"user_id" db:"user_id"`
    Expression    string     `json:"expression" db:"expression"`
    FieldName     *string    `json:"field_name" db:"field_name"`
    ResultValue   int        `json:"result_value" db:"result_value"`
    ResultDetails string     `json:"-" db:"result_details"` // JSON
    Comment       string     `json:"comment" db:"comment"`
    Success       *bool      `json:"success" db:"success"`
    CreatedAt     time.Time  `json:"created_at" db:"created_at"`
}
```

### 2.2 Relacionamentos
- **Users** ← 1:N → **GameTables** (via owner_id)
- **Users** ← 1:N → **PlayerSheets** (via owner_id)
- **GameTables** ← 1:N → **PlayerSheets** (via table_id)
- **PlayerSheets** ← 1:N → **Rolls** (via sheet_id)
- **Users** ← 1:N → **Rolls** (via user_id)

## 3. Sistema de Rolagem de Dados

### 3.1 Engine de Parsing
```go
type DiceExpression struct {
    Count    int // Número de dados
    Sides    int // Número de lados
    Modifier int // Modificador (+/-)
}

// Regex para parsing: "^(\d+)d(\d+)([\+\-]\d+)?$"
// Exemplos: "1d20+3", "2d6-1", "3d8"
```

### 3.2 Funcionalidades Avançadas
- **Validação de Limites:** 1-100 dados, 2-1000 lados
- **Detecção de Críticos:** d20 = 20 (crítico), d20 = 1 (fumble)
- **Integração com Fichas:** Expressões como "1d20+{strength}"
- **Histórico Persistente:** Todas as rolagens salvas no banco

### 3.3 Algoritmo de Geração
```go
func (s *DiceService) RollDice(expression string, userID int) (*models.DiceRollResponse, error) {
    // 1. Parse da expressão
    dice, err := s.ParseDiceExpression(expression)
    
    // 2. Validação
    if dice.Count > 100 || dice.Sides > 1000 { return error }
    
    // 3. Geração de números aleatórios
    rand.Seed(time.Now().UnixNano())
    for i := 0; i < dice.Count; i++ {
        roll := rand.Intn(dice.Sides) + 1
    }
    
    // 4. Cálculo do resultado final
    finalResult := total + dice.Modifier
    
    // 5. Detecção de críticos/fumbles
    // 6. Persistência no banco
    // 7. Retorno da resposta
}
```

## 4. Sistema de Autenticação

### 4.1 JWT Implementation
```go
type Claims struct {
    UserID int    `json:"user_id"`
    Email  string `json:"email"`
    jwt.StandardClaims
}

// Token válido por 24 horas
expirationTime := time.Now().Add(24 * time.Hour)
```

### 4.2 Middleware de Autenticação
```go
func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Extração do token do header Authorization
        // 2. Validação do formato "Bearer <token>"
        // 3. Verificação da assinatura JWT
        // 4. Extração dos claims
        // 5. Injeção do userID no contexto
        c.Set("userID", claims.UserID)
        c.Next()
    }
}
```

## 5. Banco de Dados

### 5.1 Esquema Relacional
```sql
-- Usuários
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Mesas de Jogo
CREATE TABLE game_tables (
    id TEXT PRIMARY KEY, -- UUID
    name TEXT NOT NULL,
    system TEXT NOT NULL,
    owner_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (owner_id) REFERENCES users(id)
);

-- Fichas de Personagem
CREATE TABLE player_sheets (
    id TEXT PRIMARY KEY, -- UUID
    table_id TEXT NOT NULL,
    template_id INTEGER NOT NULL,
    owner_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    data TEXT, -- JSON
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (table_id) REFERENCES game_tables(id),
    FOREIGN KEY (owner_id) REFERENCES users(id),
    FOREIGN KEY (template_id) REFERENCES sheet_templates(id)
);

-- Rolagens
CREATE TABLE rolls (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sheet_id INTEGER,
    table_id TEXT,
    user_id INTEGER NOT NULL,
    expression TEXT NOT NULL,
    field_name TEXT,
    result_value INTEGER NOT NULL,
    result_details TEXT, -- JSON
    comment TEXT,
    success BOOLEAN,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sheet_id) REFERENCES player_sheets(rowid),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

### 5.2 Índices para Performance
```sql
-- Índices críticos para consultas frequentes
CREATE INDEX idx_player_sheets_table_id ON player_sheets(table_id);
CREATE INDEX idx_player_sheets_owner_id ON player_sheets(owner_id);
CREATE INDEX idx_rolls_user_id ON rolls(user_id);
CREATE INDEX idx_rolls_sheet_id ON rolls(sheet_id);
CREATE INDEX idx_rolls_created_at ON rolls(created_at);
```

## 6. API Design

### 6.1 Padrões RESTful
```
GET    /api/v1/resource       # Listar
POST   /api/v1/resource       # Criar
GET    /api/v1/resource/{id}  # Buscar
PUT    /api/v1/resource/{id}  # Atualizar
DELETE /api/v1/resource/{id}  # Deletar
```

### 6.2 Códigos de Status HTTP
- **200 OK:** Operação bem-sucedida
- **201 Created:** Recurso criado com sucesso
- **400 Bad Request:** Dados inválidos
- **401 Unauthorized:** Token inválido/ausente
- **403 Forbidden:** Sem permissão
- **404 Not Found:** Recurso não encontrado
- **500 Internal Server Error:** Erro interno

### 6.3 Estrutura de Resposta
```json
// Sucesso
{
  "data": { ... },
  "meta": {
    "page": 1,
    "limit": 20,
    "total": 100
  }
}

// Erro
{
  "error": "Mensagem do erro",
  "message": "Detalhes adicionais",
  "code": "ERROR_CODE"
}
```

## 7. Segurança

### 7.1 Autenticação e Autorização
- **Hash de Senhas:** bcrypt com salt automático
- **JWT:** Assinatura HMAC-SHA256
- **Middleware:** Validação automática em rotas protegidas
- **UUIDs:** Para evitar enumeration attacks

### 7.2 Validação de Dados
```go
// Validação de entrada
type CreatePlayerSheetRequest struct {
    TableID    string          `json:"table_id" validate:"required,uuid"`
    TemplateID int             `json:"template_id" validate:"required,min=1"`
    Name       string          `json:"name" validate:"required,min=3,max=100"`
    Data       PlayerSheetData `json:"data" validate:"required"`
}
```

## 8. Performance e Escalabilidade

### 8.1 Otimizações Implementadas
- **Connection Pool:** sqlx com pool de conexões
- **Índices:** Criados para consultas críticas
- **Paginação:** Implementada em listagens
- **Prepared Statements:** Para consultas repetitivas

### 8.2 Pontos de Melhoria
- **Cache:** Redis para sessões e dados frequentes
- **Rate Limiting:** Para proteger contra abuso
- **Background Jobs:** Para processamento assíncrono
- **Load Balancing:** Para múltiplas instâncias

## 9. Monitoramento e Observabilidade

### 9.1 Logging
```go
// Gin middleware para logs estruturados
[GIN] 2025/07/11 - 14:46:52 | 201 | 25.9045ms | 127.0.0.1 | POST "/api/v1/tables/"
```

### 9.2 Health Check
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

## 10. Análise de Qualidade

### 10.1 Pontos Fortes
- ✅ Arquitetura limpa e bem estruturada
- ✅ Separação clara de responsabilidades
- ✅ Código legível e maintível
- ✅ Padrões consistentes
- ✅ Documentação abrangente

### 10.2 Áreas de Melhoria
- 🔄 Testes automatizados (unitários e integração)
- 🔄 Monitoring e alertas
- 🔄 Performance profiling
- 🔄 Deploy automatizado
- 🔄 Backup e disaster recovery

### 10.3 Métricas de Código
- **Complexidade:** Baixa a Média
- **Acoplamento:** Baixo
- **Coesão:** Alta
- **Testabilidade:** Boa (devido à DI)
- **Maintibilidade:** Excelente
