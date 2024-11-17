// handlers/websocket_handler.go
package handlers

import (
	"net/http"

	"github.com/Dostonlv/hackathon-nt/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	notificationService *utils.NotificationService
	upgrader            websocket.Upgrader
}

func NewWebSocketHandler(notificationService *utils.NotificationService) *WebSocketHandler {
	return &WebSocketHandler{
		notificationService: notificationService,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Add appropriate origin checking for production
			},
		},
	}
}

// HandleWebSocket upgrades the HTTP connection to WebSocket
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Message: "Authentication required"})
		return
	}

	clientID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid user ID"})
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	// Register the client
	h.notificationService.RegisterClient(clientID, conn)

	// Handle disconnection
	go func() {
		defer func() {
			conn.Close()
			h.notificationService.UnregisterClient(clientID)
		}()

		for {
			// Keep the connection alive and handle any incoming messages
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}()
}
