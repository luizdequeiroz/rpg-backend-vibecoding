package bff

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/luizdequeiroz/rpg-backend/internal/app/middleware"
	"github.com/luizdequeiroz/rpg-backend/internal/app/models"
	"github.com/luizdequeiroz/rpg-backend/internal/app/services"
)

// PlayerSheetHandler gerencia requisições HTTP para fichas
type PlayerSheetHandler struct {
	sheetService *services.PlayerSheetService
}

// NewPlayerSheetHandler cria novo handler
func NewPlayerSheetHandler(sheetService *services.PlayerSheetService) *PlayerSheetHandler {
	return &PlayerSheetHandler{
		sheetService: sheetService,
	}
}

// CreateSheet cria nova ficha
// @Summary Criar ficha de personagem
// @Description Cria nova ficha de personagem em uma mesa
// @Tags Player Sheets
// @Accept json
// @Produce json
// @Param body body models.CreatePlayerSheetRequest true "Dados da ficha"
// @Security BearerAuth
// @Success 201 {object} models.PlayerSheetResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/sheets [post]
func (h *PlayerSheetHandler) CreateSheet(c *gin.Context) {
	userID, _, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "ID do usuário não encontrado no contexto",
		})
		return
	}

	var req models.CreatePlayerSheetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Dados inválidos",
			"details": err.Error(),
		})
		return
	}

	// O tableID agora vem no corpo da requisição
	sheet, err := h.sheetService.Create(req, req.TableID, userID)
	if err != nil {
		if err.Error() == "acesso negado à mesa" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao criar ficha",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, sheet)
}

// GetSheetsByTable lista fichas da mesa
// @Summary Listar fichas da mesa
// @Description Lista todas as fichas de personagens de uma mesa
// @Tags Player Sheets
// @Produce json
// @Param table_id query string true "ID da mesa"
// @Param page query int false "Página" default(1)
// @Param limit query int false "Itens por página" default(20)
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/sheets [get]
func (h *PlayerSheetHandler) GetSheetsByTable(c *gin.Context) {
	tableID := c.Query("table_id")
	if tableID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "table_id é obrigatório",
		})
		return
	}

	userID, _, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	sheets, err := h.sheetService.GetByTableID(tableID, userID, page, limit)
	if err != nil {
		if err.Error() == "acesso negado à mesa" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao buscar fichas",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sheets": sheets,
		"total":  len(sheets),
		"page":   page,
		"limit":  limit,
	})
}

// GetSheet busca ficha por ID
// @Summary Buscar ficha
// @Description Busca ficha de personagem por ID
// @Tags Player Sheets
// @Produce json
// @Param id path string true "ID da ficha"
// @Security BearerAuth
// @Success 200 {object} models.PlayerSheetResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/sheets/{id} [get]
func (h *PlayerSheetHandler) GetSheet(c *gin.Context) {
	sheetID := c.Param("id")
	userID, _, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	sheet, err := h.sheetService.GetByID(sheetID, userID)
	if err != nil {
		if err.Error() == "ficha não encontrada" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "acesso negado" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao buscar ficha",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, sheet)
}

// UpdateSheet atualiza ficha
// @Summary Atualizar ficha
// @Description Atualiza dados de uma ficha de personagem
// @Tags Player Sheets
// @Accept json
// @Produce json
// @Param id path string true "ID da ficha"
// @Param body body models.UpdatePlayerSheetRequest true "Dados para atualização"
// @Security BearerAuth
// @Success 200 {object} models.PlayerSheetResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/sheets/{id} [put]
func (h *PlayerSheetHandler) UpdateSheet(c *gin.Context) {
	sheetID := c.Param("id")
	userID, _, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	var req models.UpdatePlayerSheetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Dados inválidos",
			"details": err.Error(),
		})
		return
	}

	sheet, err := h.sheetService.Update(sheetID, req, userID)
	if err != nil {
		if err.Error() == "ficha não encontrada" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "apenas o proprietário pode atualizar a ficha" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao atualizar ficha",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, sheet)
}

// DeleteSheet remove ficha
// @Summary Remover ficha
// @Description Remove uma ficha de personagem
// @Tags Player Sheets
// @Produce json
// @Param id path string true "ID da ficha"
// @Security BearerAuth
// @Success 204
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/sheets/{id} [delete]
func (h *PlayerSheetHandler) DeleteSheet(c *gin.Context) {
	sheetID := c.Param("id")
	userID, _, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	err := h.sheetService.Delete(sheetID, userID)
	if err != nil {
		if err.Error() == "apenas o proprietário da ficha ou da mesa pode removê-la" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao remover ficha",
			"details": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// RollDice executa rolagem de dados
// @Summary Rolar dados
// @Description Executa rolagem de dados baseada em expressão ou campo da ficha
// @Tags Player Sheets
// @Accept json
// @Produce json
// @Param body body models.CreateRollRequest true "Dados da rolagem"
// @Security BearerAuth
// @Success 200 {object} models.RollResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/rolls [post]
func (h *PlayerSheetHandler) RollDice(c *gin.Context) {
	userID, _, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	var req models.CreateRollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Dados inválidos",
			"details": err.Error(),
		})
		return
	}

	if req.SheetID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "sheet_id é obrigatório",
		})
		return
	}

	roll, err := h.sheetService.CreateRoll(req.SheetID, req, userID)
	if err != nil {
		if err.Error() == "ficha não encontrada" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "acesso negado à mesa" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Erro na rolagem",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, roll)
}

// GetRollsByTable lista rolagens da mesa
// @Summary Listar rolagens da mesa
// @Description Lista histórico de rolagens de uma mesa
// @Tags Player Sheets
// @Produce json
// @Param tableID path string true "ID da mesa"
// @Param page query int false "Página" default(1)
// @Param limit query int false "Itens por página" default(20)
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/rolls/table/{tableID} [get]
func (h *PlayerSheetHandler) GetRollsByTable(c *gin.Context) {
	tableID := c.Param("tableID")
	userID := c.GetInt("userID")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	rolls, err := h.sheetService.GetRollsByTableID(tableID, userID, page, limit)
	if err != nil {
		if err.Error() == "acesso negado à mesa" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao buscar rolagens",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rolls": rolls,
		"total": len(rolls),
		"page":  page,
		"limit": limit,
	})
}

// GetRollsBySheet lista rolagens da ficha
// @Summary Listar rolagens da ficha
// @Description Lista histórico de rolagens de uma ficha específica
// @Tags Player Sheets
// @Produce json
// @Param sheetID path string true "ID da ficha"
// @Param page query int false "Página" default(1)
// @Param limit query int false "Itens por página" default(20)
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/rolls/sheet/{sheetID} [get]
func (h *PlayerSheetHandler) GetRollsBySheet(c *gin.Context) {
	sheetID := c.Param("sheetID")
	userID := c.GetInt("userID")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	rolls, err := h.sheetService.GetRollsBySheetID(sheetID, userID, page, limit)
	if err != nil {
		if err.Error() == "acesso negado à mesa" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao buscar rolagens",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rolls": rolls,
		"total": len(rolls),
		"page":  page,
		"limit": limit,
	})
}
