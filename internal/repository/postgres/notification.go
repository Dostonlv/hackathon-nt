package postgres

import (
	"context"
	"database/sql"

	"github.com/Dostonlv/hackathon-nt/internal/models"
	"github.com/google/uuid"
)

type NotificationRepo struct {
	db *sql.DB
}

func NewNotificationRepo(db *sql.DB) *NotificationRepo {
	return &NotificationRepo{db: db}
}

func (r *NotificationRepo) Create(ctx context.Context, notification *models.Notification) error {
	query := `
        INSERT INTO notifications (id, user_id, message, relation_id, type, read, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
	_, err := r.db.ExecContext(ctx, query,
		notification.ID,
		notification.UserID,
		notification.Message,
		notification.RelationID,
		notification.Type,
		notification.Read,
		notification.CreatedAt,
	)
	return err
}

func (r *NotificationRepo) ListByUserID(ctx context.Context, userID uuid.UUID) ([]models.Notification, error) {
	query := `
        SELECT id, user_id, message, relation_id, type, read, created_at
        FROM notifications
        WHERE user_id = $1
        ORDER BY created_at DESC
    `

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var n models.Notification
		err := rows.Scan(
			&n.ID,
			&n.UserID,
			&n.Message,
			&n.RelationID,
			&n.Type,
			&n.Read,
			&n.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}
	return notifications, nil
}

func (r *NotificationRepo) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE notifications SET read = true WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
