package postgres

import (
	"database/sql"
	"fmt"

	"github.com/Dostonlv/hackathon-nt/internal/models"
	"github.com/google/uuid"
)

// Tender represents a tender history entry.
type HistoryRepo struct {
	db *sql.DB
}

// NewHistoryRepo creates a new instance of HistoryRepo.
func NewHistoryRepo(db *sql.DB) *HistoryRepo {
	return &HistoryRepo{
		db: db,
	}
}

// GetTenderHistory retrieves the tender history for a specific user.
func (h *HistoryRepo) GetTenderHistory(userID uuid.UUID) ([]models.Tender, error) {
	query := `
		SELECT t.id, t.title, t.description, t.status, t.created_at, t.updated_at
		FROM tenders t
		INNER JOIN users u ON t.user_id = u.id
		WHERE u.id = $1
	`

	rows, err := h.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	var tenders []models.Tender
	for rows.Next() {
		var tender models.Tender
		err := rows.Scan(
			&tender.ID,
			&tender.Title,
			&tender.Description,
			&tender.Status,
			&tender.CreatedAt,
			&tender.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		tenders = append(tenders, tender)
	}

	return tenders, nil

}

// GetBidHistory retrieves the bid history for a specific contractor.
func (h *HistoryRepo) GetBidHistory(userID uuid.UUID) ([]models.Bid, error) {
	query := `
		SELECT b.id, b.tender_id, b.contractor_id, b.price, b.status, b.created_at, b.updated_at
		FROM bids b
		INNER JOIN contractors c ON b.contractor_id = c.id
		WHERE c.id = $1
	`

	rows, err := h.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	var bids []models.Bid
	for rows.Next() {
		var bid models.Bid
		err := rows.Scan(
			&bid.ID,
			&bid.TenderID,
			&bid.ContractorID,
			&bid.Price,
			&bid.Status,
			&bid.CreatedAt,
			&bid.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		bids = append(bids, bid)
	}

	return bids, nil
}
