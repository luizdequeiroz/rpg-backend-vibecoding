package bff

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luizdequeiroz/rpg-backend/internal/app/middleware"
	"github.com/luizdequeiroz/rpg-backend/internal/app/models"
	"github.com/luizdequeiroz/rpg-backend/internal/app/services"
)

// AuthHandler gerencia endpoints de autenticação
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler cria uma nova instância do handler de autenticação
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Signup godoc
// @Summary Registro de novo usuário
// @Description Registra um novo usuário no sistema com email e senha
// @Tags auth
// @Accept json
// @Produce json
// @Param signup body models.UserSignupRequest true "Dados de registro"
// @Success 201 {object} models.AuthResponse "Usuário criado com sucesso"
// @Failure 400 {object} models.ErrorResponse "Dados inválidos"
// @Failure 409 {object} models.ErrorResponse "Email já está em uso"
// @Failure 500 {object} models.ErrorResponse "Erro interno do servidor"
// @Router /api/v1/auth/signup [post]
func (h *AuthHandler) Signup(c *gin.Context) {
	var req models.UserSignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Dados inválidos",
			Message: err.Error(),
		})
		return
	}

	// Validações básicas
	if req.Email == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Email é obrigatório",
		})
		return
	}

	if len(req.Password) < 6 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Senha deve ter pelo menos 6 caracteres",
		})
		return
	}

	// Registrar usuário
	response, err := h.authService.Signup(req)
	if err != nil {
		if err == services.ErrEmailAlreadyExists {
			c.JSON(http.StatusConflict, models.ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Erro interno do servidor",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// Login godoc
// @Summary Login de usuário
// @Description Autentica um usuário e retorna token JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param login body models.UserLoginRequest true "Credenciais de login"
// @Success 200 {object} models.AuthResponse "Login realizado com sucesso"
// @Failure 400 {object} models.ErrorResponse "Dados inválidos"
// @Failure 401 {object} models.ErrorResponse "Credenciais inválidas"
// @Failure 500 {object} models.ErrorResponse "Erro interno do servidor"
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Dados inválidos",
			Message: err.Error(),
		})
		return
	}

	// Validações básicas
	if req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Email e senha são obrigatórios",
		})
		return
	}

	// Autenticar usuário
	response, err := h.authService.Login(req)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "Email ou senha incorretos",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Erro interno do servidor",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Me godoc
// @Summary Informações do usuário atual
// @Description Retorna as informações do usuário autenticado
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User "Informações do usuário"
// @Failure 401 {object} models.ErrorResponse "Token inválido ou ausente"
// @Failure 404 {object} models.ErrorResponse "Usuário não encontrado"
// @Failure 500 {object} models.ErrorResponse "Erro interno do servidor"
// @Router /api/v1/auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	userID, _, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "Usuário não autenticado",
		})
		return
	}

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		if err == services.ErrUserNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "Usuário não encontrado",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Erro interno do servidor",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// SetupAuthRoutes configura as rotas de autenticação
func (h *AuthHandler) SetupAuthRoutes(router *gin.RouterGroup, authService *services.AuthService) {
	auth := router.Group("/auth")
	{
		auth.POST("/signup", h.Signup)
		auth.POST("/login", h.Login)
		auth.GET("/me", middleware.AuthMiddleware(authService), h.Me)
	}
}
