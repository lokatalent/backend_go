package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"path"
	"time"
)

// TimeRange represents available hours for a day
type TimeRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// Availability represents weekly availability schedule
type Availability struct {
	Monday    TimeRange `json:"monday"`
	Tuesday   TimeRange `json:"tuesday"`
	Wednesday TimeRange `json:"wednesday"`
	Thursday  TimeRange `json:"thursday"`
	Friday    TimeRange `json:"friday"`
	Saturday  TimeRange `json:"saturday"`
	Sunday    TimeRange `json:"sunday"`
}

// Service represents a service offered by a user
type UserService struct {
	ID              string       `json:"id"`
	UserID          string       `json:"user_id"`
	ServiceType     string       `json:"service_type"`
	ServiceDesc     string       `json:"service_desc"`
	RatePerHour     float64      `json:"rate_per_hour"`
	ExperienceYears int          `json:"experience_years"`
	Availability    Availability `json:"availability"`
	Address         string       `json:"address"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
}

type ServiceImage struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	ServiceType string `json:"service_type"`
	URL         string `json:"url"`
}

// Implement sql.Scanner interface for Availability to handle JSONB
func (a *Availability) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("error parsing service availability.")
	}

	return json.Unmarshal(bytes, &a)
}

// Implement driver.Valuer interface for Availability to save to JSONB
func (a Availability) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (s *ServiceImage) ServiceImagePath() string {
	return path.Join("profiles", "services", s.ServiceType, s.UserID, s.ID)
}
