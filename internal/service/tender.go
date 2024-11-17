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
	Attachment  *string
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
	tender := &models.Tender{
		ID:          uuid.New(),
		ClientID:    input.ClientID,
		Title:       input.Title,
		Description: input.Description,
		Deadline:    input.Deadline,
		Budget:      input.Budget,
		Attachment:  input.Attachment,
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

func (s *TenderService) ListTenders(ctx context.Context, clientID uuid.UUID) ([]models.Tender, error) {
	return s.repo.ListByClientID(ctx, clientID)
}

func (s *TenderService) GetTenderByID(ctx context.Context, id uuid.UUID) (*models.Tender, error) {
	if id == uuid.Nil {
		return nil, errors.Join(ErrInvalidInput, errors.New("tender not found"))
	}

	tender, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrTenderNotFound
		}
		return nil, err
	}

	return tender, nil
}

type UpdateTenderInput struct {
	ID     uuid.UUID
	Status *string
}

func (s *TenderService) UpdateTender(ctx context.Context, input UpdateTenderInput) (*models.Tender, error) {
	// Validate input
	if input.ID == uuid.Nil {
		return nil, errors.Join(ErrInvalidInput, errors.New("invalid tender ID"))
	}

	// Get existing tender
	tender, err := s.repo.GetByID(ctx, input.ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrTenderNotFound
		}
		return nil, err
	}

	if input.Status != nil {
		newStatus := models.TenderStatus(*input.Status)
		if !newStatus.IsValid() {
			return nil, errors.Join(ErrInvalidInput, errors.New("invalid status"))
		}
		tender.Status = newStatus
	}

	tender.UpdatedAt = time.Now()

	// Save updates
	err = s.repo.UpdateStatus(ctx, input.ID, string(tender.Status))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrTenderNotFound
		}
		return nil, err
	}
	return tender, nil
}

// DeleteTender deletes a tender
func (s *TenderService) DeleteTender(ctx context.Context, tenderID, clientID uuid.UUID) error {
	if tenderID == uuid.Nil {
		return errors.Join(ErrInvalidInput, errors.New("invalid tender ID"))
	}

	// Get existing tender
	tender, err := s.repo.GetByID(ctx, tenderID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrTenderNotFound
		}
		return err
	}

	// Check authorization
	if tender.ClientID != clientID {
		return ErrUnauthorized
	}

	return s.repo.Delete(ctx, tenderID)
}

func (s *TenderService) ListTendersFiltering(ctx context.Context, filters repository.TenderFilters) ([]models.Tender, error) {
	return s.repo.List(ctx, filters)
}
