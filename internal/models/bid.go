package models

import (
	"time"

	"github.com/google/uuid"
)

type Bid struct {
	ID           uuid.UUID `json:"id" db:"id"`
	TenderID     uuid.UUID `json:"tender_id" db:"tender_id"`
	ContractorID uuid.UUID `json:"contractor_id" db:"contractor_id"`
	Price        float64   `json:"price" db:"price"`
	DeliveryTime int       `json:"delivery_time" db:"delivery_time"`
	Comments     string    `json:"comments" db:"comments"`
	Status       string    `json:"status" db:"status"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
