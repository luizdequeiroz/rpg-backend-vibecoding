package bff

import (
	"github.com/gin-gonic/gin"
	"github.com/luizdequeiroz/rpg-backend/internal/app/handlers"
	"github.com/luizdequeiroz/rpg-backend/internal/app/middleware"
	"github.com/luizdequeiroz/rpg-backend/internal/app/services"
)

// DiceHandler wrapper para o sistema de dados
type DiceHandler struct {
	diceHandler *handlers.DiceHandler
}

// NewDiceHandler cria novo handler para dados
func NewDiceHandler(diceService *services.DiceService, playerSheetService *services.PlayerSheetService) *DiceHandler {
	return &DiceHandler{
		diceHandler: handlers.NewDiceHandler(diceService, playerSheetService),
	}
}

// SetupDiceRoutes configura rotas de dados
func (h *DiceHandler) SetupDiceRoutes(router *gin.RouterGroup, authService *services.AuthService) {
	dice := router.Group("/dice")
	dice.Use(middleware.AuthMiddleware(authService))
	{
		dice.POST("/roll", h.diceHandler.RollDice)
		dice.POST("/roll-with-sheet", h.diceHandler.RollWithSheet)
		dice.GET("/history", h.diceHandler.GetHistory)
	}
}

// Métodos auxiliares para integração com o BFF se necessário
func (h *DiceHandler) RollDice(c *gin.Context) {
	h.diceHandler.RollDice(c)
}

func (h *DiceHandler) RollWithSheet(c *gin.Context) {
	h.diceHandler.RollWithSheet(c)
}

func (h *DiceHandler) GetHistory(c *gin.Context) {
	h.diceHandler.GetHistory(c)
}
