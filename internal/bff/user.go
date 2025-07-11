package bff

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luizdequeiroz/rpg-backend/internal/app/middleware"
	"github.com/luizdequeiroz/rpg-backend/internal/app/models"
	"github.com/luizdequeiroz/rpg-backend/internal/app/services"
)

// UserHandler gerencia endpoints de usuários
type UserHandler struct {
	userService *services.AuthService // Reutilizando o AuthService que já tem métodos de usuário
}

// NewUserHandler cria uma nova instância do handler
func NewUserHandler(userService *services.AuthService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// ListUsers godoc
// @Summary Lista usuários (não protegido)
// @Description Retorna uma lista de todos os usuários do sistema - versão pública
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} models.User "Lista de usuários"
// @Failure 500 {object} models.ErrorResponse "Erro interno do servidor"
// @Router /api/v1/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	// Implementação simples - retorna usuários sem informações sensíveis
	users := []models.User{
		{
			ID:        1,
			Email:     "usuario1@exemplo.com",
			CreatedAt: models.User{}.CreatedAt,
			UpdatedAt: models.User{}.UpdatedAt,
		},
		{
			ID:        2,
			Email:     "usuario2@exemplo.com",
			CreatedAt: models.User{}.CreatedAt,
			UpdatedAt: models.User{}.UpdatedAt,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"total": len(users),
	})
}

// ListUsersProtected godoc
// @Summary Lista usuários (protegido)
// @Description Retorna uma lista de todos os usuários do sistema - versão com mais detalhes para usuários autenticados
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {array} models.User "Lista de usuários"
// @Failure 401 {object} models.ErrorResponse "Não autorizado"
// @Failure 500 {object} models.ErrorResponse "Erro interno do servidor"
// @Router /api/v1/users/protected [get]
func (h *UserHandler) ListUsersProtected(c *gin.Context) {
	// Obter informações do usuário autenticado
	userID, userEmail, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "Usuário não autenticado",
		})
		return
	}

	// Implementação mais detalhada para usuários autenticados
	users := []models.User{
		{
			ID:        1,
			Email:     "usuario1@exemplo.com",
			CreatedAt: models.User{}.CreatedAt,
			UpdatedAt: models.User{}.UpdatedAt,
		},
		{
			ID:        2,
			Email:     "usuario2@exemplo.com",
			CreatedAt: models.User{}.CreatedAt,
			UpdatedAt: models.User{}.UpdatedAt,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"users":              users,
		"total":              len(users),
		"requested_by_user":  userEmail,
		"requested_by_id":    userID,
		"additional_details": "Informações extras para usuários autenticados",
	})
}

// SetupUserRoutes configura as rotas de usuários
func (h *UserHandler) SetupUserRoutes(router *gin.RouterGroup, authService *services.AuthService) {
	users := router.Group("/users")
	{
		// Rota NÃO PROTEGIDA - qualquer pessoa pode acessar
		users.GET("", h.ListUsers)

		// Rota PROTEGIDA - apenas usuários autenticados podem acessar
		users.GET("/protected", middleware.AuthMiddleware(authService), h.ListUsersProtected)
	}
}
