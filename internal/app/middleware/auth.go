package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/luizdequeiroz/rpg-backend/internal/app/models"
	"github.com/luizdequeiroz/rpg-backend/internal/app/services"
)

// AuthMiddleware middleware para autenticação JWT
func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extrair token do header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "Token de autorização requerido",
			})
			c.Abort()
			return
		}

		// Verificar se o header tem o formato "Bearer <token>"
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "Formato de token inválido. Use: Bearer <token>",
			})
			c.Abort()
			return
		}

		token := bearerToken[1]

		// Validar token
		claims, err := authService.ValidateJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error:   "Token inválido",
				Message: err.Error(),
			})
			c.Abort()
			return
		}

		// Extrair informações do usuário dos claims
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "Token malformado: user_id inválido",
			})
			c.Abort()
			return
		}

		userID := int(userIDFloat)
		email, ok := claims["email"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "Token malformado: email inválido",
			})
			c.Abort()
			return
		}

		// Adicionar informações do usuário ao contexto
		c.Set("user_id", userID)
		c.Set("user_email", email)

		c.Next()
	}
}

// GetUserFromContext extrai as informações do usuário do contexto
func GetUserFromContext(c *gin.Context) (userID int, email string, exists bool) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		return 0, "", false
	}

	emailValue, exists := c.Get("user_email")
	if !exists {
		return 0, "", false
	}

	userID, ok := userIDValue.(int)
	if !ok {
		return 0, "", false
	}

	email, ok = emailValue.(string)
	if !ok {
		return 0, "", false
	}

	return userID, email, true
}

// OptionalAuthMiddleware middleware opcional para autenticação (não aborta se não autenticado)
func OptionalAuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != "bearer" {
			c.Next()
			return
		}

		token := bearerToken[1]
		claims, err := authService.ValidateJWT(token)
		if err != nil {
			c.Next()
			return
		}

		if userIDFloat, ok := claims["user_id"].(float64); ok {
			if email, ok := claims["email"].(string); ok {
				c.Set("user_id", int(userIDFloat))
				c.Set("user_email", email)
			}
		}

		c.Next()
	}
}
