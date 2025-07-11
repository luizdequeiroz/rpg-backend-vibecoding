// Package main RPG Session Backend API
//
// # Documentação da API para gerenciamento de sessões de RPG
//
// Termos de Serviço: http://swagger.io/terms/
//
// Contato: Luiz de Queiroz <luiz@example.com>
//
// Versão: 1.0.0
// Host: localhost:8080
// BasePath: /
// Schemes: http
//
// @title RPG Session Backend API
// @version 1.0.0
// @description API para gerenciamento de sessões de RPG, incluindo usuários, campanhas, personagens e sessões de jogo.
// @termsOfService http://swagger.io/terms/
// @contact.name Luiz de Queiroz
// @contact.email luiz@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Digite 'Bearer' seguido do token JWT (ex: Bearer eyJhbGciOiJIUzI1...)
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/luizdequeiroz/rpg-backend/docs" // docs gerados pelo swag
	"github.com/luizdequeiroz/rpg-backend/internal/bff"
	"github.com/luizdequeiroz/rpg-backend/pkg/config"
	"github.com/luizdequeiroz/rpg-backend/pkg/db"
)

func main() {
	// Carregar configurações
	cfg := config.Load()
	log.Printf("Configurações carregadas: Host=%s, Port=%d", cfg.Server.Host, cfg.Server.Port)

	// Conectar ao banco de dados
	var database *db.DB
	var err error

	if cfg.Database.URL != "" {
		database, err = db.NewDBWithDSN(cfg.Database.URL)
	} else {
		database, err = db.NewDB()
	}

	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer database.Close()

	// Executar migrações automaticamente na inicialização
	if err := database.RunMigrations(); err != nil {
		log.Fatalf("Erro ao executar migrações: %v", err)
	}

	// Verificar saúde da conexão
	if err := database.Health(); err != nil {
		log.Fatalf("Erro na verificação de saúde do banco: %v", err)
	}

	// Configurar modo do Gin baseado no ambiente
	if cfg.Log.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Criar router HTTP
	router := gin.New()

	// Middlewares globais
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// CORS middleware simples
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// Criar handler BFF
	bffHandler := bff.NewHandler(database)

	// Healthcheck endpoint
	router.GET("/health", bffHandler.HealthHandler)

	// Swagger documentation
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	v1 := router.Group("/api/v1")
	bffHandler.SetupRoutes(v1)

	// Endpoint raiz com informações da API
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":        "RPG Session Backend",
			"version":     "1.0.0",
			"description": "Backend para gerenciamento de sessões de RPG",
			"endpoints": gin.H{
				"health": "/health",
				"docs":   "/docs/index.html",
				"api_v1": "/api/v1",
				"auth": gin.H{
					"signup": "/api/v1/auth/signup",
					"login":  "/api/v1/auth/login",
					"me":     "/api/v1/auth/me",
				},
				"templates":  "/api/v1/templates",
				"users":      "/api/v1/users",
				"campaigns":  "/api/v1/campaigns",
				"characters": "/api/v1/characters",
				"sessions":   "/api/v1/sessions",
			},
		})
	})

	// Configurar servidor
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Iniciar servidor em goroutine
	go func() {
		log.Printf("Servidor RPG Backend iniciado com sucesso!")
		log.Printf("Banco de dados: %s", database.GetDSN())
		log.Printf("Servidor rodando em: http://%s", addr)
		log.Printf("Healthcheck: http://%s/health", addr)
		log.Printf("API v1: http://%s/api/v1", addr)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao iniciar servidor: %v", err)
		}
	}()

	// Aguardar sinal de interrupção
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Desligando servidor...")

	// Graceful shutdown
	if err := server.Close(); err != nil {
		log.Printf("Erro durante shutdown do servidor: %v", err)
	}

	log.Println("Servidor finalizado com sucesso.")
}
