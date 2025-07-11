package bff

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/luizdequeiroz/rpg-backend/internal/app/middleware"
	"github.com/luizdequeiroz/rpg-backend/internal/app/models"
	"github.com/luizdequeiroz/rpg-backend/internal/app/services"
)

// GameTableHandler gerencia endpoints de mesas de jogo
type GameTableHandler struct {
	service *services.GameTableService
}

// NewGameTableHandler cria uma nova instância do handler
func NewGameTableHandler(service *services.GameTableService) *GameTableHandler {
	return &GameTableHandler{
		service: service,
	}
}

// SetupGameTableRoutes configura as rotas de mesas de jogo
func (h *GameTableHandler) SetupGameTableRoutes(router *gin.RouterGroup, authService *services.AuthService) {
	tables := router.Group("/tables")

	// Rotas de mesas (todas requerem autenticação)
	tables.Use(middleware.AuthMiddleware(authService))
	{
		tables.POST("/", h.CreateTable)
		tables.GET("/", h.ListTables)
		tables.GET("/:id", h.GetTable)
		tables.PUT("/:id", h.UpdateTable)
		tables.DELETE("/:id", h.DeleteTable)

		// Rotas de convites
		invites := tables.Group("/:id/invites")
		{
			invites.POST("/", h.CreateInvite)
			invites.GET("/", h.ListInvites)
			invites.POST("/:inviteId/accept", h.AcceptInvite)
			invites.POST("/:inviteId/decline", h.DeclineInvite)
		}
	}
}

// CreateTable godoc
// @Summary Criar nova mesa de jogo
// @Description Cria uma nova mesa de jogo. O usuário autenticado se torna o proprietário.
// @Tags GameTables
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body models.CreateGameTableRequest true "Dados da mesa"
// @Success 201 {object} models.GameTableResponse
// @Failure 400 {object} map[string]interface{} "Dados inválidos"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/v1/tables [post]
func (h *GameTableHandler) CreateTable(c *gin.Context) {
	userID, _, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não encontrado no contexto"})
		return
	}

	var req models.CreateGameTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos", "details": err.Error()})
		return
	}

	table, err := h.service.Create(req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, table)
}

// ListTables godoc
// @Summary Listar mesas do usuário
// @Description Lista todas as mesas onde o usuário é proprietário ou convidado aceito
// @Tags GameTables
// @Produce json
// @Security Bearer
// @Param page query int false "Número da página" default(1)
// @Param limit query int false "Itens por página" default(20)
// @Success 200 {object} map[string]interface{} "Lista de mesas"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/v1/tables [get]
func (h *GameTableHandler) ListTables(c *gin.Context) {
	userID, _, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não encontrado no contexto"})
		return
	}

	// Parâmetros de paginação
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	tables, err := h.service.GetTablesForUser(userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tables": tables,
		"page":   page,
		"limit":  limit,
		"total":  len(tables),
	})
}

// GetTable godoc
// @Summary Buscar mesa por ID
// @Description Retorna detalhes de uma mesa específica, incluindo convites se o usuário for o proprietário
// @Tags GameTables
// @Produce json
// @Security Bearer
// @Param id path string true "ID da mesa"
// @Success 200 {object} models.GameTableResponse
// @Failure 400 {object} map[string]interface{} "ID inválido"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 403 {object} map[string]interface{} "Acesso negado"
// @Failure 404 {object} map[string]interface{} "Mesa não encontrada"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/v1/tables/{id} [get]
func (h *GameTableHandler) GetTable(c *gin.Context) {
	userID, _, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não encontrado no contexto"})
		return
	}

	tableID := c.Param("id")
	if tableID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID da mesa é obrigatório"})
		return
	}

	table, err := h.service.GetByID(tableID, userID)
	if err != nil {
		if err.Error() == "mesa não encontrada" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "acesso negado" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, table)
}

// UpdateTable godoc
// @Summary Atualizar mesa
// @Description Atualiza nome e/ou sistema de uma mesa. Apenas o proprietário pode atualizar.
// @Tags GameTables
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "ID da mesa"
// @Param request body models.UpdateGameTableRequest true "Dados para atualização"
// @Success 200 {object} models.GameTableResponse
// @Failure 400 {object} map[string]interface{} "Dados inválidos"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 403 {object} map[string]interface{} "Apenas proprietário pode atualizar"
// @Failure 404 {object} map[string]interface{} "Mesa não encontrada"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/v1/tables/{id} [put]
func (h *GameTableHandler) UpdateTable(c *gin.Context) {
	userID, _, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não encontrado no contexto"})
		return
	}

	tableID := c.Param("id")
	if tableID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID da mesa é obrigatório"})
		return
	}

	var req models.UpdateGameTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos", "details": err.Error()})
		return
	}

	table, err := h.service.Update(tableID, req, userID)
	if err != nil {
		if err.Error() == "mesa não encontrada" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "apenas o proprietário pode atualizar a mesa" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, table)
}

