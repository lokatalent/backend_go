package repository

import "github.com/lokatalent/backend_go/internal/models"

type CommissionRepository interface {
	CreateServiceCommission(serviceCom *models.ServiceCommission) error
	GetServiceCommission() (models.ServiceCommission, error)
	UpdateServiceCommission(serviceCom *models.ServiceCommission) error
}
