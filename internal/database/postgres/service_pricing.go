package postgres

import (
	"context"
	"database/sql"
	// "errors"
	// "strings"

	"github.com/google/uuid"

	"github.com/lokatalent/backend_go/internal/models"
	"github.com/lokatalent/backend_go/internal/repository"
)

type servicePricingImplementation struct {
	DB *sql.DB
}

func NewServicePricingImplementation(db *sql.DB) repository.ServicePricingRepository {
	return &servicePricingImplementation{DB: db}
}

func (s *servicePricingImplementation) CreateServicePricing(servicePricing *models.ServicePricing) error {
	if servicePricing.ID == "" {
		servicePricing.ID = uuid.NewString()
	}
	stmt := `
    INSERT INTO services_pricing (
        id,
        service_type,
        rate_per_hour
    ) VALUES (
        $1, $2, $3
    ) RETURNING
        id,
        service_type,
        rate_per_hour,
        created_at,
        updated_at;
    `
	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	err := s.DB.QueryRowContext(
		ctx,
		stmt,
		servicePricing.ID,
		servicePricing.ServiceType,
		servicePricing.RatePerHour,
	).Scan(
		&servicePricing.ID,
		&servicePricing.ServiceType,
		&servicePricing.RatePerHour,
		&servicePricing.CreatedAt,
		&servicePricing.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *servicePricingImplementation) GetServicePricing(serviceType string) (models.ServicePricing, error) {
	stmt := `
    SELECT
        id,
        service_type,
        rate_per_hour,
        created_at,
        updated_at
    FROM services_pricing
    WHERE service_type = $1
    LIMIT 1
    `
	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	servicePricing := models.ServicePricing{}
	err := s.DB.QueryRowContext(ctx, stmt, serviceType).Scan(
		&servicePricing.ID,
		&servicePricing.ServiceType,
		&servicePricing.RatePerHour,
		&servicePricing.CreatedAt,
		&servicePricing.UpdatedAt,
	)
	if err != nil {
		return models.ServicePricing{}, err
	}

	return servicePricing, nil
}

func (s *servicePricingImplementation) GetAllServicesPricing() ([]models.ServicePricing, error) {
	stmt := `
    SELECT
        id,
        service_type,
        rate_per_hour,
        created_at,
        updated_at
    FROM services_pricing
    `
	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	rows, err := s.DB.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}

	servicesPricing := []models.ServicePricing{}
	for rows.Next() {
		servicePricing := models.ServicePricing{}
		err := rows.Scan(
			&servicePricing.ID,
			&servicePricing.ServiceType,
			&servicePricing.RatePerHour,
			&servicePricing.CreatedAt,
			&servicePricing.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		servicesPricing = append(servicesPricing, servicePricing)
	}

	return servicesPricing, nil
}

func (s *servicePricingImplementation) UpdateServicePricing(servicePricing *models.ServicePricing) error {
	stmt := `
    UPDATE services_pricing
    SET
        rate_per_hour = $2,
        updated_at = now()
    WHERE service_type = $1
    RETURNING
        id,
        service_type,
        rate_per_hour,
        created_at,
        updated_at;
    `

	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	err := s.DB.QueryRowContext(
		ctx,
		stmt,
		servicePricing.ServiceType,
		servicePricing.RatePerHour,
	).Scan(
		&servicePricing.ID,
		&servicePricing.ServiceType,
		&servicePricing.RatePerHour,
		&servicePricing.CreatedAt,
		&servicePricing.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *servicePricingImplementation) DeleteServicePricing(serviceType string) error {
	stmt := `
    DELETE FROM services_pricing
    WHERE service_type = $1
    `

	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	_, err := s.DB.ExecContext(ctx, stmt, serviceType)
	if err != nil {
		return err
	}

	return nil
}
