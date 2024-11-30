package models

import (
	"database/sql"
	"time"
)

type Booking struct {
	ID            string         `json:"id"`
	RequesterID   string         `json:"requester_id"`
	ProviderID    sql.NullString `json:"provider_id"`
	RequesterAddr string         `json:"requester_addr"`
	ServiceType   string         `json:"service_type"`
	BookingType   string         `json:"booking_type"`
	ServiceDesc   string         `json:"service_desc"`
	StartTime     time.Time      `json:"start_time"`
	EndTime       time.Time      `json:"end_time"`
	StartDate     time.Time      `json:"start_date"`
	EndDate       time.Time      `json:"end_date"`
	TotalPrice    float64        `json:"total_price"`
	ActualPrice   float64        `json:"actual_price"`
	Status        string         `json:"status"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}
