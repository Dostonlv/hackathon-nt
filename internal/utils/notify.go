package utils

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type NotificationService struct {
	// Map of client ID to their WebSocket connection
	clients sync.Map
}

type BidNotification struct {
	Type     string    `json:"type"`
	TenderID uuid.UUID `json:"tender_id"`
	BidID    uuid.UUID `json:"bid_id"`
	Price    float64   `json:"price"`
	Message  string    `json:"message"`
}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

// RegisterClient registers a new WebSocket client
func (s *NotificationService) RegisterClient(clientID uuid.UUID, conn *websocket.Conn) {
	s.clients.Store(clientID.String(), conn)
}

// UnregisterClient removes a WebSocket client
func (s *NotificationService) UnregisterClient(clientID uuid.UUID) {
	s.clients.Delete(clientID)
}

// NotifyNewBid sends a notification to a specific client about a new bid
func (s *NotificationService) NotifyNewBid(ctx context.Context, clientID uuid.UUID, notification BidNotification) error {
	conn, ok := s.clients.Load(clientID.String())
	if !ok {
		return nil // Client not connected, silently ignore
	}

	wsConn := conn.(*websocket.Conn)
	message, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	return wsConn.WriteMessage(websocket.TextMessage, message)
}

type AwardNotification struct {
	Type     string    `json:"type"`
	TenderID uuid.UUID `json:"tender_id"`
	AwardID  uuid.UUID `json:"award_id"`
	Message  string    `json:"message"`
}

// NotifyAward sends a notification to a specific client about an award
func (s *NotificationService) NotifyAward(ctx context.Context, clientID uuid.UUID, notification BidNotification) error {
	conn, ok := s.clients.Load(clientID.String())
	if !ok {
		return nil // Client not connected, silently ignore
	}

	wsConn := conn.(*websocket.Conn)
	message, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	return wsConn.WriteMessage(websocket.TextMessage, message)
}
