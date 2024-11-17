package models

import (
	"time"

	"github.com/google/uuid"
)

type TenderStatus string

func (s TenderStatus) IsValid() bool {
	switch s {
	case TenderStatusOpen, TenderStatusClosed, TenderStatusAwarded:
		return true
	}
	return false
}

const (
	TenderStatusOpen    TenderStatus = "open"
	TenderStatusClosed  TenderStatus = "closed"
	TenderStatusAwarded TenderStatus = "awarded"
)

type Tender struct {
	ID          uuid.UUID    `json:"id" db:"id"`
	ClientID    uuid.UUID    `json:"client_id" db:"client_id"`
	Title       string       `json:"title" db:"title"`
	Description string       `json:"description" db:"description"`
	Deadline    time.Time    `json:"deadline" db:"deadline"`
	Budget      float64      `json:"budget" db:"budget"`
	Status      TenderStatus `json:"status" db:"status"`
	Attachment  *string      `json:"attachment,omitempty" db:"attachment"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
}
