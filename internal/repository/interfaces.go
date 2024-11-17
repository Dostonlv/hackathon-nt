package repository

import (
	"context"

	"github.com/Dostonlv/hackathon-nt/internal/models"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	// Update(ctx context.Context, user *models.User) error
}

type TenderRepository interface {
	Create(ctx context.Context, tender *models.Tender) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Tender, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByClientID(ctx context.Context, clientID uuid.UUID) ([]models.Tender, error)
	List(ctx context.Context, filters TenderFilters) ([]models.Tender, error)
	GetClientIDByTenderID(ctx context.Context, tenderID uuid.UUID) (uuid.UUID, error)
}

type BidRepository interface {
	Create(ctx context.Context, bid *models.Bid) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Bid, error)
	ListByTenderID(ctx context.Context, tenderID uuid.UUID, filters BidFilters) ([]models.Bid, error)
	ListByContractorID(ctx context.Context, contractorID uuid.UUID) ([]models.Bid, error)
	Update(ctx context.Context, bid *models.Bid) error
	ListByClientTenderID(ctx context.Context, clientID, tenderID uuid.UUID) ([]models.Bid, error)
	AwardBidByTenderID(ctx context.Context, clientID, tenderID, bidID uuid.UUID) error
	DeleteByContractorID(ctx context.Context, contractorID, bidID uuid.UUID) error
}

type NotificationRepository interface {
	Create(ctx context.Context, notification *models.Notification) error
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]models.Notification, error)
	MarkAsRead(ctx context.Context, id uuid.UUID) error
}

type HistoryRepository interface {
	GetTenderHistory(userID uuid.UUID) ([]models.Tender, error)
	GetBidHistory(userID uuid.UUID) ([]models.Bid, error)
}

type TenderFilters struct {
	Status string
	Search string
}

type BidFilters struct {
	Status          string
	Search          string
	Price           *float64
	DeliveryTime    *int
	MinDeliveryTime *int
	MaxDeliveryTime *int
	MinPrice        *float64
	MaxPrice        *float64
	MinDelivery     *int
	MaxDelivery     *int
	SortBy          string
	SortOrder       string
}
