package websocket

import (
	"log"
	
	"github.com/luizdequeiroz/rpg-backend/internal/app/interfaces"
)

// WebSocketService integra WebSocket com outras camadas
// Implementa interfaces.NotificationService
type WebSocketService struct {
	hub *Hub
}

// Verificar em tempo de compilação se implementa a interface
var _ interfaces.NotificationService = (*WebSocketService)(nil)

// NewWebSocketService cria novo serviço WebSocket
func NewWebSocketService(hub *Hub) *WebSocketService {
	return &WebSocketService{
		hub: hub,
	}
}

// NotifyInviteCreated notifica criação de convite
func (ws *WebSocketService) NotifyInviteCreated(tableID string, inviteData interface{}) {
	log.Printf("WebSocket: Notificando criação de convite na mesa %s", tableID)
	ws.hub.BroadcastToTable(tableID, EventInviteCreated, 0, "sistema", inviteData)
}

// NotifyInviteAccepted notifica aceite de convite
func (ws *WebSocketService) NotifyInviteAccepted(tableID string, inviteData interface{}) {
	log.Printf("WebSocket: Notificando aceite de convite na mesa %s", tableID)
	ws.hub.BroadcastToTable(tableID, EventInviteAccepted, 0, "sistema", inviteData)
}

// NotifyInviteDeclined notifica recusa de convite
func (ws *WebSocketService) NotifyInviteDeclined(tableID string, inviteData interface{}) {
	log.Printf("WebSocket: Notificando recusa de convite na mesa %s", tableID)
	ws.hub.BroadcastToTable(tableID, EventInviteDeclined, 0, "sistema", inviteData)
}

// NotifySheetCreated notifica criação de ficha
func (ws *WebSocketService) NotifySheetCreated(tableID string, userID int, userEmail string, sheetData interface{}) {
	log.Printf("WebSocket: Notificando criação de ficha na mesa %s por usuário %d", tableID, userID)
	ws.hub.BroadcastToTable(tableID, EventSheetCreated, userID, userEmail, sheetData)
}

// NotifySheetUpdated notifica atualização de ficha
func (ws *WebSocketService) NotifySheetUpdated(tableID string, userID int, userEmail string, sheetData interface{}) {
	log.Printf("WebSocket: Notificando atualização de ficha na mesa %s por usuário %d", tableID, userID)
	ws.hub.BroadcastToTable(tableID, EventSheetUpdated, userID, userEmail, sheetData)
}

// NotifySheetDeleted notifica exclusão de ficha
func (ws *WebSocketService) NotifySheetDeleted(tableID string, userID int, userEmail string, sheetData interface{}) {
	log.Printf("WebSocket: Notificando exclusão de ficha na mesa %s por usuário %d", tableID, userID)
	ws.hub.BroadcastToTable(tableID, EventSheetDeleted, userID, userEmail, sheetData)
}

// NotifyRollPerformed notifica rolagem de dados
func (ws *WebSocketService) NotifyRollPerformed(tableID string, userID int, userEmail string, rollData interface{}) {
	log.Printf("WebSocket: Notificando rolagem na mesa %s por usuário %d", tableID, userID)
	ws.hub.BroadcastToTable(tableID, EventRollPerformed, userID, userEmail, rollData)
}

// NotifyTableUpdated notifica atualização da mesa
func (ws *WebSocketService) NotifyTableUpdated(tableID string, userID int, userEmail string, tableData interface{}) {
	log.Printf("WebSocket: Notificando atualização da mesa %s por usuário %d", tableID, userID)
	ws.hub.BroadcastToTable(tableID, EventTableUpdated, userID, userEmail, tableData)
}

// GetConnectedClients retorna clientes conectados por mesa
func (ws *WebSocketService) GetConnectedClients() map[string]int {
	return ws.hub.GetConnectedClients()
}

// GetTotalConnections retorna total de conexões ativas
func (ws *WebSocketService) GetTotalConnections() int {
	stats := ws.hub.GetConnectedClients()
	total := 0
	for _, count := range stats {
		total += count
	}
	return total
}
