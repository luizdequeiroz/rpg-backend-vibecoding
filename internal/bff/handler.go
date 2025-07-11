package bff

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luizdequeiroz/rpg-backend/internal/app/repositories"
	"github.com/luizdequeiroz/rpg-backend/internal/app/services"
	"github.com/luizdequeiroz/rpg-backend/pkg/db"
)

// Handler contém as dependências da camada BFF
type Handler struct {
	db                   *db.DB
	authService          *services.AuthService
	authHandler          *AuthHandler
	sheetTemplateService *services.SheetTemplateService
	sheetTemplateHandler *SheetTemplateHandler
	userHandler          *UserHandler
	gameTableService     *services.GameTableService
	gameTableHandler     *GameTableHandler
}

// NewHandler cria um novo handler BFF
func NewHandler(database *db.DB) *Handler {
	authService := services.NewAuthService(database)
	authHandler := NewAuthHandler(authService)

	sheetTemplateService := services.NewSheetTemplateService(database)
	sheetTemplateHandler := NewSheetTemplateHandler(sheetTemplateService)

	userHandler := NewUserHandler(authService)

	// Inicializar repositórios e serviços para GameTable
	gameTableRepo := repositories.NewGameTableRepository(database)
	inviteRepo := repositories.NewInviteRepository(database)
	gameTableService := services.NewGameTableService(gameTableRepo, inviteRepo)
	gameTableHandler := NewGameTableHandler(gameTableService)

	return &Handler{
		db:                   database,
		authService:          authService,
		authHandler:          authHandler,
		sheetTemplateService: sheetTemplateService,
		sheetTemplateHandler: sheetTemplateHandler,
		userHandler:          userHandler,
		gameTableService:     gameTableService,
		gameTableHandler:     gameTableHandler,
	}
}

// SetupRoutes configura todas as rotas da API v1
func (h *Handler) SetupRoutes(router *gin.RouterGroup) {
	// Rotas de autenticação
	h.authHandler.SetupAuthRoutes(router, h.authService)

	// Rotas de templates de ficha
	h.sheetTemplateHandler.SetupTemplateRoutes(router)

	// Rotas de usuários (usando o novo UserHandler)
	h.userHandler.SetupUserRoutes(router, h.authService)

	// Rotas de mesas de jogo
	h.gameTableHandler.SetupGameTableRoutes(router, h.authService)

	// Rotas de campanhas
	campaigns := router.Group("/campaigns")
	{
		campaigns.GET("", h.listCampaigns)
		campaigns.POST("", h.createCampaign)
		campaigns.GET("/:id", h.getCampaign)
		campaigns.PUT("/:id", h.updateCampaign)
		campaigns.DELETE("/:id", h.deleteCampaign)
	}

	// Rotas de personagens
	characters := router.Group("/characters")
	{
		characters.GET("", h.listCharacters)
		characters.POST("", h.createCharacter)
		characters.GET("/:id", h.getCharacter)
		characters.PUT("/:id", h.updateCharacter)
		characters.DELETE("/:id", h.deleteCharacter)
	}

	// Rotas de sessões
	sessions := router.Group("/sessions")
	{
		sessions.GET("", h.listSessions)
		sessions.POST("", h.createSession)
		sessions.GET("/:id", h.getSession)
		sessions.PUT("/:id", h.updateSession)
		sessions.DELETE("/:id", h.deleteSession)
	}
}

// Handlers temporários - serão implementados com a lógica de negócio real

// Users handlers
func (h *Handler) listUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Lista de usuários",
		"data":    []interface{}{},
	})
}

func (h *Handler) createUser(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"message": "Usuário criado com sucesso",
	})
}

func (h *Handler) getUser(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Detalhes do usuário",
		"id":      id,
	})
}

func (h *Handler) updateUser(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Usuário atualizado",
		"id":      id,
	})
}

func (h *Handler) deleteUser(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Usuário removido",
		"id":      id,
	})
}

// Campaigns handlers
func (h *Handler) listCampaigns(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Lista de campanhas",
		"data":    []interface{}{},
	})
}

func (h *Handler) createCampaign(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"message": "Campanha criada com sucesso",
	})
}

func (h *Handler) getCampaign(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Detalhes da campanha",
		"id":      id,
	})
}

func (h *Handler) updateCampaign(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Campanha atualizada",
		"id":      id,
	})
}

func (h *Handler) deleteCampaign(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Campanha removida",
		"id":      id,
	})
}

// Characters handlers
func (h *Handler) listCharacters(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Lista de personagens",
		"data":    []interface{}{},
	})
}

func (h *Handler) createCharacter(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"message": "Personagem criado com sucesso",
	})
}

func (h *Handler) getCharacter(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Detalhes do personagem",
		"id":      id,
	})
}

func (h *Handler) updateCharacter(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Personagem atualizado",
		"id":      id,
	})
}

func (h *Handler) deleteCharacter(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Personagem removido",
		"id":      id,
	})
}

// Sessions handlers
func (h *Handler) listSessions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Lista de sessões",
		"data":    []interface{}{},
	})
}

func (h *Handler) createSession(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"message": "Sessão criada com sucesso",
	})
}

func (h *Handler) getSession(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Detalhes da sessão",
		"id":      id,
	})
}

func (h *Handler) updateSession(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Sessão atualizada",
		"id":      id,
	})
}

func (h *Handler) deleteSession(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Sessão removida",
		"id":      id,
	})
}
