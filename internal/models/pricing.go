package models

import "time"

type ServiceCommission struct {
	ID         string    `json:"id"`
	Percentage int       `json:"percentage"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ServicePricing struct {
	ID          string    `json:"id"`
	ServiceType string    `json:"service_type"`
	RatePerHour float64   `json:"rate_per_hour"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
