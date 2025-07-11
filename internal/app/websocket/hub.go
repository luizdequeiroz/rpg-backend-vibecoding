package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// EventType representa tipos de eventos WebSocket
type EventType string

const (
	EventInviteCreated  EventType = "invite_created"
	EventInviteAccepted EventType = "invite_accepted"
	EventInviteDeclined EventType = "invite_declined"
	EventSheetCreated   EventType = "sheet_created"
	EventSheetUpdated   EventType = "sheet_updated"
	EventSheetDeleted   EventType = "sheet_deleted"
	EventRollPerformed  EventType = "roll_performed"
	EventTableUpdated   EventType = "table_updated"
)

// Event representa um evento WebSocket
type Event struct {
	Type      EventType   `json:"type"`
	UserID    int         `json:"user_id"`
	UserEmail string      `json:"user_email"`
	TableID   string      `json:"table_id"`
	Data      interface{} `json:"data"`
	Timestamp string      `json:"timestamp"`
}

// Client representa uma conexão WebSocket
type Client struct {
	conn    *websocket.Conn
	send    chan []byte
	hub     *Hub
	userID  int
	email   string
	tableID string
}

// Hub gerencia todas as conexões WebSocket
type Hub struct {
	// Clientes registrados agrupados por mesa
	clients map[string]map[*Client]bool

	// Canal para registrar novos clientes
	register chan *Client

	// Canal para desregistrar clientes
	unregister chan *Client

	// Mutex para operações thread-safe
	mutex sync.RWMutex
}

// NewHub cria um novo hub WebSocket
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run executa o hub em loop infinito
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			if h.clients[client.tableID] == nil {
				h.clients[client.tableID] = make(map[*Client]bool)
			}
			h.clients[client.tableID][client] = true
			h.mutex.Unlock()

			log.Printf("Cliente conectado: UserID=%d, Email=%s, TableID=%s",
				client.userID, client.email, client.tableID)

		case client := <-h.unregister:
			h.mutex.Lock()
			if tableClients, exists := h.clients[client.tableID]; exists {
				if _, exists := tableClients[client]; exists {
					delete(tableClients, client)
					close(client.send)

					// Remove mesa se não há mais clientes
					if len(tableClients) == 0 {
						delete(h.clients, client.tableID)
					}
				}
			}
			h.mutex.Unlock()

			log.Printf("Cliente desconectado: UserID=%d, Email=%s, TableID=%s",
				client.userID, client.email, client.tableID)
		}
	}
}

// BroadcastToTable envia evento para todos os clientes de uma mesa
func (h *Hub) BroadcastToTable(tableID string, eventType EventType, userID int, userEmail string, data interface{}) {
	event := Event{
		Type:      eventType,
		UserID:    userID,
		UserEmail: userEmail,
		TableID:   tableID,
		Data:      data,
		Timestamp: getTimestamp(),
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		log.Printf("Erro ao serializar evento: %v", err)
		return
	}

	h.mutex.RLock()
	tableClients := h.clients[tableID]
	h.mutex.RUnlock()

	if tableClients == nil {
		log.Printf("Nenhum cliente conectado na mesa %s", tableID)
		return
	}

	for client := range tableClients {
		select {
		case client.send <- eventJSON:
			// Evento enviado com sucesso
		default:
			// Cliente não está respondendo, desconectar
			h.mutex.Lock()
			delete(tableClients, client)
			close(client.send)
			if len(tableClients) == 0 {
				delete(h.clients, tableID)
			}
			h.mutex.Unlock()
		}
	}

	log.Printf("Evento %s enviado para %d clientes na mesa %s",
		eventType, len(tableClients), tableID)
}

// GetConnectedClients retorna número de clientes por mesa
func (h *Hub) GetConnectedClients() map[string]int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	result := make(map[string]int)
	for tableID, tableClients := range h.clients {
		result[tableID] = len(tableClients)
	}

	return result
}

// readPump lê mensagens da conexão WebSocket
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	// Configurar timeouts
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Erro WebSocket: %v", err)
			}
			break
		}
	}
}

// writePump escreve mensagens para a conexão WebSocket
func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Adicionar mensagens em fila ao mesmo writer
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// getTimestamp retorna timestamp atual em formato RFC3339
func getTimestamp() string {
	return time.Now().Format(time.RFC3339)
}
