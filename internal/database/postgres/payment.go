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

type paymentImplementation struct {
	DB *sql.DB
}

func NewPaymentImplementation(db *sql.DB) repository.PaymentRepository {
	return &paymentImplementation{DB: db}
}

// users wallets

func (p *paymentImplementation) CreateWallet(wallet *models.UserWallet) error {
	if wallet.ID == "" {
		wallet.ID = uuid.NewString()
	}
	stmt := `
    INSERT INTO wallets (
        id,
        user_id,
        credits,
        debits
    ) VALUES (
        $1, $2, 0, 0
    ) RETURNING
        id,
        user_id,
        credits,
        debits,
        created_at,
        updated_at;
    `
	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	var walletCredits, walletDebits float64
	err := p.DB.QueryRowContext(ctx, stmt, wallet.ID, wallet.UserID).Scan(
		&wallet.ID,
		&wallet.UserID,
		&walletCredits,
		&walletDebits,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)
	if err != nil {
		return err
	}

	wallet.Balance = walletCredits - walletDebits
	wallet.Debits = walletDebits

	return nil
}

func (p *paymentImplementation) GetWallet(userID string) (models.UserWallet, error) {
	stmt := `
    SELECT
        id,
        user_id,
        credits,
        debits,
        created_at,
        updated_at
    FROM wallets
    WHERE user_id = $1
    LIMIT 1
    `
	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	var walletCredits, walletDebits float64
	wallet := models.UserWallet{}
	err := p.DB.QueryRowContext(ctx, stmt, userID).Scan(
		&wallet.ID,
		&wallet.UserID,
		&walletCredits,
		&walletDebits,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.UserWallet{}, repository.ErrRecordNotFound
		}
		return models.UserWallet{}, err
	}

	wallet.Balance = walletCredits - walletDebits
	wallet.Debits = walletDebits
	if wallet.Balance < 0 {
		return models.UserWallet{}, repository.ErrInvalidWalletBalance
	}

	return wallet, nil
}

func (p *paymentImplementation) GetUserDebits(userID string) (float64, error) {
	stmt := `
    SELECT debits
    FROM wallets
    WHERE user_id = $1
    LIMIT = 1;
    `
	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	var walletDebits float64
	err := p.DB.QueryRowContext(ctx, stmt, userID).Scan(&walletDebits)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0.0, repository.ErrRecordNotFound
		}
		return 0.0, err
	}

	return walletDebits, nil
}

func (p *paymentImplementation) UpdateWallet(userID, paymentType string, amount float64) error {
	stmt := ``
	switch paymentType {
	case models.PAYMENT_TYPE_CREDIT:
		stmt = `
        UPDATE wallets
        SET
            credits = credits + $2,
            debits = debits + $2
        WHERE user_id = $1;
        `
	case models.PAYMENT_TYPE_DEBIT:
		stmt = `
        UPDATE wallets
        SET
            -- credits = credits - $2,
            debits = debits + $2
        WHERE user_id = $1;
        `
	case models.PAYMENT_TYPE_REFUND:
		stmt = `
        UPDATE wallets
        SET
            -- credits = credits + $2,
            debits = debits - $2
        WHERE user_id = $1;
        `
	default:
		return errors.New("unknown payment type")
	}

	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	_, err := p.DB.ExecContext(ctx, stmt, userID, amount)
	if err != nil {
		return err
	}

	return nil
}

// transactions

