package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/Dostonlv/hackathon-nt/internal/service"
)

type TenderHandler struct {
	tenderService *service.TenderService
}

func NewTenderHandler(tenderService *service.TenderService) *TenderHandler {
	return &TenderHandler{
		tenderService: tenderService,
	}
}

type CreateTenderRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Deadline    time.Time `json:"deadline" binding:"required,gtfield=time.Now"`
	Budget      float64   `json:"budget" binding:"required,gt=0"`
	FileURL     *string   `json:"file_url"`
}

func (h *TenderHandler) CreateTender(c *gin.Context) {
	var req CreateTenderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")
	tender, err := h.tenderService.CreateTender(c.Request.Context(), service.CreateTenderInput{
		ClientID:    userID.(uuid.UUID),
		Title:       req.Title,
		Description: req.Description,
		Deadline:    req.Deadline,
		Budget:      req.Budget,
		FileURL:     req.FileURL,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tender)
}

func (h *TenderHandler) ListTenders(c *gin.Context) {
	filters := service.TenderFilters{
		Status: c.Query("status"),
		Search: c.Query("search"),
	}

	tenders, err := h.tenderService.ListTenders(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tenders)
}
