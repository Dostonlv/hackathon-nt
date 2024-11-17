package handlers

import (
	"net/http"
	"time"

	"github.com/Dostonlv/hackathon-nt/internal/repository"
	"github.com/Dostonlv/hackathon-nt/internal/service"
	"github.com/Dostonlv/hackathon-nt/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TenderHandler struct {
	tenderService *service.TenderService
}

func NewTenderHandler(tenderService *service.TenderService) *TenderHandler {
	if tenderService == nil {
		panic("tenderService cannot be nil")
	}
	return &TenderHandler{
		tenderService: tenderService,
	}
}

type CreateTenderRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Deadline    string  `json:"deadline" datetime:"2006-01-02T15:04:05Z07:00"`
	Budget      float64 `json:"budget" `
	Attachment  *string `json:"attachment"`
}

// CreateTender godoc
// @Summary Create a new tender
// @Description Create a new tender with the provided details
// @Tags tenders
// @Accept json
// @Produce json
// @Param tender body CreateTenderRequest true "Tender details"
// @Success 201 {object} models.Tender
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/client/tenders [post]
func (h *TenderHandler) CreateTender(c *gin.Context) {
	var req CreateTenderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Missing token"})
		return
	}

	if req.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	if req.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	if req.Deadline == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	if req.Budget <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid tender data"})
		return
	}

	claims, err := utils.ParseToken(authHeader, []byte("secreeet"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	if claims.UserID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	deadline, err := time.Parse(time.RFC3339, req.Deadline)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tender data"})
		return
	}

	if time.Now().After(deadline) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Deadline has already passed"})
		return
	}

	tender, err := h.tenderService.CreateTender(c.Request.Context(), service.CreateTenderInput{
		ClientID:    claims.UserID,
		Title:       req.Title,
		Description: req.Description,
		Deadline:    deadline,
		Budget:      req.Budget,
		Attachment:  req.Attachment,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": tender.ID, "title": tender.Title})
}

// ListTenders handles the request to list tenders based on provided filters.
// @Summary List Tenders
// @Description Retrieves a list of tenders filtered by status and search query.
// @Tags tenders
// @Accept json
// @Produce json
// @Param status query string false "Filter tenders by status"
// @Param search query string false "Search tenders by keyword"
// @Success 200 {array} repository.Tender "List of tenders"
// @Failure 500 {object} []models.Tender "Internal Server Error"
// @Router /api/client/tenders [get]
func (h *TenderHandler) ListTenders(c *gin.Context) {
	clientID, _ := c.Get("userId")

	clientUUID, err := uuid.Parse(clientID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID format"})
		return
	}
	tenders, err := h.tenderService.ListTenders(c.Request.Context(), clientUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tenders)
}

type UpdateTenderStatusRequest struct {
	Status string `json:"status"`
}

// UpdateTenderStatus godoc
// @Summary Update the status of a tender
// @Description Update the status of a tender by its ID
// @Tags tenders
// @Accept json
// @Produce json
// @Param id path string true "Tender ID"
// @Param status body UpdateTenderStatusRequest true "New status"
// @Success 200 {object} models.Tender
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/client/tenders/{id} [put]
func (h *TenderHandler) UpdateTenderStatus(c *gin.Context) {
	var req UpdateTenderStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Status == "" || req.Status != "open" && req.Status != "closed" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid tender status"})
		return
	}

	tenderID := c.Param("id")
	if tenderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Tender ID is required"})
		return
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Missing token"})
		return
	}

	claims, err := utils.ParseToken(authHeader, []byte("secreeet"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	if claims.UserID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	tenderUUID, err := uuid.Parse(tenderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Tender not found"})
		return
	}

	_, err = h.tenderService.UpdateTender(c.Request.Context(), service.UpdateTenderInput{ID: tenderUUID, Status: &req.Status})
	if err != nil {
		if err == service.ErrTenderNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tender not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tender status updated"})
}

// GetTenderByID godoc
// @Summary Get a tender by ID
// @Description Retrieve a tender by its ID
// @Tags tenders
// @Accept json
// @Produce json
// @Param id path string true "Tender ID"
// @Success 200 {object} models.Tender
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/client/tenders/{id} [get]
func (h *TenderHandler) GetTenderByID(c *gin.Context) {
	tenderID := c.Param("id")
	if tenderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Tender ID is required"})
		return
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Missing token"})
		return
	}

	claims, err := utils.ParseToken(authHeader, []byte("secreeet"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	if claims.UserID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	tenderUUID, err := uuid.Parse(tenderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tender ID format"})
		return
	}

	tender, err := h.tenderService.GetTenderByID(c.Request.Context(), tenderUUID)
	if err != nil {
		if err == service.ErrTenderNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "Tender not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tender)
}

// DeleteTender godoc
// @Summary Delete a tender
// @Description Delete a tender by its ID
// @Tags tenders
// @Accept json
// @Produce json
// @Param id path string true "Tender ID"
// @Success 200 {object} gin.H{"message": "Tender deleted"}
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/client/tenders/{id} [delete]
func (h *TenderHandler) DeleteTender(c *gin.Context) {
	tenderID := c.Param("id")
	if tenderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Tender ID is required"})
		return
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Missing token"})
		return
	}

	claims, err := utils.ParseToken(authHeader, []byte("secreeet"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	if claims.UserID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	tenderUUID, err := uuid.Parse(tenderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Tender not found or access denied"})
		return
	}

	err = h.tenderService.DeleteTender(c.Request.Context(), tenderUUID, claims.UserID)
	if err != nil {
		if err == service.ErrTenderNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "Tender not found or access denied"})
			return
		}
		if err == service.ErrUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tender deleted successfully"})
}

// ListTendersFiltering handles the request to list tenders with filters.
// @Summary List Tenders with Filters
// @Description Retrieves a list of tenders filtered by various criteria.
// @Tags tenders
// @Accept json
// @Produce json
// @Param status query string false "Filter tenders by status"
// @Param search query string false "Search tenders by keyword"
// @Param minBudget query number false "Minimum budget"
// @Param maxBudget query number false "Maximum budget"
// @Success 200 {array} models.Tender "List of tenders"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /api/client/tenders/filter [get]
func (h *TenderHandler) ListTendersFiltering(c *gin.Context) {
	var filters repository.TenderFilters

	if status := c.Query("status"); status != "" {
		filters.Status = status
	}
	if search := c.Query("search"); search != "" {
		filters.Search = search
	}

	tenders, err := h.tenderService.ListTendersFiltering(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tenders)
}
