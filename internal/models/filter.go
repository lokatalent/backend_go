package models

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

func (f Filter) Offset() int {
	return (f.Page - 1) * f.Limit
}
