package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	// "time"

	"github.com/google/uuid"

	"github.com/lokatalent/backend_go/internal/models"
	"github.com/lokatalent/backend_go/internal/repository"
)

func (u *userImplementation) CreateService(service *models.UserService) error {
	if service.ID == "" {
		service.ID = uuid.NewString()
	}
	stmt := `
    INSERT INTO services (
        id,
        user_id,
        service_type,
        service_desc,
        rate_per_hour,
        experience_years,
        availability,
        address
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8
    ) RETURNING
        id,
        user_id,
        service_type,
        service_desc,
        rate_per_hour,
        experience_years,
        availability,
        address,
        created_at,
        updated_at;
    `

	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	err := u.DB.QueryRowContext(
		ctx,
		stmt,
		service.ID,
		service.UserID,
		service.ServiceType,
		service.ServiceDesc,
		service.RatePerHour,
		service.ExperienceYears,
		service.Availability,
		service.Address,
	).Scan(
		&service.ID,
		&service.UserID,
		&service.ServiceType,
		&service.ServiceDesc,
		&service.RatePerHour,
		&service.ExperienceYears,
		&service.Availability,
		&service.Address,
		&service.CreatedAt,
		&service.UpdatedAt,
	)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), duplicateService):
			return repository.ErrDuplicateService
		default:
			return err
		}
	}

	return nil
}

func (u *userImplementation) GetService(userID, serviceType string) (models.UserService, error) {
	stmt := `
    SELECT
        id,
        user_id,
        service_type,
        service_desc,
        rate_per_hour,
        experience_years,
        availability,
        address,
        created_at,
        updated_at
    FROM services
    WHERE user_id = $1 AND service_type=$2;
    `

	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	newService := models.UserService{}
	err := u.DB.QueryRowContext(ctx, stmt, userID, serviceType).Scan(
		&newService.ID,
		&newService.UserID,
		&newService.ServiceType,
		&newService.ServiceDesc,
		&newService.RatePerHour,
		&newService.ExperienceYears,
		&newService.Availability,
		&newService.Address,
		&newService.CreatedAt,
		&newService.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.UserService{}, repository.ErrRecordNotFound
		}
		return models.UserService{}, err
	}

	return newService, nil
}

func (u *userImplementation) GetAllServices(userID string) ([]models.UserService, error) {
	stmt := `
    SELECT
        id,
        user_id,
        service_type,
        service_desc,
        rate_per_hour,
        experience_years,
        availability,
        address,
        created_at,
        updated_at
    FROM services
    WHERE user_id = $1;
    `

	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	rows, err := u.DB.QueryContext(ctx, stmt, userID)
	if err != nil {
		return nil, err
	}
	services := []models.UserService{}
	for rows.Next() {
		newService := models.UserService{}
		err := rows.Scan(
			&newService.ID,
			&newService.UserID,
			&newService.ServiceType,
			&newService.ServiceDesc,
			&newService.RatePerHour,
			&newService.ExperienceYears,
			&newService.Availability,
			&newService.Address,
			&newService.CreatedAt,
			&newService.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		services = append(services, newService)
	}
	return services, nil
}

func (u *userImplementation) UpdateService(service *models.UserService) error {
	stmt := `
    UPDATE services
    SET
        service_desc = $3,
        rate_per_hour = $4,
        experience_years = $5,
        availability = $6,
        address = $7,
        updated_at = now()
    WHERE user_id = $1 AND service_type = $2
    RETURNING
        id,
        user_id,
        service_type,
        service_desc,
        rate_per_hour,
        experience_years,
        availability,
        address,
        created_at,
        updated_at
    `

	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	err := u.DB.QueryRowContext(
		ctx,
		stmt,
		service.UserID,
		service.ServiceType,
		service.ServiceDesc,
		service.RatePerHour,
		service.ExperienceYears,
		service.Availability,
		service.Address,
	).Scan(
		&service.ID,
		&service.UserID,
		&service.ServiceType,
		&service.ServiceDesc,
		&service.RatePerHour,
		&service.ExperienceYears,
		&service.Availability,
		&service.Address,
		&service.CreatedAt,
		&service.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (u *userImplementation) DeleteService(userID, serviceType string) error {
	stmt := `
    DELETE FROM services
    WHERE user_id = $1 AND service_type = $2;
    `

	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	_, err := u.DB.ExecContext(ctx, stmt, userID, serviceType)
	if err != nil {
		return err
	}

	return nil
}

// Service images
func (u *userImplementation) CreateServiceImage(serviceImg *models.ServiceImage) error {
	if serviceImg.ID == "" {
		serviceImg.ID = uuid.NewString()
	}
	stmt := `
    INSERT INTO service_images (
        id,
        user_id,
        service_type,
        url
    ) VALUES (
        $1, $2, $3, $4
    ) RETURNING
        id,
        user_id,
        service_type,
        url;
    `
	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	err := u.DB.QueryRowContext(
		ctx,
		stmt,
		serviceImg.ID,
		serviceImg.UserID,
		serviceImg.ServiceType,
		serviceImg.URL,
	).Scan(
		&serviceImg.ID,
		&serviceImg.UserID,
		&serviceImg.ServiceType,
		&serviceImg.URL,
	)
	if err != nil {
		return err
	}

	return nil
}

func (u *userImplementation) GetServiceImages(userID, serviceType string) ([]models.ServiceImage, error) {
	stmt := `
    SELECT
        id,
        user_id,
        service_type,
        url
    FROM service_images
    WHERE user_id = $1 AND service_type = $2
    `

	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	rows, err := u.DB.QueryContext(ctx, stmt, userID, serviceType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrRecordNotFound
		}
		return nil, err
	}

	imgs := []models.ServiceImage{}
	for rows.Next() {
		newImg := models.ServiceImage{}
		err := rows.Scan(&newImg.ID, &newImg.UserID, &newImg.ServiceType, &newImg.URL)
		if err != nil {
			return nil, err
		}
		imgs = append(imgs, newImg)
	}

	return imgs, nil
}

func (u *userImplementation) DeleteServiceImage(id, userID, serviceType string) error {
	stmt := `
    DELETE FROM service_images
    WHERE id = $1 AND user_id = $2 AND service_type = $3
    `

	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	_, err := u.DB.ExecContext(ctx, stmt, id, userID, serviceType)
	if err != nil {
		return err
	}

	return nil
}
