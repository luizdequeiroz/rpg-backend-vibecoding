# Guia de Desenvolvimento - RPG Backend

## 🛠️ Configuração do Ambiente de Desenvolvimento

### Pré-requisitos
- **Go 1.21+**: [Download](https://golang.org/dl/)
- **Git**: [Download](https://git-scm.com/)
- **VS Code** (recomendado): [Download](https://code.visualstudio.com/)
- **Make** (opcional): [Windows](http://gnuwin32.sourceforge.net/packages/make.htm)

### Extensões VS Code Recomendadas
```json
{
  "recommendations": [
    "golang.go",
    "ms-vscode.vscode-json",
    "humao.rest-client",
    "42crunch.vscode-openapi",
    "bradlc.vscode-tailwindcss"
  ]
}
```

### Configuração Inicial
```bash
# 1. Clone o repositório
git clone https://github.com/luizdequeiroz/rpg-backend.git
cd rpg-backend

# 2. Instale dependências
go mod tidy

# 3. Instale ferramentas de desenvolvimento
go install github.com/swaggo/swag/cmd/swag@latest
go install github.com/pressly/goose/v3/cmd/goose@latest

# 4. Configure variáveis de ambiente
cp .env.example .env  # Edite conforme necessário
```

## 🚀 Comandos de Desenvolvimento

### Scripts Make (Linux/Mac)
```bash
# Desenvolvimento
make run              # Executar aplicação
make dev              # Executar com hot reload
make build            # Compilar aplicação
make test             # Executar todos os testes
make test-coverage    # Testes com cobertura

# Banco de Dados
make migrate-up       # Executar migrações
make migrate-down     # Desfazer última migração
make migrate-create   # Criar nova migração
make migrate-status   # Status das migrações
make migrate-reset    # Resetar todas as migrações

# Documentação
make swagger-generate # Gerar docs Swagger
make docs-serve      # Servir documentação

# Qualidade de Código
make lint            # Executar linter
make fmt             # Formatar código
make vet             # Análise estática
make clean           # Limpar arquivos temporários
```

### Scripts Windows (scripts.bat)
```batch
# Desenvolvimento
scripts.bat run              # Executar aplicação
scripts.bat dev              # Executar com debug
scripts.bat build            # Compilar aplicação
scripts.bat test             # Executar testes

# Banco de Dados
scripts.bat migrate-up       # Executar migrações
scripts.bat migrate-down     # Desfazer migração
scripts.bat migrate-create   # Criar nova migração

# Documentação
scripts.bat swagger-generate # Gerar docs Swagger
```

## 📝 Padrões de Código

### 1. **Estrutura de Arquivos**
```go
// Ordem de imports
package main

import (
    // 1. Standard library
    "fmt"
    "net/http"
    
    // 2. Third party
    "github.com/gin-gonic/gin"
    
    // 3. Local imports
    "github.com/luizdequeiroz/rpg-backend/internal/app/models"
)
```

### 2. **Naming Conventions**
```go
// Variáveis: camelCase
var userID int
var emailAddress string

// Constantes: UPPER_CASE ou camelCase
const MaxRetries = 3
const defaultTimeout = 30 * time.Second

// Funções públicas: PascalCase
func CreateUser() {}

// Funções privadas: camelCase
func validateEmail() {}

// Interfaces: PascalCase + "er" suffix
type UserCreator interface {
    CreateUser() error
}

// Structs: PascalCase
type UserHandler struct {}
```

### 3. **Error Handling**
```go
// Definir erros específicos
var (
    ErrUserNotFound  = errors.New("usuário não encontrado")
    ErrInvalidEmail  = errors.New("email inválido")
)

// Wrapping de erros
func (s *UserService) GetUser(id int) (*User, error) {
    user, err := s.repo.GetByID(id)
    if err != nil {
        return nil, fmt.Errorf("falha ao buscar usuário %d: %w", id, err)
    }
    return user, nil
}

// Tratamento em handlers
func (h *UserHandler) GetUser(c *gin.Context) {
    user, err := h.service.GetUser(id)
    if err != nil {
        if errors.Is(err, ErrUserNotFound) {
            c.JSON(http.StatusNotFound, models.ErrorResponse{
                Error: "Usuário não encontrado",
            })
            return
        }
        
        c.JSON(http.StatusInternalServerError, models.ErrorResponse{
            Error: "Erro interno do servidor",
        })
        return
    }
    
    c.JSON(http.StatusOK, user)
}
```

### 4. **Logging**
```go
import "log/slog"

// Logging estruturado
slog.Info("User created", 
    "user_id", user.ID,
    "email", user.Email,
)

slog.Error("Database connection failed",
    "error", err,
    "database_url", dbURL,
)
```

### 5. **Validation**
```go
// Tags de validação nos models
type CreateUserRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
    Name     string `json:"name" validate:"required,min=2,max=100"`
}

// Validação customizada
func (r *CreateUserRequest) Validate() error {
    if len(r.Password) < 6 {
        return errors.New("senha deve ter pelo menos 6 caracteres")
    }
    return nil
}
```

## 🧪 Testes

### 1. **Estrutura de Testes**
```
├── internal/
│   ├── app/
│   │   ├── services/
│   │   │   ├── auth.go
│   │   │   └── auth_test.go      # Unit tests
│   │   ├── repositories/
│   │   │   ├── user.go
│   │   │   └── user_test.go      # Integration tests
│   └── bff/
│       ├── auth.go
│       └── auth_test.go          # Handler tests
├── tests/
│   ├── integration/              # End-to-end tests
│   └── fixtures/                 # Test data
```

### 2. **Unit Tests (Services)**
```go
func TestAuthService_Login(t *testing.T) {
    tests := []struct {
        name        string
        email       string
        password    string
        expectError bool
    }{
        {
            name:        "valid login",
            email:       "user@test.com",
            password:    "password123",
            expectError: false,
        },
        {
            name:        "invalid email",
            email:       "invalid-email",
            password:    "password123",
            expectError: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            service := NewAuthService(mockDB)
            
            token, user, err := service.Login(models.UserLoginRequest{
                Email:    tt.email,
                Password: tt.password,
            })
            
            if tt.expectError {
                assert.Error(t, err)
                assert.Empty(t, token)
                assert.Nil(t, user)
            } else {
                assert.NoError(t, err)
                assert.NotEmpty(t, token)
                assert.NotNil(t, user)
            }
        })
    }
}
```

### 3. **Integration Tests (Repositories)**
```go
func TestUserRepository_Create(t *testing.T) {
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    repo := NewUserRepository(db)
    
    user := &models.User{
        Email:        "test@example.com",
        PasswordHash: "hashed_password",
    }
    
    createdUser, err := repo.Create(user)
    
    assert.NoError(t, err)
    assert.NotZero(t, createdUser.ID)
    assert.Equal(t, user.Email, createdUser.Email)
}
```

### 4. **Handler Tests**
```go
func TestAuthHandler_Login(t *testing.T) {
    gin.SetMode(gin.TestMode)
    
    mockService := &MockAuthService{}
    handler := NewAuthHandler(mockService)
    
    router := gin.New()
    router.POST("/login", handler.Login)
    
    payload := `{"email":"test@example.com","password":"password123"}`
    req := httptest.NewRequest("POST", "/login", strings.NewReader(payload))
    req.Header.Set("Content-Type", "application/json")
    
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
    
    var response models.AuthResponse
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.NotEmpty(t, response.Token)
}
```

### 5. **Test Helpers**
```go
// tests/helpers/database.go
func SetupTestDB(t *testing.T) *db.DB {
    testDB, err := db.NewDB(":memory:")
    require.NoError(t, err)
    
    // Executar migrações
    err = testDB.RunMigrations()
    require.NoError(t, err)
    
    return testDB
}

func CleanupTestDB(t *testing.T, db *db.DB) {
    err := db.Close()
    require.NoError(t, err)
}

// tests/helpers/auth.go
func CreateTestUser(t *testing.T, db *db.DB) *models.User {
    user := &models.User{
        Email:        "test@example.com",
        PasswordHash: "hashed_password",
    }
    
    repo := repositories.NewUserRepository(db)
    createdUser, err := repo.Create(user)
    require.NoError(t, err)
    
    return createdUser
}
```

## 📊 Monitoramento e Observabilidade

### 1. **Logging**
```go
// Configuração de logging
func setupLogging() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    }))
    slog.SetDefault(logger)
}

// Logging em handlers
func (h *UserHandler) CreateUser(c *gin.Context) {
    slog.Info("Creating user", 
        "request_id", c.GetHeader("X-Request-ID"),
        "user_agent", c.GetHeader("User-Agent"),
    )
    
    // ... lógica do handler
    
    slog.Info("User created successfully",
        "user_id", user.ID,
        "email", user.Email,
    )
}
```

### 2. **Métricas**
```go
// pkg/metrics/metrics.go
var (
    requestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    requestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request duration in seconds",
        },
        []string{"method", "endpoint"},
    )
)

// Middleware de métricas
func MetricsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        c.Next()
        
        duration := time.Since(start).Seconds()
        requestDuration.WithLabelValues(
            c.Request.Method,
            c.FullPath(),
        ).Observe(duration)
        
        requestsTotal.WithLabelValues(
            c.Request.Method,
            c.FullPath(),
            strconv.Itoa(c.Writer.Status()),
        ).Inc()
    }
}
```

### 3. **Health Checks**
```go
// internal/bff/health.go
type HealthResponse struct {
    Status    string                 `json:"status"`
    Timestamp time.Time              `json:"timestamp"`
    Version   string                 `json:"version"`
    Services  map[string]interface{} `json:"services"`
}

func (h *Handler) HealthHandler(c *gin.Context) {
    services := make(map[string]interface{})
    
    // Check database
    if err := h.db.Health(); err != nil {
        services["database"] = map[string]interface{}{
            "status": "unhealthy",
            "error":  err.Error(),
        }
    } else {
        services["database"] = "ok"
    }
    
    // Determinar status geral
    status := "ok"
    for _, service := range services {
        if serviceMap, ok := service.(map[string]interface{}); ok {
            if serviceMap["status"] == "unhealthy" {
                status = "unhealthy"
                break
            }
        }
    }
    
    response := HealthResponse{
        Status:    status,
        Timestamp: time.Now(),
        Version:   "1.0.0",
        Services:  services,
    }
    
    statusCode := http.StatusOK
    if status == "unhealthy" {
        statusCode = http.StatusServiceUnavailable
    }
    
    c.JSON(statusCode, response)
}
```

## 🔧 Debugging

### 1. **VS Code Debug Configuration**
```json
// .vscode/launch.json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch API",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/cmd/api",
            "env": {
                "LOG_LEVEL": "debug",
                "DATABASE_URL": "file:./data/test.db"
            },
            "args": []
        },
        {
            "name": "Debug Test",
            "type": "go",
            "request": "launch",
            "mode": "test",
            "program": "${workspaceFolder}",
            "env": {
                "DATABASE_URL": ":memory:"
            }
        }
    ]
}
```

### 2. **Profiling**
```go
// cmd/api/main.go
import _ "net/http/pprof"

func main() {
    // ... configuração normal
    
    // Endpoint de profiling (apenas em desenvolvimento)
    if gin.Mode() == gin.DebugMode {
        go func() {
            log.Println("Profiling server running on :6060")
            log.Println(http.ListenAndServe(":6060", nil))
        }()
    }
    
    // ... resto da aplicação
}
```

### 3. **Debug Utilities**
```go
// pkg/debug/debug.go
func PrintJSON(v interface{}) {
    b, _ := json.MarshalIndent(v, "", "  ")
    fmt.Printf("%s\n", b)
}

func PrintStruct(v interface{}) {
    fmt.Printf("%+v\n", v)
}

// Middleware de debug
func DebugMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        if gin.Mode() == gin.DebugMode {
            fmt.Printf("Request: %s %s\n", c.Request.Method, c.Request.URL.Path)
            for k, v := range c.Request.Header {
                fmt.Printf("Header: %s = %v\n", k, v)
            }
        }
        c.Next()
    }
}
```

## 📋 Checklist de Pull Request

### ✅ Código
- [ ] Código segue as convenções do projeto
- [ ] Funções têm documentação adequada
- [ ] Tratamento de erros implementado
- [ ] Logs apropriados adicionados
- [ ] Código não contém TODOs ou FIXMEs

### ✅ Testes
- [ ] Unit tests para novas funções
- [ ] Integration tests para mudanças de banco
- [ ] Testes passando localmente
- [ ] Cobertura de testes mantida/melhorada

### ✅ Documentação
- [ ] README atualizado se necessário
- [ ] Swagger docs atualizadas
- [ ] Comentários de código adequados
- [ ] Exemplos de uso fornecidos

### ✅ Performance
- [ ] Queries de banco otimizadas
- [ ] Memory leaks verificados
- [ ] Timeouts apropriados configurados

### ✅ Segurança
- [ ] Inputs validados
- [ ] Autorização verificada
- [ ] Dados sensíveis protegidos
- [ ] SQL injection prevenido

---

**Happy Coding! 🚀**
