package repository

import "github.com/lokatalent/backend_go/internal/models"

type ServicePricingRepository interface {
	CreateServicePricing(servicePrice *models.ServicePricing) error
	GetServicePricing(serviceType string) (models.ServicePricing, error)
	GetAllServicesPricing() ([]models.ServicePricing, error)
	UpdateServicePricing(servicePrice *models.ServicePricing) error
	DeleteServicePricing(serviceType string) error
}
