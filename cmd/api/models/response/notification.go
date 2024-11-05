package response

import (
	"time"

	"github.com/lokatalent/backend_go/internal/models"
)

type NotificationResponse struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	UserID    string    `json:"user_id"`
	BookingID string    `json:"booking_id"`
	Seen      bool      `json:"seen"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

func NotificationResponseFromModel(notification models.Notification) NotificationResponse {
	response := NotificationResponse{
		ID:        notification.ID,
		Type:      notification.Type,
		UserID:    notification.UserID,
		Seen:      notification.Seen,
		Message:   notification.Message,
		CreatedAt: notification.CreatedAt,
	}

	if notification.BookingID.Valid {
		response.BookingID = notification.BookingID.String
	}

	return response
}
