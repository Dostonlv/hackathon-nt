package handlers

import (
	"net/http"
	"strconv"

	"github.com/Dostonlv/hackathon-nt/internal/repository"
	"github.com/Dostonlv/hackathon-nt/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BidHandler struct {
	bidService *service.BidService
}

func NewBidHandler(bidService *service.BidService) *BidHandler {
	return &BidHandler{
		bidService: bidService,
	}
}

type CreateBidRequest struct {
	Price        float64 `json:"price" binding:"required,gt=0"`
	DeliveryTime int     `json:"delivery_time" binding:"required,gt=0"`
	Comments     string  `json:"comments"`
}

func (h *BidHandler) CreateBid(c *gin.Context) {
	tenderID, err := uuid.Parse(c.Param("tender_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tender ID"})
		return
	}

	var req CreateBidRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")
	bid, err := h.bidService.CreateBid(c.Request.Context(), service.CreateBidInput{
		TenderID:     tenderID,
		ContractorID: userID.(uuid.UUID),
		Price:        req.Price,
		DeliveryTime: req.DeliveryTime,
		Comments:     req.Comments,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, bid)
}

func (h *BidHandler) ListBids(c *gin.Context) {
	tenderID, err := uuid.Parse(c.Param("tender_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tender ID"})
		return
	}

	filters := repository.BidFilters{
		Status:      c.Query("status"),
		Search:      c.Query("search"),
		MinPrice:    parseFloatQuery(c, "min_price"),
		MaxPrice:    parseFloatQuery(c, "max_price"),
		MinDelivery: parseIntQuery(c, "min_delivery_time"),
		MaxDelivery: parseIntQuery(c, "max_delivery_time"),
		SortBy:      c.Query("sort_by"),
		SortOrder:   c.Query("sort_order"),
	}

	bids, err := h.bidService.ListBids(c.Request.Context(), tenderID, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bids)
}

func parseFloatQuery(c *gin.Context, key string) *float64 {
	if value := c.Query(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return &floatValue
		}
	}
	return nil
}

func parseIntQuery(c *gin.Context, key string) *int {
	if value := c.Query(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return &intValue
		}
	}
	return nil
}
