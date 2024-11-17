package postgres

import (
	"context"
	"database/sql"
	"fmt"

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

func (r *BidRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Bid, error) {
	query := `
		SELECT id, tender_id, contractor_id, price, delivery_time, comments, status, created_at, updated_at
		FROM bids
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var b models.Bid
	err := row.Scan(
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
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &b, nil
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

func (r *BidRepo) ListByContractorID(ctx context.Context, contractorID uuid.UUID) ([]models.Bid, error) {
	query := `
		SELECT id, tender_id, contractor_id, price, delivery_time, comments, status, created_at, updated_at
		FROM bids
		WHERE contractor_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, contractorID)
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

func (r *BidRepo) Update(ctx context.Context, bid *models.Bid) error {
	query := `
		UPDATE bids
		SET tender_id = $2, contractor_id = $3, price = $4, delivery_time = $5, comments = $6, status = $7, updated_at = $8
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		bid.ID,
		bid.TenderID,
		bid.ContractorID,
		bid.Price,
		bid.DeliveryTime,
		bid.Comments,
		bid.Status,
		bid.UpdatedAt,
	)
	return err
}

func (r *BidRepo) ListByClientTenderID(ctx context.Context, clientID, tenderID uuid.UUID) ([]models.Bid, error) {
	query := `
		SELECT b.id, b.tender_id, b.contractor_id, b.price, b.delivery_time, b.comments, b.status, b.created_at, b.updated_at
		FROM bids b
		INNER JOIN tenders t ON b.tender_id = t.id
		WHERE t.client_id = $1 AND b.tender_id = $2
	`
	rows, err := r.db.QueryContext(ctx, query, clientID, tenderID)
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

func (r *BidRepo) AwardBidByTenderID(ctx context.Context, clientID, tenderID, bidID uuid.UUID) error {
	// Check if the tender belongs to the client
	var existingClientID uuid.UUID
	query := `
		SELECT client_id
		FROM tenders
		WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, tenderID).Scan(&existingClientID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("tender not found")
		}
		return err
	}

	if existingClientID != clientID {
		return fmt.Errorf("unauthorized: client does not own the tender")
	}

	// Check if the bid belongs to the tender
	var existingTenderID uuid.UUID
	query = `
		SELECT tender_id
		FROM bids
		WHERE id = $1
	`
	err = r.db.QueryRowContext(ctx, query, bidID).Scan(&existingTenderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("bid not found")
		}
		return err
	}

	if existingTenderID != tenderID {
		return fmt.Errorf("unauthorized: bid does not belong to the tender")
	}

	// Award the bid
	query = `
		UPDATE bids
		SET status = 'awarded'
		WHERE id = $1 AND tender_id = $2
	`
	_, err = r.db.ExecContext(ctx, query, bidID, tenderID)
	return err
}

func (r *BidRepo) DeleteByContractorID(ctx context.Context, contractorID, bidID uuid.UUID) error {
	// Check if the bid belongs to the contractor
	var existingContractorID uuid.UUID
	query := `
		SELECT contractor_id
		FROM bids
		WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, bidID).Scan(&existingContractorID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("bid not found")
		}
		return err
	}

	if existingContractorID != contractorID {
		return fmt.Errorf("unauthorized: contractor does not own the bid")
	}

	// Delete the bid
	query = `
		DELETE FROM bids
		WHERE id = $1 AND contractor_id = $2
	`
	_, err = r.db.ExecContext(ctx, query, bidID, contractorID)
	return err
}
