package repository

import "github.com/lokatalent/backend_go/internal/models"

type BookingRepository interface {
	Create(booking *models.Booking) error
	GetByID(id string) (models.Booking, error)
	GetAll(filter models.BookingFilter) ([]models.Booking, error)

	UpdateStatus(id, status string) (models.Booking, error)
	AssignProvider(id, userID string) (models.Booking, error)

	RejectBooking(id, userID string) error
	CheckRejected(id, userID string) (bool, error)

	MatchServices(booking *models.Booking, filter models.ServiceFilter) ([]models.UserService, error)
}
