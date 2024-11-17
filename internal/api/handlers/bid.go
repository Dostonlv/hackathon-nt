package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Dostonlv/hackathon-nt/internal/repository"
	"github.com/Dostonlv/hackathon-nt/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/k0kubun/pp"
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
	Price        float64 `json:"price"`
	DeliveryTime int     `json:"delivery_time"`
	Comments     string  `json:"comments"`
}

type Bid struct {
	ID           uuid.UUID `json:"id"`
	TenderID     uuid.UUID `json:"tender_id"`
	ContractorID uuid.UUID `json:"contractor_id"`
	Price        float64   `json:"price"`
	DeliveryTime int       `json:"delivery_time"`
	Comments     string    `json:"comments"`
	Status       string    `json:"status"`
}

// CreateBid handles the creation of a new bid for a specific tender.
//
// @Summary Create a new bid
// @Description This endpoint allows a contractor to create a new bid for a specified tender. The contractor must provide the bid details in the request body.
// @Tags bids
// @Accept json
// @Produce json
// @Param tender_id path string true "Tender ID"
// @Param bid body CreateBidRequest true "Bid details"
// @Success 201 {object} Bid "Successfully created bid"
// @Failure 400 {object} ErrorResponse "Invalid tender ID or bad request body"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Security ApiKeyAuth
// @Router /api/contractor/tenders/{tender_id}/bid [post]
func (h *BidHandler) CreateBid(c *gin.Context) {
	tenderID, err := uuid.Parse(c.Param("tender_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Message: "Tender not found"})
		return
	}

	var req CreateBidRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}

	if req.Price <= 0 || req.DeliveryTime <= 0 || req.Comments == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid bid data"})
		return
	}

	userID, _ := c.Get("userId")
	contractorID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid user ID"})
		return
	}

	bid, err := h.bidService.CreateBid(c.Request.Context(), service.CreateBidInput{
		TenderID:     tenderID,
		ContractorID: contractorID,
		Price:        req.Price,
		DeliveryTime: req.DeliveryTime,
		Comments:     req.Comments,
	})

	if err != nil {
		if errors.Is(err, service.ErrInvalidTender) {
			c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Tender is not open for bids"})
			return
		}
		pp.Print(err.Error())
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, bid)
}

