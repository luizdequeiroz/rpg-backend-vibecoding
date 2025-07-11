package bff

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthResponse representa a resposta do healthcheck
type HealthResponse struct {
	Status    string            `json:"status" example:"ok"`
	Timestamp time.Time         `json:"timestamp" example:"2025-07-10T15:47:08Z"`
	Version   string            `json:"version" example:"1.0.0"`
	Services  map[string]string `json:"services" example:"database:ok"`
}

// HealthHandler implementa o endpoint de healthcheck
// @Summary Verificação de saúde da API
// @Description Verifica o status de saúde da API e seus serviços dependentes
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse "API funcionando normalmente"
// @Success 503 {object} HealthResponse "API com problemas em algum serviço"
// @Router /health [get]
func (h *Handler) HealthHandler(c *gin.Context) {
	// Verificar saúde do banco de dados
	dbStatus := "ok"
	if err := h.db.Health(); err != nil {
		dbStatus = "error: " + err.Error()
	}

	response := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now(),
		Version:   "1.0.0", // TODO: obter da build ou config
		Services: map[string]string{
			"database": dbStatus,
		},
	}

	// Se algum serviço está com problema, marcar como degraded
	statusCode := http.StatusOK
	for _, status := range response.Services {
		if status != "ok" {
			response.Status = "degraded"
			statusCode = http.StatusServiceUnavailable
			break
		}
	}

	c.JSON(statusCode, response)
}
