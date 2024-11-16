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
			string(models.TenderStatusClosed):  true,
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
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrTenderNotFound
		}
		return nil, err
	}

	return tender, nil
}

type UpdateTenderInput struct {
	ID          uuid.UUID
	ClientID    uuid.UUID // for authorization
	Title       *string
	Description *string
	Deadline    *time.Time
	Budget      *float64
	FileURL     *string
	Status      *string
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

	// Check authorization
	if tender.ClientID != input.ClientID {
		return nil, ErrUnauthorized
	}

	// Update fields if provided
	if input.Title != nil {
		if *input.Title == "" {
			return nil, errors.Join(ErrInvalidInput, errors.New("title cannot be empty"))
		}
		tender.Title = *input.Title
	}

	if input.Description != nil {
		if *input.Description == "" {
			return nil, errors.Join(ErrInvalidInput, errors.New("description cannot be empty"))
		}
		tender.Description = *input.Description
	}

	if input.Deadline != nil {
		if input.Deadline.Before(time.Now()) {
			return nil, errors.Join(ErrInvalidInput, errors.New("deadline must be in the future"))
		}
		tender.Deadline = *input.Deadline
	}

	if input.Budget != nil {
		if *input.Budget <= 0 {
			return nil, errors.Join(ErrInvalidInput, errors.New("budget must be greater than zero"))
		}
		tender.Budget = *input.Budget
	}

	if input.FileURL != nil {
		tender.FileURL = input.FileURL
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
	err = s.repo.Update(ctx, tender)
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
