package models

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID         uuid.UUID `json:"id" db:"id"`
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	Message    string    `json:"message" db:"message"`
	RelationID uuid.UUID `json:"relation_id,omitempty" db:"relation_id"`
	Type       string    `json:"type" db:"type"`
	Read       bool      `json:"read" db:"read"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
