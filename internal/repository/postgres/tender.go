package postgres

import (
	"context"
	"database/sql"

	"github.com/Dostonlv/hackathon-nt/internal/models"
	"github.com/Dostonlv/hackathon-nt/internal/repository"
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
            id, client_id, title, description, deadline, budget, status, file_url, created_at, updated_at
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
		tender.FileURL,
		tender.CreatedAt,
		tender.UpdatedAt,
	)
	return err
}

func (r *TenderRepo) List(ctx context.Context, filters repository.TenderFilters) ([]models.Tender, error) {
	query := `
        SELECT id, client_id, title, description, deadline, budget, status, file_url, created_at, updated_at
        FROM tenders
        WHERE ($1::text IS NULL OR status = $1)
        AND ($2::text IS NULL OR 
            title ILIKE '%' || $2 || '%' OR 
            description ILIKE '%' || $2 || '%')
        ORDER BY created_at DESC
    `

	rows, err := r.db.QueryContext(ctx, query, filters.Status, filters.Search)
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
			&t.FileURL,
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
