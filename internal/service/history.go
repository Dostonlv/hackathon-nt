package service

import (
	"github.com/Dostonlv/hackathon-nt/internal/models"
	"github.com/Dostonlv/hackathon-nt/internal/repository"
	"github.com/google/uuid"
)

// HistoryService represents the service for retrieving history data.
type HistoryService struct {
	historyRepo repository.HistoryRepository
}

// NewHistoryService creates a new instance of HistoryService.
func NewHistoryService(historyRepo repository.HistoryRepository) *HistoryService {
	return &HistoryService{
		historyRepo: historyRepo,
	}
}

// GetTenderHistory retrieves the tender history for a specific user.
func (s *HistoryService) GetTenderHistory(userID uuid.UUID) ([]models.Tender, error) {
	return s.historyRepo.GetTenderHistory(userID)
}

// GetBidHistory retrieves the bid history for a specific contractor.
func (s *HistoryService) GetBidHistory(userID uuid.UUID) ([]models.Bid, error) {
	return s.historyRepo.GetBidHistory(userID)
}
