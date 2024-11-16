package postgres

import (
	"context"
	"database/sql"

	"github.com/Dostonlv/hackathon-nt/internal/models"
	"github.com/Dostonlv/hackathon-nt/internal/repository"
	"github.com/google/uuid"
)

type BidRepo struct {
	db *sql.DB
}

func NewBidRepo(db *sql.DB) *BidRepo {
	return &BidRepo{db: db}
}

func (r *BidRepo) Create(ctx context.Context, bid *models.Bid) error {
	query := `
        INSERT INTO bids (
            id, tender_id, contractor_id, price, delivery_time, comments, status, created_at, updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `
	_, err := r.db.ExecContext(ctx, query,
		bid.ID,
		bid.TenderID,
		bid.ContractorID,
		bid.Price,
		bid.DeliveryTime,
		bid.Comments,
		bid.Status,
		bid.CreatedAt,
		bid.UpdatedAt,
	)
	return err
}

func (r *BidRepo) ListByTenderID(ctx context.Context, tenderID uuid.UUID, filters repository.BidFilters) ([]models.Bid, error) {
	query := `
        SELECT id, tender_id, contractor_id, price, delivery_time, comments, status, created_at, updated_at
        FROM bids
        WHERE tender_id = $1
        AND ($2::float8 IS NULL OR price >= $2)
        AND ($3::float8 IS NULL OR price <= $3)
    `

	if filters.SortBy != "" {
		query += " ORDER BY " + filters.SortBy
		if filters.SortOrder == "desc" {
			query += " DESC"
		}
	}

	rows, err := r.db.QueryContext(ctx, query, tenderID, filters.MinPrice, filters.MaxPrice)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bids []models.Bid
	for rows.Next() {
		var b models.Bid
		err := rows.Scan(
			&b.ID,
			&b.TenderID,
			&b.ContractorID,
			&b.Price,
			&b.DeliveryTime,
			&b.Comments,
			&b.Status,
			&b.CreatedAt,
			&b.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		bids = append(bids, b)
	}
	return bids, nil
}
