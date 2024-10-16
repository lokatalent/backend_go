package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"github.com/lokatalent/backend_go/internal/models"
	"github.com/lokatalent/backend_go/internal/repository"
)

func (u *userImplementation) CreateVerificationCode(verCode *models.UserVerificationCode) error {
	if verCode.ID == "" {
		verCode.ID = uuid.NewString()
	}
	stmt := `
    INSERT INTO contact_verifications (
        id,
        user_id,
        code,
        contact_type
    ) VALUES (
        $1, $2, $3, $4
    ) RETURNING
        id,
        user_id,
        code,
        contact_type,
        created_at,
        expires_at;
    `

	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	err := u.DB.QueryRowContext(
		ctx,
		stmt,
		verCode.ID,
		verCode.UserID,
		verCode.Code,
		verCode.ContactType,
	).Scan(
		&verCode.ID,
		&verCode.UserID,
		&verCode.Code,
		&verCode.ContactType,
		&verCode.CreatedAt,
		&verCode.ExpiresAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (u *userImplementation) GetVerificationCode(userID, contactType string) (models.UserVerificationCode, error) {
	stmt := `
    SELECT
        id,
        user_id,
        code,
        contact_type,
        created_at,
        expires_at
    FROM contact_verifications
    WHERE user_id=$1 AND contact_type=$2
    `

	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	newCode := models.UserVerificationCode{}
	err := u.DB.QueryRowContext(
		ctx,
		stmt,
		userID,
		contactType,
	).Scan(
		&newCode.ID,
		&newCode.UserID,
		&newCode.Code,
		&newCode.ContactType,
		&newCode.CreatedAt,
		&newCode.ExpiresAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.UserVerificationCode{}, repository.ErrRecordNotFound
		}
		return models.UserVerificationCode{}, err
	}

	return newCode, nil
}

func (u *userImplementation) DeleteVerificationCode(userID, contactType string) error {
	stmt := `
    DELETE FROM contact_verifications
        WHERE user_id=$1 AND contact_type=$2
    `

	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	_, err := u.DB.ExecContext(ctx, stmt, userID, contactType)
	if err != nil {
		return err
	}

	return nil
}
