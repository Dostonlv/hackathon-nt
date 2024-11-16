package service

import (
	"context"
	"errors"
	"time"

	"github.com/Dostonlv/hackathon-nt/internal/models"
	"github.com/Dostonlv/hackathon-nt/internal/repository"
	"github.com/google/uuid"
)

var (
	ErrInvalidTender     = errors.New("invalid tender")
	ErrInvalidContractor = errors.New("invalid contractor")
	ErrBidNotFound       = errors.New("bid not found")
)

type CreateBidInput struct {
	TenderID     uuid.UUID
	ContractorID uuid.UUID
	Price        float64
	DeliveryTime int
	Comments     string
}

type BidService struct {
	bidRepo repository.BidRepository
}

func NewBidService(
	bidRepo repository.BidRepository,
) *BidService {
	return &BidService{
		bidRepo: bidRepo,
	}
}

func (s *BidService) CreateBid(ctx context.Context, input CreateBidInput) (*models.Bid, error) {
	bid := &models.Bid{
		ID:           uuid.New(),
		TenderID:     input.TenderID,
		ContractorID: input.ContractorID,
		Price:        input.Price,
		DeliveryTime: input.DeliveryTime,
		Comments:     input.Comments,
		Status:       "pending", // Initial status
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.bidRepo.Create(ctx, bid); err != nil {
		return nil, err
	}

	return bid, nil
}

func (s *BidService) ListBids(ctx context.Context, tenderID uuid.UUID, filters repository.BidFilters) ([]models.Bid, error) {
	// Get bids with filters
	bids, err := s.bidRepo.ListByTenderID(ctx, tenderID, filters)
	if err != nil {
		return nil, err
	}

	return bids, nil
}

// Additional helper methods that might be useful

func (s *BidService) GetBidByID(ctx context.Context, bidID uuid.UUID) (*models.Bid, error) {
	bid, err := s.bidRepo.GetByID(ctx, bidID)
	if err != nil {
		return nil, err
	}
	if bid == nil {
		return nil, ErrBidNotFound
	}
	return bid, nil
}

func (s *BidService) UpdateBidStatus(ctx context.Context, bidID uuid.UUID, status string) error {
	bid, err := s.GetBidByID(ctx, bidID)
	if err != nil {
		return err
	}

	bid.Status = status
	bid.UpdatedAt = time.Now()

	return s.bidRepo.Update(ctx, bid)
}
