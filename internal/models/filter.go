package models

// import "time"

type Filter struct {
	FirstName   string
	LastName    string
	PhoneNum    string
	Email       string
	Gender      string
	Role        string
	ServiceRole string
	Page        int
	Limit       int
}

type BookingFilter struct {
	RequesterID string
	ProviderID  string
	ServiceType string
	BookingType string
	Status      string
	StartDate   string // time.Time
	EndDate     string // time.Time
	StartTime   string // time.Time
	EndTime     string // time.Time
	Page        int
	Limit       int
}

type NotificationFilter struct {
	Type      string
	UserID    string
	BookingID string
	Seen      bool
	Page      int
	Limit     int
}

type ServiceFilter struct {
	Page  int
	Limit int
}

func (f Filter) Offset() int {
	return (f.Page - 1) * f.Limit
}

func (b BookingFilter) Offset() int {
	return (b.Page - 1) * b.Limit
}

func (n NotificationFilter) Offset() int {
	return (n.Page - 1) * n.Limit
}

func (s ServiceFilter) Offset() int {
	return (s.Page - 1) * s.Limit
}
