package postgres

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"strings"

	"github.com/google/uuid"

	"github.com/lokatalent/backend_go/internal/models"
	"github.com/lokatalent/backend_go/internal/repository"
)

type bookingImplementation struct {
	DB *sql.DB
}

func NewBookingImplementation(db *sql.DB) repository.BookingRepository {
	return &bookingImplementation{DB: db}
}

func (b *bookingImplementation) Create(booking *models.Booking) error {
	if booking.ID == "" {
		booking.ID = uuid.NewString()
	}
	stmt := `
    INSERT INTO bookings (
        id,
        requester_id,
        requester_addr,
        service_type,
        booking_type,
        service_desc,
        start_time,
        end_time,
        start_date,
        end_date,
        total_price,
        status
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
        $11, $12
    ) RETURNING 
        id,
        requester_id,
        requester_addr,
        service_type,
        booking_type,
        service_desc,
        start_time,
        end_time,
        start_date,
        end_date,
        total_price,
        actual_price,
        status,
        created_at,
        updated_at;
    `
	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	err := b.DB.QueryRowContext(
		ctx,
		stmt,
		booking.ID,
		booking.RequesterID,
		booking.RequesterAddr,
		booking.ServiceType,
		booking.BookingType,
		booking.ServiceDesc,
		booking.StartTime,
		booking.EndTime,
		booking.StartDate,
		booking.EndDate,
		booking.TotalPrice,
		booking.Status,
	).Scan(
		&booking.ID,
		&booking.RequesterID,
		&booking.RequesterAddr,
		&booking.ServiceType,
		&booking.BookingType,
		&booking.ServiceDesc,
		&booking.StartTime,
		&booking.EndTime,
		&booking.StartDate,
		&booking.EndDate,
		&booking.TotalPrice,
		&booking.ActualPrice,
		&booking.Status,
		&booking.CreatedAt,
		&booking.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (b *bookingImplementation) GetByID(id string) (models.Booking, error) {
	stmt := `
    SELECT
        id,
        requester_id,
        provider_id,
        requester_addr,
        service_type,
        booking_type,
        service_desc,
        start_time,
        end_time,
        start_date,
        end_date,
        total_price,
        actual_price,
        status,
        created_at,
        updated_at
    FROM bookings
    WHERE id = $1
    LIMIT 1
    `
	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	booking := models.Booking{}
	err := b.DB.QueryRowContext(ctx, stmt, id).Scan(
		&booking.ID,
		&booking.RequesterID,
		&booking.ProviderID,
		&booking.RequesterAddr,
		&booking.ServiceType,
		&booking.BookingType,
		&booking.ServiceDesc,
		&booking.StartTime,
		&booking.EndTime,
		&booking.StartDate,
		&booking.EndDate,
		&booking.TotalPrice,
		&booking.ActualPrice,
		&booking.Status,
		&booking.CreatedAt,
		&booking.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Booking{}, repository.ErrRecordNotFound
		}
		return models.Booking{}, err
	}

	return booking, nil
}

func (b *bookingImplementation) GetAll(filter models.BookingFilter) ([]models.Booking, error) {
	stmt := `
    WITH filtered_bookings AS (
        SELECT *
        FROM bookings
        WHERE
            ($1 = '' OR requester_id = $1::UUID)
            -- Provider filter (handle NULL provider_id cases)
            AND (
                NULLIF($2, '')::UUID IS NULL 
                OR (
                    CASE 
                        WHEN NULLIF($2, '') IS NULL THEN provider_id IS NULL
                        ELSE provider_id = $2::UUID
                    END
                )
            )
            AND ($3 = '' OR booking_type = $3)
            AND ($4 = '' OR service_type = $4)
            AND ($5 = '' OR status = $5)
            AND ($6 = '' OR start_date >= $6::DATE)
            AND ($7 = '' OR end_date <= $7::DATE)
            AND ($8 = '' OR start_time >= $8::TIMETZ)
            AND ($9 = '' OR end_time <= $9::TIMETZ)
    )
    SELECT 
        id,
        requester_id,
        provider_id,
        requester_addr,
        service_type,
        booking_type,
        service_desc,
        start_time,
        end_time,
        start_date,
        end_date,
        total_price,
        actual_price,
        status,
        created_at,
        updated_at
    FROM filtered_bookings
    ORDER BY start_date ASC, start_time ASC
    LIMIT $10 OFFSET $11;
    `
	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	rows, err := b.DB.QueryContext(
		ctx,
		stmt,
		filter.RequesterID,
		filter.ProviderID,
		filter.BookingType,
		filter.ServiceType,
		filter.Status,
		filter.StartDate,
		filter.EndDate,
		filter.StartTime,
		filter.EndTime,
		filter.Limit,
		filter.Offset(),
	)
	if err != nil {
		return nil, err
	}

	bookings := []models.Booking{}
	for rows.Next() {
		booking := models.Booking{}
		err := rows.Scan(
			&booking.ID,
			&booking.RequesterID,
			&booking.ProviderID,
			&booking.RequesterAddr,
			&booking.ServiceType,
			&booking.BookingType,
			&booking.ServiceDesc,
			&booking.StartTime,
			&booking.EndTime,
			&booking.StartDate,
			&booking.EndDate,
			&booking.TotalPrice,
			&booking.ActualPrice,
			&booking.Status,
			&booking.CreatedAt,
			&booking.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}

	return bookings, nil
}

func (b *bookingImplementation) UpdateStatus(id, status string) (models.Booking, error) {
	stmt := `
	UPDATE bookings
	SET
	    status = $2,
	    updated_at = now()
	WHERE id = $1
    RETURNING 
        id,
        requester_id,
        provider_id,
        requester_addr,
        service_type,
        booking_type,
        service_desc,
        start_time,
        end_time,
        start_date,
        end_date,
        total_price,
        actual_price,
        status,
        created_at,
        updated_at;
    `
	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	booking := models.Booking{}
	err := b.DB.QueryRowContext(ctx, stmt, id, status).Scan(
		&booking.ID,
		&booking.RequesterID,
		&booking.ProviderID,
		&booking.RequesterAddr,
		&booking.ServiceType,
		&booking.BookingType,
		&booking.ServiceDesc,
		&booking.StartTime,
		&booking.EndTime,
		&booking.StartDate,
		&booking.EndDate,
		&booking.TotalPrice,
		&booking.ActualPrice,
		&booking.Status,
		&booking.CreatedAt,
		&booking.UpdatedAt,
	)
	if err != nil {
		return models.Booking{}, err
	}
	return booking, nil
}

func (b *bookingImplementation) AssignProvider(id, userID string) (models.Booking, error) {
	stmt := `
	UPDATE bookings
	SET
	    provider_id = $2,
	    updated_at = now()
	WHERE id = $1
    RETURNING 
        id,
        requester_id,
        provider_id,
        requester_addr,
        service_type,
        booking_type,
        service_desc,
        start_time,
        end_time,
        start_date,
        end_date,
        total_price,
        actual_price,
        status,
        created_at,
        updated_at;
    `
	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	booking := models.Booking{}
	err := b.DB.QueryRowContext(ctx, stmt, id, userID).Scan(
		&booking.ID,
		&booking.RequesterID,
		&booking.ProviderID,
		&booking.RequesterAddr,
		&booking.ServiceType,
		&booking.BookingType,
		&booking.ServiceDesc,
		&booking.StartTime,
		&booking.EndTime,
		&booking.StartDate,
		&booking.EndDate,
		&booking.TotalPrice,
		&booking.ActualPrice,
		&booking.Status,
		&booking.CreatedAt,
		&booking.UpdatedAt,
	)
	if err != nil {
		return models.Booking{}, err
	}
	return booking, nil
}

func (b *bookingImplementation) RejectBooking(id, userID string) error {
	stmt := `
    INSERT INTO rejected_bookings (
        id,
        user_id,
        booking_id
    ) VALUES (
        $1, $2, $3
    );
    `
	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	_, err := b.DB.ExecContext(ctx, stmt, uuid.NewString(), userID, id)
	if err != nil {
		return err
	}

	return nil
}

// checks if a booking has been rejected by a user
func (b *bookingImplementation) CheckRejected(id, userID string) (bool, error) {
	stmt := `
    SELECT
        user_id,
        booking_id
    FROM rejected_bookings
    WHERE user_id = $1 AND booking_id = $2
    LIMIT 1
    `
	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	var user_id, booking_id string
	err := b.DB.QueryRowContext(ctx, stmt, userID, id).Scan(&user_id, &booking_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// match available service providers with a booking
func (b *bookingImplementation) MatchServices(
	booking *models.Booking,
	filter models.ServiceFilter,
) ([]models.UserService, error) {
	stmt := `
    SELECT DISTINCT s.*
    FROM services s
    WHERE 
        -- Match address by city, state, country
        -- s.address LIKE '%' || $2 || ', ' || $3 || ', ' || $4
        s.address ~* ($2 || ', ' || $3 || ', ' || $4 || '$')
        
        -- Match service type
        AND s.service_type = $1
        
        -- Exclude services with users who have in-progress bookings
        AND s.user_id NOT IN (
            SELECT provider_id 
            FROM bookings 
            WHERE 
                provider_id IS NOT NULL 
                AND status = $5
        )
    LIMIT $6 OFFSET $7;
    `
	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	splittedAddr := strings.Split(booking.RequesterAddr, ",")
	slices.Reverse(splittedAddr)

	rows, err := b.DB.QueryContext(
		ctx,
		stmt,
		booking.ServiceType,
		strings.TrimSpace(splittedAddr[2]), // city
		strings.TrimSpace(splittedAddr[1]), // state
		strings.TrimSpace(splittedAddr[0]), // country
		models.BOOKING_IN_PROGRESS,
		filter.Limit,
		filter.Offset(),
	)
	if err != nil {
		return nil, err
	}

	services := []models.UserService{}
	for rows.Next() {
		service := models.UserService{}
		err := rows.Scan(
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
			return nil, err
		}
		services = append(services, service)
	}
	return services, nil
}
