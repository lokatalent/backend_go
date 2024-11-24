package models

import (
	"database/sql"
	"time"
)

type UserWallet struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Balance   float64   `json:"balance"`
	Debits    float64   `json:"debits"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Payment struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	BookingID  sql.NullString `json:"booking_id"`
	Amount     float64        `json:"amount"`
	PaymentRef string         `json:"payment_ref"`
	Status     string         `json:"status"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

type PaymentAccessCodes struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	AccessCode string    `json:"access_code"`
	CreatedAt  time.Time `json:"created_at"`
}

type PaymentRecipientCode struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	RecipientCode string    `json:"recipient_code"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
