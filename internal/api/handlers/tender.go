package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Dostonlv/hackathon-nt/internal/repository"
	"github.com/Dostonlv/hackathon-nt/internal/service"
	"github.com/Dostonlv/hackathon-nt/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/k0kubun/pp"
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
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Deadline    time.Time `json:"deadline" binding:"required,gtfield=time.Now"`
	Budget      float64   `json:"budget" binding:"required,gt=0"`
}

// CreateTender godoc
// @Summary Create a new tender
// @Description Create a new tender with the provided details and a PDF file
// @Tags tenders
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "Title of the tender"
// @Param description formData string true "Description of the tender"
// @Param deadline formData string true "Deadline for the tender"
// @Param budget formData float64 true "Budget for the tender"
// @Param file formData file true "PDF file for the tender"
// @Success 201 {object} models.Tender
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/tenders [post]
func (h *TenderHandler) CreateTender(c *gin.Context) {
	var req CreateTenderRequest

	req.Title = c.Request.FormValue("title")
	req.Description = c.Request.FormValue("description")
	budget, err := strconv.ParseFloat(c.Request.FormValue("budget"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid budget value"})
		return
	}
	req.Budget = budget

	deadlineStr := c.Request.FormValue("deadline")
	deadline, err := time.Parse(time.RFC3339, deadlineStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deadline format"})
		return
	}
	req.Deadline = deadline

	pp.Print("req: ", req)
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	// Validate file type
	if file.Header.Get("Content-Type") != "application/pdf" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only PDF files are allowed"})
		return
	}

	// Save the file to a specific location
	filePath := "./uploads/" + file.Filename
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
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

	tender, err := h.tenderService.CreateTender(c.Request.Context(), service.CreateTenderInput{
		ClientID:    claims.UserID,
		Title:       req.Title,
		Description: req.Description,
		Deadline:    req.Deadline,
		Budget:      req.Budget,
		FileURL:     &filePath,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tender)
}

func (h *TenderHandler) ListTenders(c *gin.Context) {
	filters := repository.TenderFilters{
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
