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
	ErrInvalidInput   = errors.New("invalid input parameters")
	ErrTenderNotFound = errors.New("tender not found")
	ErrUnauthorized   = errors.New("unauthorized action")
)

type TenderService struct {
	repo repository.TenderRepository
}

func NewTenderService(repo repository.TenderRepository) *TenderService {
	if repo == nil {
		panic("tender repository cannot be nil")
	}
	return &TenderService{
		repo: repo,
	}
}

type CreateTenderInput struct {
	ClientID    uuid.UUID
	Title       string
	Description string
	Deadline    time.Time
	Budget      float64
	FileURL     *string
}

func (s *TenderService) validateCreateTenderInput(input CreateTenderInput) error {
	if input.ClientID == uuid.Nil {
		return errors.New("client ID is required")
	}
	if input.Title == "" {
		return errors.New("title is required")
	}
	if input.Description == "" {
		return errors.New("description is required")
	}
	if input.Deadline.Before(time.Now()) {
		return errors.New("deadline must be in the future")
	}
	if input.Budget <= 0 {
		return errors.New("budget must be greater than zero")
	}
	return nil
}

// CreateTender creates a new tender
func (s *TenderService) CreateTender(ctx context.Context, input CreateTenderInput) (*models.Tender, error) {
	if err := s.validateCreateTenderInput(input); err != nil {
		return nil, errors.Join(ErrInvalidInput, err)
	}

	tender := &models.Tender{
		ID:          uuid.New(),
		ClientID:    input.ClientID,
		Title:       input.Title,
		Description: input.Description,
		Deadline:    input.Deadline,
		Budget:      input.Budget,
		FileURL:     input.FileURL,
		Status:      models.TenderStatusOpen,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.repo.Create(ctx, tender)
	if err != nil {
		return nil, err
	}

	return tender, nil
}



func (s *TenderService) ListTenders(ctx context.Context, filters repository.TenderFilters) ([]models.Tender, error) {
	// Validate status if provided
	if filters.Status != "" {
		validStatuses := map[string]bool{
			string(models.TenderStatusOpen):    true,
			string(models.TenderStatusClosed): true,
			string(models.TenderStatusAwarded): true,
		}
		if !validStatuses[filters.Status] {
			return nil, errors.Join(ErrInvalidInput, errors.New("invalid status filter"))
		}
	}

	return s.repo.List(ctx, repository.TenderFilters{
		Status: filters.Status,
		Search: filters.Search,
	})
}


func (s *TenderService) GetTenderByID(ctx context.Context, id uuid.UUID) (*models.Tender, error) {
	if id == uuid.Nil {
		return nil, errors.Join(ErrInvalidInput, errors.New("invalid tender ID"))
	}

	tender, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return tender, nil
}