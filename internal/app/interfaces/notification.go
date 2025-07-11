package interfaces

// NotificationService define interface para notificações em tempo real
type NotificationService interface {
	// Notificações de convites
	NotifyInviteCreated(tableID string, inviteData interface{})
	NotifyInviteAccepted(tableID string, inviteData interface{})
	NotifyInviteDeclined(tableID string, inviteData interface{})
	
	// Notificações de fichas
	NotifySheetCreated(tableID string, userID int, userEmail string, sheetData interface{})
	NotifySheetUpdated(tableID string, userID int, userEmail string, sheetData interface{})
	NotifySheetDeleted(tableID string, userID int, userEmail string, sheetData interface{})
	
	// Notificações de rolagens
	NotifyRollPerformed(tableID string, userID int, userEmail string, rollData interface{})
	
	// Notificações de mesa
	NotifyTableUpdated(tableID string, userID int, userEmail string, tableData interface{})
}
