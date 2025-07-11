package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/luizdequeiroz/rpg-backend/internal/app/interfaces"
	"github.com/luizdequeiroz/rpg-backend/internal/app/models"
	"github.com/luizdequeiroz/rpg-backend/internal/app/services"
)

type DiceHandler struct {
	diceService        *services.DiceService
	playerSheetService *services.PlayerSheetService
	notificationService interfaces.NotificationService
}

func NewDiceHandler(diceService *services.DiceService, playerSheetService *services.PlayerSheetService, notificationService interfaces.NotificationService) *DiceHandler {
	return &DiceHandler{
		diceService:        diceService,
		playerSheetService: playerSheetService,
		notificationService: notificationService,
	}
}

// RollDice executa uma rolagem de dados
// @Summary Rolar dados
// @Description Executa uma rolagem de dados com expressão personalizada
// @Tags dice
// @Accept json
// @Produce json
// @Param request body models.DiceRollRequest true "Dados da rolagem"
// @Security ApiKeyAuth
// @Success 200 {object} models.DiceRollResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/dice/roll [post]
func (h *DiceHandler) RollDice(c *gin.Context) {
	var req models.DiceRollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Dados inválidos",
			Message: err.Error(),
		})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "Não autorizado",
			Message: "Token inválido",
		})
		return
	}

	result, err := h.diceService.RollDice(req.Expression, userID.(int))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Erro na rolagem",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// RollWithSheet executa uma rolagem usando dados da ficha do personagem
// @Summary Rolar dados com ficha
// @Description Executa uma rolagem usando atributos da ficha do personagem
// @Tags dice
// @Accept json
// @Produce json
// @Param request body models.DiceRollWithSheetRequest true "Dados da rolagem com ficha"
// @Security ApiKeyAuth
// @Success 200 {object} models.DiceRollResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/dice/roll-with-sheet [post]
func (h *DiceHandler) RollWithSheet(c *gin.Context) {
	var req models.DiceRollWithSheetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Dados inválidos",
			Message: err.Error(),
		})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "Não autorizado",
			Message: "Token inválido",
		})
		return
	}

	// Verificar se a ficha existe e pertence ao usuário
	sheet, err := h.playerSheetService.GetByID(req.SheetID, userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "Ficha não encontrada",
			Message: "Ficha não existe ou não pertence ao usuário",
		})
		return
	}

	result, err := h.diceService.RollWithSheet(req.Expression, req.AttributeField, sheet)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Erro na rolagem",
			Message: err.Error(),
		})
		return
	}

	// Notificar via WebSocket se há serviço de notificação configurado
	if h.notificationService != nil {
		userEmail, _ := c.Get("userEmail")
		h.notificationService.NotifyRollPerformed(
			sheet.TableID, 
			userID.(int), 
			userEmail.(string), 
			result,
		)
	}

	c.JSON(http.StatusOK, result)
}

// GetHistory recupera histórico de rolagens
// @Summary Histórico de rolagens
// @Description Recupera histórico de rolagens do usuário
// @Tags dice
// @Produce json
// @Param page query int false "Página" default(1)
// @Param limit query int false "Limite por página" default(10)
// @Security ApiKeyAuth
// @Success 200 {object} models.DiceHistoryResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/dice/history [get]
func (h *DiceHandler) GetHistory(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "Não autorizado",
			Message: "Token inválido",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	rolls, total, err := h.diceService.GetUserHistory(userID.(int), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Erro interno",
			Message: "Não foi possível recuperar o histórico",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Histórico de rolagens",
		"rolls":   rolls,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

