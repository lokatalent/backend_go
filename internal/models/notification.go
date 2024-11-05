package models

import (
	"database/sql"
	"time"
)

type Notification struct {
	ID        string
	Type      string
	UserID    string
	BookingID sql.NullString
	Message   string
	Seen      bool
	CreatedAt time.Time
}
