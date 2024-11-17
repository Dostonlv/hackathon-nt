package handlers

import (
	"net/http"

	"github.com/Dostonlv/hackathon-nt/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// HistoryHandler represents the handler for history-related API endpoints.
type HistoryHandler struct {
	historyService *service.HistoryService
}

// NewHistoryHandler creates a new instance of HistoryHandler.
func NewHistoryHandler(historyService *service.HistoryService) *HistoryHandler {
	return &HistoryHandler{
		historyService: historyService,
	}
}

// GetTenderHistory handles the GET /users/:id/tenders endpoint.
// @Summary Get tender history
// @Description Get the tender history for a user
// @Tags History
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {array} models.Tender "OK"
// @Security BearerAuth
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /api/users/{id}/tenders [get]
func (h *HistoryHandler) GetTenderHistory(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	tenders, err := h.historyService.GetTenderHistory(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tender history"})
		return
	}

	c.JSON(http.StatusOK, tenders)
}

// GetBidHistory handles the GET /users/:id/bids endpoint.
// @Summary Get bid history
// @Description Get the bid history for a user
// @Tags History
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {array} models.Bid "OK"
// @Security BearerAuth
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /api/users/{id}/bids [get]
func (h *HistoryHandler) GetBidHistory(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	bids, err := h.historyService.GetBidHistory(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve bid history"})
		return
	}

	c.JSON(http.StatusOK, bids)
}
