package response

import (
	"time"

	"github.com/lokatalent/backend_go/internal/models"
)

type DistanceEstimation struct {
	Distance string `json:"distance"`
	Duration string `json:"duration"`
}

type ServiceDistanceResponse struct {
	ServiceID       string
	ProviderID      string
	ServiceType     string
	ServiceDesc     string
	ExperienceYears int
	DistanceEstimation
}

type BookingResponse struct {
	ID            string    `json:"id"`
	RequesterID   string    `json:"requester_id"`
	ProviderID    string    `json:"provider_id"`
	RequesterAddr string    `json:"requester_addr"`
	ServiceType   string    `json:"service_type"`
	BookingType   string    `json:"booking_type"`
	ServiceDesc   string    `json:"service_desc"`
	StartTime     string    `json:"start_time"`
	EndTime       string    `json:"end_time"`
	StartDate     string    `json:"start_date"`
	EndDate       string    `json:"end_date"`
	TotalPrice    float64   `json:"total_price"`
	ActualPrice   float64   `json:"actual_price"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func BookingResponseFromModel(booking models.Booking) BookingResponse {
	response := BookingResponse{
		ID:            booking.ID,
		RequesterID:   booking.RequesterID,
		RequesterAddr: booking.RequesterAddr,
		ServiceType:   booking.ServiceType,
		BookingType:   booking.BookingType,
		ServiceDesc:   booking.ServiceDesc,
		TotalPrice:    booking.TotalPrice,
		ActualPrice:   booking.ActualPrice,
		Status:        booking.Status,
		CreatedAt:     booking.CreatedAt,
		UpdatedAt:     booking.UpdatedAt,
	}

	if booking.ProviderID.Valid {
		response.ProviderID = booking.ProviderID.String
	}

	response.StartTime = booking.StartTime.Format(time.Kitchen)
	response.EndTime = booking.EndTime.Format(time.Kitchen)
	response.StartDate = booking.StartDate.Format(time.DateOnly)
	response.EndDate = booking.EndDate.Format(time.DateOnly)

	return response
}

func BookingServiceProviderResponse(service models.UserService, distance DistanceEstimation) ServiceDistanceResponse {
	response := ServiceDistanceResponse{
		ServiceID:       service.ID,
		ProviderID:      service.UserID,
		ServiceType:     service.ServiceType,
		ServiceDesc:     service.ServiceDesc,
		ExperienceYears: service.ExperienceYears,
		DistanceEstimation: DistanceEstimation{
			Distance: distance.Distance,
			Duration: distance.Duration,
		},
	}

	return response
}
