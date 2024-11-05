package postgres

import (
	"context"
	"database/sql"
	// "errors"

	"github.com/google/uuid"

	"github.com/lokatalent/backend_go/internal/models"
	"github.com/lokatalent/backend_go/internal/repository"
)

type notificationImplementation struct {
	DB *sql.DB
}

func NewNotificationImplementation(db *sql.DB) repository.NotificationRepository {
	return &notificationImplementation{DB: db}
}

func (n *notificationImplementation) Create(notification *models.Notification) error {
	if notification.ID == "" {
		notification.ID = uuid.NewString()
	}
	stmt := `
    INSERT INTO notifications (
        id,
        type,
        user_id,
        booking_id,
        message
    ) VALUES (
        $1, $2, $3, $4, $5
    ) RETURNING
        id,
        type,
        user_id,
        booking_id,
        message,
        seen,
        created_at;
    `
	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	err := n.DB.QueryRowContext(
		ctx,
		stmt,
		notification.ID,
		notification.Type,
		notification.UserID,
		notification.BookingID,
		notification.Message,
	).Scan(
		&notification.ID,
		&notification.Type,
		&notification.UserID,
		&notification.BookingID,
		&notification.Message,
		&notification.Seen,
		&notification.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (n *notificationImplementation) GetForUser(filter models.NotificationFilter) ([]models.Notification, error) {
	stmt := `
    SELECT
        id,
        type,
        user_id,
        booking_id,
        message,
        seen,
        created_at
    FROM notifications
    WHERE
        ($1 = '' OR type = $1) AND
        ($2 = '' OR user_id = $2::UUID) AND
        (
            NULLIF($3, '') IS NULL 
            OR (
                CASE 
                    WHEN NULLIF($3, '') IS NULL THEN booking_id IS NULL
                    ELSE booking_id = $3::UUID
                END
            )
        ) -- AND
        -- ($3 = '' OR booking_id = $3::UUID) AND
        -- seen = $4
    ORDER BY seen, created_at DESC
    LIMIT $4 OFFSET $5;
    `
	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	rows, err := n.DB.QueryContext(
		ctx,
		stmt,
		filter.Type,
		filter.UserID,
		filter.BookingID,
		// filter.Seen,
		filter.Limit,
		filter.Offset(),
	)
	if err != nil {
		return nil, err
	}

	notifications := []models.Notification{}
	for rows.Next() {
		notification := models.Notification{}
		err := rows.Scan(
			&notification.ID,
			&notification.Type,
			&notification.UserID,
			&notification.BookingID,
			&notification.Message,
			&notification.Seen,
			&notification.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func (n *notificationImplementation) MarkAllAsRead(userID string) error {
	stmt := `
	UPDATE notifications
	SET seen = true
	WHERE user_id = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	_, err := n.DB.ExecContext(ctx, stmt, userID)
	if err != nil {
		return err
	}
	return nil
}
func (n notificationImplementation) MarkAsRead(userID, notificationID string) error {
	stmt := `
	UPDATE notifications
	SET seen = true
	WHERE id=$1 AND user_id = $2;
	`

	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	_, err := n.DB.ExecContext(ctx, stmt, notificationID, userID)
	if err != nil {
		return err
	}
	return nil
}

// CountForUser returns the number of unseen notifications for a user.
func (n notificationImplementation) CountUnseen(userID string) (int64, error) {
	stmt := `
	SELECT COUNT(*)
	FROM notifications
	WHERE user_id = $1 AND seen = FALSE
	`

	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	var count int64
	err := n.DB.QueryRowContext(ctx, stmt, userID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
