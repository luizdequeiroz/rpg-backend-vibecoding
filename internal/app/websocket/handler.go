package websocket

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/luizdequeiroz/rpg-backend/internal/app/middleware"
)

// WebSocketHandler gerencia conexões WebSocket
type WebSocketHandler struct {
	hub      *Hub
	upgrader websocket.Upgrader
}

// NewWebSocketHandler cria novo handler WebSocket
func NewWebSocketHandler(hub *Hub) *WebSocketHandler {
	return &WebSocketHandler{
		hub: hub,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				// Em produção, implementar verificação de origin adequada
				return true
			},
		},
	}
}

// HandleWebSocket gerencia upgrade para WebSocket
// @Summary Conectar WebSocket
// @Description Estabelece conexão WebSocket para receber notificações em tempo real de uma mesa
// @Tags WebSocket
// @Param table_id query string true "ID da mesa para receber notificações"
// @Security BearerAuth
// @Success 101 {string} string "Switching Protocols"
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/ws [get]
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	// Verificar autenticação JWT
	userID, userEmail, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token JWT inválido"})
		return
	}

	// Obter tableID dos query parameters
	tableID := c.Query("table_id")
	if tableID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "table_id é obrigatório"})
		return
	}

	// TODO: Verificar se usuário tem acesso à mesa
	// Por enquanto, permitir acesso a qualquer mesa para simplificar

	// Fazer upgrade para WebSocket
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao estabelecer conexão WebSocket",
			"details": err.Error(),
		})
		return
	}

	// Criar cliente
	client := &Client{
		conn:    conn,
		send:    make(chan []byte, 256),
		hub:     h.hub,
		userID:  userID,
		email:   userEmail,
		tableID: tableID,
	}

	// Registrar cliente no hub
	h.hub.register <- client

	// Iniciar goroutines para leitura e escrita
	go client.writePump()
	go client.readPump()
}

// GetStats retorna estatísticas das conexões WebSocket
// @Summary Estatísticas WebSocket
// @Description Retorna número de clientes conectados por mesa
// @Tags WebSocket
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/ws/stats [get]
func (h *WebSocketHandler) GetStats(c *gin.Context) {
	// Verificar autenticação (apenas usuários autenticados podem ver stats)
	_, _, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token JWT inválido"})
		return
	}

	stats := h.hub.GetConnectedClients()

	totalClients := 0
	for _, count := range stats {
		totalClients += count
	}

	c.JSON(http.StatusOK, gin.H{
		"total_clients":     totalClients,
		"clients_per_table": stats,
		"active_tables":     len(stats),
		"timestamp":         getTimestamp(),
	})
}

// BroadcastTestEvent envia evento de teste (apenas para desenvolvimento)
// @Summary Evento de teste WebSocket
// @Description Envia evento de teste para uma mesa específica (apenas desenvolvimento)
// @Tags WebSocket
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "Dados do evento de teste"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/ws/test [post]
func (h *WebSocketHandler) BroadcastTestEvent(c *gin.Context) {
	// Verificar autenticação
	userID, userEmail, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token JWT inválido"})
		return
	}

	var req struct {
		TableID   string      `json:"table_id" binding:"required"`
		EventType string      `json:"event_type" binding:"required"`
		Data      interface{} `json:"data"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Dados inválidos",
			"details": err.Error(),
		})
		return
	}

	// Enviar evento de teste
	h.hub.BroadcastToTable(req.TableID, EventType(req.EventType), userID, userEmail, req.Data)

	c.JSON(http.StatusOK, gin.H{
		"message":    "Evento de teste enviado",
		"table_id":   req.TableID,
		"event_type": req.EventType,
		"user_id":    userID,
		"timestamp":  getTimestamp(),
	})
}
