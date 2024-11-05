package repository

import "github.com/lokatalent/backend_go/internal/models"

type NotificationRepository interface {
	Create(notification *models.Notification) error
	GetForUser(filter models.NotificationFilter) ([]models.Notification, error)
	MarkAsRead(userID, notificationID string) error
	MarkAllAsRead(userID string) error
	CountUnseen(userID string) (int64, error)
}
