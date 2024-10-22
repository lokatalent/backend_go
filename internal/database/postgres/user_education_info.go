package postgres

import (
	"context"
	"database/sql"
	"errors"
	// "strings"

	"github.com/google/uuid"

	"github.com/lokatalent/backend_go/internal/models"
	"github.com/lokatalent/backend_go/internal/repository"
)

func (u *userImplementation) CreateEducationInfo(eduInfo *models.UserEducationInfo) error {
	stmt := `
    INSERT INTO users_education_info (
        user_id,
        institute,
        degree,
        discipline,
        start,
        finish
    ) VALUES (
        $1, $2, $3, $4, $5, $6
    ) RETURNING
        user_id,
        institute,
        degree,
        discipline,
        start,
        finish,
        created_at,
        updated_at;
    `

	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	err := u.DB.QueryRowContext(
		ctx,
		stmt,
		eduInfo.UserID,
		eduInfo.Institute,
		eduInfo.Degree,
		eduInfo.Discipline,
		eduInfo.Start,
		eduInfo.Finish,
	).Scan(
		&eduInfo.UserID,
		&eduInfo.Institute,
		&eduInfo.Degree,
		&eduInfo.Discipline,
		&eduInfo.Start,
		&eduInfo.Finish,
		&eduInfo.CreatedAt,
		&eduInfo.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (u *userImplementation) GetEducationInfo(userID string) (models.UserEducationInfo, error) {
	stmt := `
    SELECT
        user_id,
        institute,
        degree,
        discipline,
        start,
        finish,
        created_at,
        updated_at
    FROM users_education_info
    WHERE user_id = $1;
    `

	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	newEduInfo := models.UserEducationInfo{}
	err := u.DB.QueryRowContext(ctx, stmt, userID).Scan(
		&newEduInfo.UserID,
		&newEduInfo.Institute,
		&newEduInfo.Degree,
		&newEduInfo.Discipline,
		&newEduInfo.Start,
		&newEduInfo.Finish,
		&newEduInfo.CreatedAt,
		&newEduInfo.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.UserEducationInfo{}, repository.ErrRecordNotFound
		}
		return models.UserEducationInfo{}, err
	}

	return newEduInfo, nil
}

func (u *userImplementation) UpdateEducationInfo(eduInfo *models.UserEducationInfo) error {
	stmt := `
    UPDATE users_education_info
    SET
        institute = $2,
        degree = $3,
        discipline = $4,
        start = $5,
        finish = $6,
        updated_at = now()
    WHERE user_id = $1
    RETURNING
        user_id,
        institute,
        degree,
        discipline,
        start,
        finish,
        created_at,
        updated_at;
    `

	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	err := u.DB.QueryRowContext(
		ctx,
		stmt,
		eduInfo.UserID,
		eduInfo.Institute,
		eduInfo.Degree,
		eduInfo.Discipline,
		eduInfo.Start,
		eduInfo.Finish,
	).Scan(
		&eduInfo.UserID,
		&eduInfo.Institute,
		&eduInfo.Degree,
		&eduInfo.Discipline,
		&eduInfo.Start,
		&eduInfo.Finish,
		&eduInfo.CreatedAt,
		&eduInfo.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

// Certifications
func (u *userImplementation) CreateCertification(cert *models.UserCertification) error {
	if cert.ID == "" {
		cert.ID = uuid.NewString()
	}
	stmt := `
    INSERT INTO users_certifications (
        id,
        user_id,
        url
    ) VALUES (
        $1, $2, $3
    ) RETURNING
        id,
        user_id,
        url;
    `
	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	err := u.DB.QueryRowContext(
		ctx,
		stmt,
		cert.ID,
		cert.UserID,
		cert.URL,
	).Scan(
		&cert.ID,
		&cert.UserID,
		&cert.URL,
	)
	if err != nil {
		return err
	}

	return nil
}

func (u *userImplementation) GetCertifications(userID string) ([]models.UserCertification, error) {
	stmt := `
    SELECT
        id,
        user_id,
        url
    FROM users_certifications
    WHERE user_id = $1
    `

	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	rows, err := u.DB.QueryContext(ctx, stmt, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrRecordNotFound
		}
		return nil, err
	}

	certs := []models.UserCertification{}
	for rows.Next() {
		newCert := models.UserCertification{}
		err := rows.Scan(&newCert.ID, &newCert.UserID, &newCert.URL)
		if err != nil {
			return nil, err
		}
		certs = append(certs, newCert)
	}

	return certs, nil
}

func (u *userImplementation) DeleteCertification(id, userID string) error {
	stmt := `
    DELETE FROM users_certifications
    WHERE id = $1 AND user_id = $2
    `

	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	_, err := u.DB.ExecContext(ctx, stmt, id, userID)
	if err != nil {
		return err
	}

	return nil
}