// DeleteTable godoc
// @Summary Remover mesa
// @Description Remove uma mesa permanentemente. Apenas o proprietário pode remover.
// @Tags GameTables
// @Produce json
// @Security Bearer
// @Param id path string true "ID da mesa"
// @Success 204 "Mesa removida com sucesso"
// @Failure 400 {object} map[string]interface{} "ID inválido"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 403 {object} map[string]interface{} "Apenas proprietário pode remover"
// @Failure 404 {object} map[string]interface{} "Mesa não encontrada"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/v1/tables/{id} [delete]
func (h *GameTableHandler) DeleteTable(c *gin.Context) {
	userID, _, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não encontrado no contexto"})
		return
	}

	tableID := c.Param("id")
	if tableID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID da mesa é obrigatório"})
		return
	}

	err := h.service.Delete(tableID, userID)
	if err != nil {
		if err.Error() == "mesa não encontrada" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "apenas o proprietário pode remover a mesa" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// CreateInvite godoc
// @Summary Criar convite para mesa
// @Description Cria um convite para um usuário participar da mesa. Apenas o proprietário pode criar convites.
// @Tags GameTables
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "ID da mesa"
// @Param request body models.CreateInviteRequest true "Email do usuário a ser convidado"
// @Success 201 {object} models.InviteDetails
// @Failure 400 {object} map[string]interface{} "Dados inválidos"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 403 {object} map[string]interface{} "Apenas proprietário pode criar convites"
// @Failure 404 {object} map[string]interface{} "Mesa ou usuário não encontrado"
// @Failure 409 {object} map[string]interface{} "Convite já existe"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/v1/tables/{id}/invites [post]
func (h *GameTableHandler) CreateInvite(c *gin.Context) {
	userID, _, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não encontrado no contexto"})
		return
	}

	tableID := c.Param("id")
	if tableID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID da mesa é obrigatório"})
		return
	}

	var req models.CreateInviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos", "details": err.Error()})
		return
	}

	invite, err := h.service.CreateInvite(tableID, req, userID)
	if err != nil {
		if err.Error() == "mesa não encontrada" || err.Error() == "usuário não encontrado" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "apenas o proprietário pode criar convites" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "convite já existe para este usuário" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, invite)
}

// ListInvites godoc
// @Summary Listar convites da mesa
// @Description Lista todos os convites de uma mesa. Acessível pelo proprietário e convidados.
// @Tags GameTables
// @Produce json
// @Security Bearer
// @Param id path string true "ID da mesa"
// @Success 200 {object} map[string]interface{} "Lista de convites"
// @Failure 400 {object} map[string]interface{} "ID inválido"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 403 {object} map[string]interface{} "Acesso negado"
// @Failure 404 {object} map[string]interface{} "Mesa não encontrada"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/v1/tables/{id}/invites [get]
func (h *GameTableHandler) ListInvites(c *gin.Context) {
	userID, _, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não encontrado no contexto"})
		return
	}

	tableID := c.Param("id")
	if tableID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID da mesa é obrigatório"})
		return
	}

	invites, err := h.service.GetInvitesForTable(tableID, userID)
	if err != nil {
		if err.Error() == "acesso negado" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"invites": invites,
		"total":   len(invites),
	})
}

// AcceptInvite godoc
// @Summary Aceitar convite
// @Description Aceita um convite para participar da mesa. Apenas o convidado pode aceitar.
// @Tags GameTables
// @Produce json
// @Security Bearer
// @Param id path string true "ID da mesa"
// @Param inviteId path string true "ID do convite"
// @Success 200 {object} map[string]interface{} "Convite aceito"
// @Failure 400 {object} map[string]interface{} "ID inválido"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 403 {object} map[string]interface{} "Apenas convidado pode aceitar"
// @Failure 404 {object} map[string]interface{} "Convite não encontrado"
// @Failure 409 {object} map[string]interface{} "Convite já respondido"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/v1/tables/{id}/invites/{inviteId}/accept [post]
func (h *GameTableHandler) AcceptInvite(c *gin.Context) {
	userID, _, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não encontrado no contexto"})
		return
	}

	inviteID := c.Param("inviteId")
	if inviteID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do convite é obrigatório"})
		return
	}

	err := h.service.AcceptInvite(inviteID, userID)
	if err != nil {
		if err.Error() == "convite não encontrado" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "apenas o convidado pode alterar o status do convite" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "convite já foi respondido" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Convite aceito com sucesso"})
}

// DeclineInvite godoc
// @Summary Recusar convite
// @Description Recusa um convite para participar da mesa. Apenas o convidado pode recusar.
// @Tags GameTables
// @Produce json
// @Security Bearer
// @Param id path string true "ID da mesa"
// @Param inviteId path string true "ID do convite"
// @Success 200 {object} map[string]interface{} "Convite recusado"
// @Failure 400 {object} map[string]interface{} "ID inválido"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 403 {object} map[string]interface{} "Apenas convidado pode recusar"
// @Failure 404 {object} map[string]interface{} "Convite não encontrado"
// @Failure 409 {object} map[string]interface{} "Convite já respondido"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/v1/tables/{id}/invites/{inviteId}/decline [post]
func (h *GameTableHandler) DeclineInvite(c *gin.Context) {
	userID, _, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não encontrado no contexto"})
		return
	}

	inviteID := c.Param("inviteId")
	if inviteID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do convite é obrigatório"})
		return
	}

	err := h.service.DeclineInvite(inviteID, userID)
	if err != nil {
		if err.Error() == "convite não encontrado" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "apenas o convidado pode alterar o status do convite" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "convite já foi respondido" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Convite recusado"})
}