// ListBids retrieves a list of bids for a specific tender.
//
// @Summary List bids for a tender
// @Description This endpoint retrieves a list of bids for a specified tender. The list can be filtered and sorted using query parameters.
// @Tags bids
// @Accept json
// @Produce json
// @Param tender_id path string true "Tender ID"
// @Param status query string false "Filter by bid status"
// @Param search query string false "Search bids by comments"
// @Param min_price query number false "Minimum bid price"
// @Param max_price query number false "Maximum bid price"
// @Param min_delivery_time query int false "Minimum delivery time"
// @Param max_delivery_time query int false "Maximum delivery time"
// @Param sort_by query string false "Sort by field"
// @Param sort_order query string false "Sort order (asc or desc)"
// @Success 200 {array} Bid "List of bids"
// @Failure 400 {object} ErrorResponse "Invalid tender ID"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Security ApiKeyAuth
// @Router /api/contractor/tenders/{tender_id}/bids [get]
func (h *BidHandler) ListBids(c *gin.Context) {
	tenderID, err := uuid.Parse(c.Param("tender_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid tender ID"})
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
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
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

// GetBidsByContractorID retrieves a list of bids made by a specific contractor.
//
// @Summary Get bids by contractor ID
// @Description This endpoint retrieves a list of bids made by a specific contractor.
// @Tags bids
// @Accept json
// @Produce json
// @Param contractor_id path string true "Contractor ID"
// @Success 200 {array} Bid "List of bids"
// @Failure 400 {object} ErrorResponse "Invalid contractor ID"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Security ApiKeyAuth
// @Router /api/contractor/bids [get]
func (h *BidHandler) GetBidsByContractorID(c *gin.Context) {

	contractorID, _ := c.Get("userId")

	contractorUUID, err := uuid.Parse(contractorID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid contractor ID"})
		return
	}

	bids, err := h.bidService.GetBidsByContractorID(c.Request.Context(), contractorUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, bids)
}

// GetBidsByClientID retrieves a list of bids made by a specific client.
//
// @Summary Get bids by client ID
// @Description This endpoint retrieves a list of bids made by a specific client for a specified tender.
// @Tags bids
// @Accept json
// @Produce json
// @Param tender_id path string true "Tender ID"
// @Success 200 {array} Bid "List of bids"
// @Failure 400 {object} ErrorResponse "Invalid tender ID or client ID"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Security ApiKeyAuth
// @Router /api/client/tenders/{tender_id}/bids [get]
func (h *BidHandler) GetBidsByClientID(c *gin.Context) {
	tenderID, err := uuid.Parse(c.Param("tender_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Message: "Tender not found or access denied"})
		return
	}

	clientID, _ := c.Get("userId")
	clientUUID, err := uuid.Parse(clientID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid client ID"})
		return
	}

	bids, err := h.bidService.GetBidsByClientID(c.Request.Context(), clientUUID, tenderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, bids)
}

// AwardBid awards a specific bid for a tender.
//
// @Summary Award a bid
// @Description This endpoint allows a client to award a specific bid for a specified tender.
// @Tags bids
// @Accept json
// @Produce json
// @Param tender_id path string true "Tender ID"
// @Param bid_id path string true "Bid ID"
// @Success 200 {object} Bid "Successfully awarded bid"
// @Failure 400 {object} ErrorResponse "Invalid tender ID or bid ID"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Security ApiKeyAuth
// @Router /api/client/tenders/{tender_id}/award/{bid_id} [post]
func (h *BidHandler) AwardBid(c *gin.Context) {
	tenderID, err := uuid.Parse(c.Param("tender_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Message: "Tender not found or access denied"})
		return
	}

	bidID, err := uuid.Parse(c.Param("bid_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Message: "Bid not found"})
		return
	}

	clientID, _ := c.Get("userId")
	clientUUID, err := uuid.Parse(clientID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid client ID"})
		return
	}

	err = h.bidService.AwardBid(c.Request.Context(), clientUUID, tenderID, bidID)
	if err != nil {
		if err.Error() == "unauthorized: client does not own the tender" || err.Error() == "bid not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{Message: "Tender not found or access denied"})
			return
		}

		if errors.Is(err, service.ErrInvalidTender) {
			c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Tender is not open for bids"})
			return
		}
		pp.Println(err.Error())
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bid awarded successfully"})
}

// DeleteBidByContractorID deletes a specific bid made by a contractor.
//
// @Summary Delete a bid by contractor ID
// @Description This endpoint allows a contractor to delete a specific bid they made.
// @Tags bids
// @Accept json
// @Produce json
// @Param bid_id path string true "Bid ID"
// @Success 200 {object} gin.H "Successfully deleted bid"
// @Failure 400 {object} ErrorResponse "Invalid bid ID or contractor ID"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Security ApiKeyAuth
// @Router /api/contractor/bids/{bid_id} [delete]
func (h *BidHandler) DeleteBidByContractorID(c *gin.Context) {
	bidID, err := uuid.Parse(c.Param("bid_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Message: "Bid not found or access denied"})
		return
	}

	contractorID, _ := c.Get("userId")
	contractorUUID, err := uuid.Parse(contractorID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid contractor ID"})
		return
	}

	err = h.bidService.DeleteBidByContractorID(c.Request.Context(), contractorUUID, bidID)
	if err != nil {
		if err.Error() == "unauthorized: contractor does not own the bid" || err.Error() == "bid not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{Message: "Bid not found or access denied"})
			return
		}
		if errors.Is(err, service.ErrInvalidContractor) {
			c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Bid not found or access denied"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bid deleted successfully"})
}
