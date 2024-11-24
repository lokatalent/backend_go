package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/lokatalent/backend_go/internal/models"
	"github.com/lokatalent/backend_go/internal/repository"
)

func (u *userImplementation) CreateBankInfo(bankInfo *models.UserBankInfo) error {
	stmt := `
    INSERT INTO users_bank_info (
        user_id,
        bank_name,
        account_name,
        account_num,
        bank_code
    ) VALUES (
        $1, $2, $3, $4, $5
    ) RETURNING
        user_id,
        bank_name,
        account_name,
        account_num,
        bank_code,
        created_at,
        updated_at;
    `

	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	err := u.DB.QueryRowContext(
		ctx,
		stmt,
		bankInfo.UserID,
		bankInfo.BankName,
		bankInfo.AccountName,
		bankInfo.AccountNum,
		bankInfo.BankCode,
	).Scan(
		&bankInfo.UserID,
		&bankInfo.BankName,
		&bankInfo.AccountName,
		&bankInfo.AccountNum,
		&bankInfo.BankCode,
		&bankInfo.CreatedAt,
		&bankInfo.UpdatedAt,
	)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), duplicateBankAcctNum):
			return repository.ErrDuplicateBankDetails
		default:
			return err
		}
	}

	return nil
}

func (u *userImplementation) GetBankInfo(userID string) (models.UserBankInfo, error) {
	stmt := `
    SELECT
        user_id,
        bank_name,
        account_name,
        account_num,
        bank_code,
        created_at,
        updated_at
    FROM users_bank_info
    WHERE user_id = $1;
    `

	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	newBankInfo := models.UserBankInfo{}
	err := u.DB.QueryRowContext(ctx, stmt, userID).Scan(
		&newBankInfo.UserID,
		&newBankInfo.BankName,
		&newBankInfo.AccountName,
		&newBankInfo.AccountNum,
		&newBankInfo.BankCode,
		&newBankInfo.CreatedAt,
		&newBankInfo.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.UserBankInfo{}, repository.ErrRecordNotFound
		}
		return models.UserBankInfo{}, err
	}

	return newBankInfo, nil
}

func (u *userImplementation) UpdateBankInfo(bankInfo *models.UserBankInfo) error {
	stmt := `
    UPDATE users_bank_info
    SET
        bank_name = $2,
        account_name = $3,
        account_num = $4,
        bank_code = $5,
        updated_at = now()
    WHERE user_id = $1
    RETURNING
        user_id,
        bank_name,
        account_name,
        account_num,
        bank_code,
        created_at,
        updated_at;
    `

	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	err := u.DB.QueryRowContext(
		ctx,
		stmt,
		bankInfo.UserID,
		bankInfo.BankName,
		bankInfo.AccountName,
		bankInfo.AccountNum,
		bankInfo.BankCode,
	).Scan(
		&bankInfo.UserID,
		&bankInfo.BankName,
		&bankInfo.AccountName,
		&bankInfo.AccountNum,
		&bankInfo.BankCode,
		&bankInfo.CreatedAt,
		&bankInfo.UpdatedAt,
	)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), duplicateBankAcctNum):
			return repository.ErrDuplicateBankDetails
		default:
			return err
		}
	}

	return nil
}
