package repository

import (
	"context"
	"database/sql"
	"pet-project/pkg/model"

	"github.com/lib/pq"
)

type PostgresNotificationRepository struct {
	DB *sql.DB
}

type NotificationRepository interface {
	Create(ctx context.Context, notif *model.Notification) error
	GetByUserID(ctx context.Context, userID int, limit, offset int) ([]model.Notification, error)
	MarkAsRead(ctx context.Context, userID int, notifID []int) error
	CountUnread(ctx context.Context, userID int) (int, error)
}

func (r *PostgresNotificationRepository) Create(ctx context.Context, notif *model.Notification) error {
	query := `INSERT INTO notification (user_id, type, message, is_read, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return r.DB.QueryRowContext(ctx, query, notif.UserID, notif.Type, notif.Message, notif.IsRead, notif.CreatedAt).Scan(&notif.ID)
}

func (r *PostgresNotificationRepository) GetByUserID(ctx context.Context, userID int, limit, offset int) ([]model.Notification, error) {
	query := `SELECT id, user_id, type, message, is_read, created_at FROM notification WHERE user_id = $1 LIMIT $2 OFFSET $3`
	rows, err := r.DB.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notifications := []model.Notification{}

	for rows.Next() {
		var notif model.Notification
		err := rows.Scan(
			&notif.ID,
			&notif.UserID,
			&notif.Type,
			&notif.Message,
			&notif.IsRead,
			&notif.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notif)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r *PostgresNotificationRepository) MarkAsRead(ctx context.Context, userID int, notifID []int) error {
	query := `UPDATE notification SET is_read = TRUE WHERE user_id = $1 AND id = ANY($2)`
	_, err := r.DB.ExecContext(ctx, query, userID, pq.Array(notifID))
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresNotificationRepository) CountUnread(ctx context.Context, userID int) (int, error) {
	query := `SELECT COUNT(id) FROM notification WHERE is_read = FALSE AND user_id = $1`
	var cnt int
	err := r.DB.QueryRowContext(ctx, query, userID).Scan(&cnt)
	if err != nil {
		return 0, err
	}
	return cnt, nil
}
