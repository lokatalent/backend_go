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

type commissionImplementation struct {
	DB *sql.DB
}

func NewCommissionImplementation(db *sql.DB) repository.CommissionRepository {
	return &commissionImplementation{DB: db}
}

func (c *commissionImplementation) CreateServiceCommission(serviceCom *models.ServiceCommission) error {
	if serviceCom.ID == "" {
		serviceCom.ID = uuid.NewString()
	}
	stmt := `
    INSERT INTO service_commission (
        id,
        percentage
    ) VALUES (
        $1, $2
    ) RETURNING
        id,
        percentage,
        created_at,
        updated_at;
    `

	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	err := c.DB.QueryRowContext(
		ctx,
		stmt,
		serviceCom.ID,
		serviceCom.Percentage,
	).Scan(
		&serviceCom.ID,
		&serviceCom.Percentage,
		&serviceCom.CreatedAt,
		&serviceCom.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *commissionImplementation) GetServiceCommission() (models.ServiceCommission, error) {
	stmt := `
    SELECT
        id,
        percentage,
        created_at,
        updated_at
    FROM service_commission
    LIMIT 1
    `

	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	serviceCommission := models.ServiceCommission{}
	err := c.DB.QueryRowContext(ctx, stmt).Scan(
		&serviceCommission.ID,
		&serviceCommission.Percentage,
		&serviceCommission.CreatedAt,
		&serviceCommission.UpdatedAt,
	)
	if err != nil {
		return models.ServiceCommission{}, err
	}

	return serviceCommission, nil
}

func (c *commissionImplementation) UpdateServiceCommission(serviceCom *models.ServiceCommission) error {
	stmt := `
    UPDATE service_commission
    SET
        percentage = $2,
        updated_at = now()
    WHERE id = $1
    RETURNING
        id,
        percentage,
        created_at,
        updated_at;
    `

	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	err := c.DB.QueryRowContext(
		ctx,
		stmt,
		serviceCom.ID,
		serviceCom.Percentage,
	).Scan(
		&serviceCom.ID,
		&serviceCom.Percentage,
		&serviceCom.CreatedAt,
		&serviceCom.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}
