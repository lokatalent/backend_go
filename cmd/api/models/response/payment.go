package response

import (
	"github.com/lokatalent/backend_go/internal/models"
)

type TrackPayment struct {
	PaymentHistory []models.TrackPayment `json:"payment_history"`
	Escrow         float64               `json:"in_escrow"`
}
