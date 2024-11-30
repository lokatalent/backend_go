package repository

import "github.com/lokatalent/backend_go/internal/models"

type PaymentRepository interface {
	// users wallets
	CreateWallet(wallet *models.UserWallet) error
	GetWallet(userID string) (models.UserWallet, error)
	GetUserDebits(userID string) (float64, error)
	UpdateWallet(userID, paymentType string, amount float64) error

	// transactions
	CreatePayment(payment *models.Payment) error
	GetPayment(filter models.PaymentFilter) (models.Payment, error)
	UpdatePaymentStatus(id, status string) (models.Payment, error)

	// paystack transaction utilities
	CreateAccessCode(id, paymentID, accessCode string) error
	CreateRecipientCode(id, userID, recipientCode string) error
	GetAccessCode(paymentID string) (string, error)
	GetRecipientCode(userID string) (string, error)
	UpdateRecipientCode(userID, recipientCode string) error
	DeleteAccessCode(paymentID string) error
}
