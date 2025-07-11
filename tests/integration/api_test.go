package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/luizdequeiroz/rpg-backend/internal/bff"
	"github.com/luizdequeiroz/rpg-backend/pkg/db"
)

// createTestDatabase cria uma base de dados em memória para testes
func createTestDatabase(t *testing.T) *db.DB {
	database, err := db.NewDBWithDSN("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Erro ao criar base de dados: %v", err)
	}
	return database
}

// setupTestRouter configura um router Gin para testes
func setupTestRouter(database *db.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	handler := bff.NewHandler(database)

	// Adicionar endpoint de health check simples para testes
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"timestamp": "2024-01-01T00:00:00Z",
		})
	})

	// Configurar rotas
	v1 := router.Group("/api/v1")
	handler.SetupRoutes(v1)

	return router
}

func TestHealthCheckIntegration(t *testing.T) {
	// Criar banco em memória
	testDB := createTestDatabase(t)
	defer testDB.Close()

	// Configurar router
	router := setupTestRouter(testDB)

	t.Run("Health check endpoint", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ok", response["status"])
		assert.Contains(t, response, "timestamp")
	})
}

func TestDatabaseConnection(t *testing.T) {
	// Criar banco em memória
	testDB := createTestDatabase(t)
	defer testDB.Close()

	t.Run("Conexão com base de dados", func(t *testing.T) {
		// Teste simples de ping à base de dados
		err := testDB.DB.DB.Ping()
		assert.NoError(t, err)
	})
}

func TestRouterSetup(t *testing.T) {
	// Criar banco em memória
	testDB := createTestDatabase(t)
	defer testDB.Close()

	t.Run("Router setup sem erros", func(t *testing.T) {
		router := setupTestRouter(testDB)
		assert.NotNil(t, router)
	})
}
