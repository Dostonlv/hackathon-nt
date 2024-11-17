package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Dostonlv/hackathon-nt/internal/models"
	"github.com/Dostonlv/hackathon-nt/internal/repository"
	"github.com/google/uuid"
)

type TenderRepo struct {
	db *sql.DB
}

func NewTenderRepo(db *sql.DB) *TenderRepo {
	return &TenderRepo{db: db}
}

func (r *TenderRepo) Create(ctx context.Context, tender *models.Tender) error {
	query := `
		INSERT INTO tenders (
			id, client_id, title, description, deadline, budget, status, attachment, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.ExecContext(ctx, query,
		tender.ID,
		tender.ClientID,
		tender.Title,
		tender.Description,
		tender.Deadline,
		tender.Budget,
		tender.Status,
		tender.Attachment,
		tender.CreatedAt,
		tender.UpdatedAt,
	)
	return err
}

func (r *TenderRepo) List(ctx context.Context, filters repository.TenderFilters) ([]models.Tender, error) {
	query := `
		SELECT id, client_id, title, description, deadline, budget, status, attachment, created_at, updated_at
		FROM tenders
		WHERE ($1 IS NULL OR (title ILIKE $1 OR description ILIKE $1))
		AND ($2 IS NULL OR status = $2)
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, nullableString(filters.Search), nullableString(filters.Status))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenders []models.Tender
	for rows.Next() {
		var t models.Tender
		err := rows.Scan(
			&t.ID,
			&t.ClientID,
			&t.Title,
			&t.Description,
			&t.Deadline,
			&t.Budget,
			&t.Status,
			&t.Attachment,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tenders = append(tenders, t)
	}
	return tenders, nil
}

func nullableString(s string) interface{} {
	if s == "" {
		return nil
	}
	return "%" + s + "%"
}

func (r *TenderRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Tender, error) {
	query := `
		SELECT id, client_id, title, description, deadline, budget, status, attachment, created_at, updated_at
		FROM tenders
		WHERE id = $1
	`

	var t models.Tender
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID,
		&t.ClientID,
		&t.Title,
		&t.Description,
		&t.Deadline,
		&t.Budget,
		&t.Status,
		&t.Attachment,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &t, nil
}

func (r *TenderRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	exists, err := r.exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return repository.ErrNotFound
	}

	query := `
		UPDATE tenders
		SET status = $1, updated_at = $2
		WHERE id = $3
	`
	_, err = r.db.ExecContext(ctx, query, status, sql.NullTime{Time: time.Now(), Valid: true}, id)
	return err
}

func (r *TenderRepo) Delete(ctx context.Context, id uuid.UUID) error {
	exists, err := r.exists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return repository.ErrNotFound
	}

	query := `
		DELETE FROM tenders
		WHERE id = $1
	`
	_, err = r.db.ExecContext(ctx, query, id)
	return err
}

func (r *TenderRepo) ListByClientID(ctx context.Context, clientID uuid.UUID) ([]models.Tender, error) {
	query := `
		SELECT id, client_id, title, description, deadline, budget, status, attachment, created_at, updated_at
		FROM tenders
		WHERE client_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenders []models.Tender
	for rows.Next() {
		var t models.Tender
		err := rows.Scan(
			&t.ID,
			&t.ClientID,
			&t.Title,
			&t.Description,
			&t.Deadline,
			&t.Budget,
			&t.Status,
			&t.Attachment,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tenders = append(tenders, t)
	}
	return tenders, nil
}

func (r *TenderRepo) exists(ctx context.Context, id uuid.UUID) (bool, error) {
	query := `SELECT 1 FROM tenders WHERE id = $1`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