func (p *paymentImplementation) CreatePayment(payment *models.Payment) error {
	if payment.ID == "" {
		payment.ID = uuid.NewString()
	}
	stmt := `
    INSERT INTO payments (
        id,
        type,
        booking_id,
        amount,
        payment_ref,
        status
    ) VALUES (
        $1, $2, $3, $4, $5, $6
    ) RETURNING
        id,
        type,
        booking_id,
        amount,
        payment_ref,
        status,
        created_at,
        updated_at;
    `
	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	err := p.DB.QueryRowContext(
		ctx,
		stmt,
		payment.ID,
		payment.Type,
		payment.BookingID,
		payment.Amount,
		payment.PaymentRef,
		payment.Status,
	).Scan(
		&payment.ID,
		&payment.Type,
		&payment.BookingID,
		&payment.Amount,
		&payment.PaymentRef,
		&payment.Status,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *paymentImplementation) GetPayment(filter models.PaymentFilter) (models.Payment, error) {
	stmt := `
    SELECT
        id,
        type,
        booking_id,
        amount,
        payment_ref,
        status,
        created_at,
        updated_at
    FROM payments
    WHERE
        ($1 = '' OR id = $1::UUID) AND
        ($2 = '' OR booking_id = $2::UUID) AND
        ($3 = '' OR payment_ref = $3::UUID) AND
        ($4 = '' OR type = $4) AND
        ($5 = '' OR status = $5)
    LIMIT 1;
    `
	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	payment := models.Payment{}
	err := p.DB.QueryRowContext(
		ctx,
		stmt,
		filter.ID,
		filter.BookingID,
		filter.PaymentRef,
		filter.Type,
		filter.Status,
	).Scan(
		&payment.ID,
		&payment.Type,
		&payment.BookingID,
		&payment.Amount,
		&payment.PaymentRef,
		&payment.Status,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Payment{}, repository.ErrRecordNotFound
		}
		return models.Payment{}, err
	}

	return payment, nil
}

func (p *paymentImplementation) UpdatePaymentStatus(id, status string) (models.Payment, error) {
	stmt := `
    UPDATE payments
    SET
        status = $2,
        updated_at = now()
    WHERE id = $1
    RETURNING
        id,
        type,
        booking_id,
        amount,
        payment_ref,
        status,
        created_at,
        updated_at;
    `
	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	payment := models.Payment{}
	err := p.DB.QueryRowContext(ctx, stmt, id, status).Scan(
		&payment.ID,
		&payment.Type,
		&payment.BookingID,
		&payment.Amount,
		&payment.PaymentRef,
		&payment.Status,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		return models.Payment{}, err
	}
	return payment, nil
}

// paystack transaction utilities

func (p *paymentImplementation) CreateAccessCode(id, paymentID, accessCode string) error {
	stmt := `
    INSERT INTO payment_access_codes (
        id,
        payment_id,
        access_code
    ) VALUES (
        $1, $2, $3
    )
    `
	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	_, err := p.DB.ExecContext(ctx, stmt, id, paymentID, accessCode)
	if err != nil {
		return err
	}
	return nil
}

func (p *paymentImplementation) CreateRecipientCode(id, userID, recipientCode string) error {
	stmt := `
    INSERT INTO payment_recipient_codes (
        id,
        user_id,
        recipient_code
    ) VALUES (
        $1, $2, $3
    )
    `
	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	_, err := p.DB.ExecContext(ctx, stmt, id, userID, recipientCode)
	if err != nil {
		return err
	}
	return nil
}

func (p *paymentImplementation) GetAccessCode(userID string) (string, error) {
	stmt := `
    SELECT access_code
    FROM payment_access_codes
    WHERE payment_id = $1
    LIMIT 1;
    `
	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	var access_code string
	err := p.DB.QueryRowContext(ctx, stmt, userID).Scan(&access_code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", repository.ErrRecordNotFound
		}
		return "", err
	}
	return access_code, nil
}

func (p *paymentImplementation) GetRecipientCode(userID string) (string, error) {
	stmt := `
    SELECT recipient_code
    FROM payment_recipient_codes
    WHERE user_id = $1
    LIMIT 1;
    `
	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	var recipient_code string
	err := p.DB.QueryRowContext(ctx, stmt, userID).Scan(&recipient_code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", repository.ErrRecordNotFound
		}
		return "", err
	}
	return recipient_code, nil
}

func (p *paymentImplementation) UpdateRecipientCode(userID, recipientCode string) error {
	stmt := `
    UPDATE payment_recipient_codes
    SET
        recipient_code = $2,
        updated_at = now()
    WHERE user_id = $1;
    `
	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	_, err := p.DB.ExecContext(ctx, stmt, userID, recipientCode)
	if err != nil {
		return err
	}
	return nil
}

func (p *paymentImplementation) DeleteAccessCode(paymentID string) error {
	stmt := `
    DELETE FROM payment_access_codes
    WHERE payment_id = $1;
    `
	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	_, err := p.DB.ExecContext(ctx, stmt, paymentID)
	if err != nil {
		return err
	}

	return nil
}
