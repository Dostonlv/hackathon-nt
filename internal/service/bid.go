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
	bidRepo    repository.BidRepository
	tenderRepo repository.TenderRepository
}

func NewBidService(bidRepo repository.BidRepository, tenderRepo repository.TenderRepository) *BidService {
	return &BidService{
		bidRepo:    bidRepo,
		tenderRepo: tenderRepo,
	}
}

func (s *BidService) CreateBid(ctx context.Context, input CreateBidInput) (*models.Bid, error) {
	// Check if tender exists
	tender, err := s.tenderRepo.GetByID(ctx, input.TenderID)
	if err != nil {
		return nil, err
	}

	if tender.Status != models.TenderStatusOpen {
		return nil, ErrInvalidTender
	}

	bid := &models.Bid{
		ID:           uuid.New(),
		TenderID:     input.TenderID,
		ContractorID: input.ContractorID,
		Price:        input.Price,
		DeliveryTime: input.DeliveryTime,
		Comments:     input.Comments,
		Status:       "open",
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
func (s *BidService) GetBidsByContractorID(ctx context.Context, contractorID uuid.UUID) ([]models.Bid, error) {
	bids, err := s.bidRepo.ListByContractorID(ctx, contractorID)
	if err != nil {
		return nil, err
	}
	return bids, nil
}

func (s *BidService) GetBidsByClientID(ctx context.Context, clientID, tenderID uuid.UUID) ([]models.Bid, error) {
	bids, err := s.bidRepo.ListByClientTenderID(ctx, clientID, tenderID)
	if err != nil {
		return nil, err
	}
	return bids, nil
}

func (s *BidService) AwardBid(ctx context.Context, clientID, tenderID, bidID uuid.UUID) error {
	// Check if tender exists
	tender, err := s.tenderRepo.GetByID(ctx, tenderID)
	if err != nil {
		return err
	}

	if tender.Status != models.TenderStatusOpen {
		return ErrInvalidTender
	}

	// Check if bid exists
	bid, err := s.GetBidByID(ctx, bidID)
	if err != nil {
		return err
	}

	if bid.TenderID != tenderID {
		return ErrInvalidTender
	}

	// Award the bid
	err = s.bidRepo.AwardBidByTenderID(ctx, clientID, tenderID, bidID)
	if err != nil {
		return err
	}

	return nil
}

func (s *BidService) DeleteBidByContractorID(ctx context.Context, contractorID, bidID uuid.UUID) error {
	// Check if bid exists
	bid, err := s.GetBidByID(ctx, bidID)
	if err != nil {
		return err
	}

	if bid.ContractorID != contractorID {
		return ErrInvalidContractor
	}

	// Delete the bid
	err = s.bidRepo.DeleteByContractorID(ctx, contractorID, bidID)
	if err != nil {
		return err
	}

	return nil
}
