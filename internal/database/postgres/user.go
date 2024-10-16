package postgres

import (
	"context"
	"database/sql"
	// "encoding/json"
	"errors"
	"strings"
	// "time"

	"github.com/google/uuid"

	"github.com/lokatalent/backend_go/internal/models"
	"github.com/lokatalent/backend_go/internal/repository"
)

type userImplementation struct {
	DB *sql.DB
}

func NewUserImplementation(db *sql.DB) repository.UserRepository {
	return &userImplementation{DB: db}
}

// Create adds new user to the database.
func (u *userImplementation) Create(user *models.User) error {
	// generate an ID for new user
	if user.ID == "" {
		user.ID = uuid.NewString()
	}

	stmt := `
    INSERT INTO users (
        id,
        first_name,
        last_name,
        email,
        phone_num,
        avatar,
        password
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7
    ) RETURNING
        id,
        first_name,
        last_name,
        email,
        phone_num,
        password,
        avatar,
        gender,
        date_of_birth,
        bio,
        address,
        role,
        service_role,
        is_verified,
        email_verified,
        phone_verified,
        created_at,
        updated_at;
    `

	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	err := u.DB.QueryRowContext(
		ctx,
		stmt,
		user.ID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.PhoneNum,
		user.Avatar,
		user.Password,
	).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.PhoneNum,
		&user.Password,
		&user.Avatar,
		&user.Gender,
		&user.DateOfBirth,
		&user.Bio,
		&user.Address,
		&user.Role,
		&user.ServiceRole,
		&user.IsVerified,
		&user.EmailVerified,
		&user.PhoneVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), duplicateEmail),
			strings.Contains(err.Error(), duplicatePhone):
			return repository.ErrDuplicateDetails
		default:
			return err
		}
	}

	return nil
}

// GetByID retrieves user from the database using their assigned ID.
func (u *userImplementation) GetByID(id string) (models.User, error) {
	stmt := `
    SELECT 
        id,
        first_name,
        last_name,
        email,
        phone_num,
        password,
        avatar,
        gender,
        date_of_birth,
        bio,
        address,
        role,
        service_role,
        is_verified,
        email_verified,
        phone_verified,
        created_at,
        updated_at
    FROM users 
    WHERE id = $1;
    `

	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	newUser := models.User{}
	err := u.DB.QueryRowContext(ctx, stmt, id).Scan(
		&newUser.ID,
		&newUser.FirstName,
		&newUser.LastName,
		&newUser.Email,
		&newUser.PhoneNum,
		&newUser.Password,
		&newUser.Avatar,
		&newUser.Gender,
		&newUser.DateOfBirth,
		&newUser.Bio,
		&newUser.Address,
		&newUser.Role,
		&newUser.ServiceRole,
		&newUser.IsVerified,
		&newUser.EmailVerified,
		&newUser.PhoneVerified,
		&newUser.CreatedAt,
		&newUser.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, repository.ErrRecordNotFound
		}
		return models.User{}, err
	}

	return newUser, nil
}

// GetByEmail retrieves user from the database using their email address.
func (u *userImplementation) GetByEmail(email string) (models.User, error) {
	stmt := `
    SELECT 
        id,
        first_name,
        last_name,
        email,
        phone_num,
        password,
        avatar,
        gender,
        date_of_birth,
        bio,
        address,
        role,
        service_role,
        is_verified,
        email_verified,
        phone_verified,
        created_at,
        updated_at
    FROM users 
    WHERE email = $1;
    `

	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	newUser := models.User{}
	err := u.DB.QueryRowContext(ctx, stmt, email).Scan(
		&newUser.ID,
		&newUser.FirstName,
		&newUser.LastName,
		&newUser.Email,
		&newUser.PhoneNum,
		&newUser.Password,
		&newUser.Avatar,
		&newUser.Gender,
		&newUser.DateOfBirth,
		&newUser.Bio,
		&newUser.Address,
		&newUser.Role,
		&newUser.ServiceRole,
		&newUser.IsVerified,
		&newUser.EmailVerified,
		&newUser.PhoneVerified,
		&newUser.CreatedAt,
		&newUser.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, repository.ErrRecordNotFound
		}
		return models.User{}, err
	}

	return newUser, nil
}

// GetByPhone retrieves user from the database using their phone number.
func (u *userImplementation) GetByPhone(phone string) (models.User, error) {
	stmt := `
    SELECT 
        id,
        first_name,
        last_name,
        email,
        phone_num,
        password,
        avatar,
        gender,
        date_of_birth,
        bio,
        address,
        role,
        service_role,
        is_verified,
        email_verified,
        phone_verified,
        created_at,
        updated_at
    FROM users 
    WHERE phone_num = $1;
    `

	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	newUser := models.User{}
	err := u.DB.QueryRowContext(ctx, stmt, phone).Scan(
		&newUser.ID,
		&newUser.FirstName,
		&newUser.LastName,
		&newUser.Email,
		&newUser.PhoneNum,
		&newUser.Password,
		&newUser.Avatar,
		&newUser.Gender,
		&newUser.DateOfBirth,
		&newUser.Bio,
		&newUser.Address,
		&newUser.Role,
		&newUser.ServiceRole,
		&newUser.IsVerified,
		&newUser.EmailVerified,
		&newUser.PhoneVerified,
		&newUser.CreatedAt,
		&newUser.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, repository.ErrRecordNotFound
		}
		return models.User{}, err
	}

	return newUser, nil
}

// GetAllUser retrieves all users in the database.
func (u *userImplementation) GetAllUsers(filter models.Filter) ([]models.User, error) {
	stmt := `
    SELECT 
        id,
        first_name,
        last_name,
        email,
        phone_num,
        avatar,
        gender,
        date_of_birth,
        bio,
        address,
        role,
        service_role,
        is_verified,
        email_verified,
        phone_verified,
        created_at,
        updated_at
    FROM users
    ORDER by first_name
    LIMIT $1 OFFSET $2
    `

	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	rows, err := u.DB.QueryContext(
		ctx,
		stmt,
		filter.Limit,
		filter.Offset())
	if err != nil {
		return nil, err
	}
	users := []models.User{}
	for rows.Next() {
		newUser := models.User{}
		err := rows.Scan(
			&newUser.ID,
			&newUser.FirstName,
			&newUser.LastName,
			&newUser.Email,
			&newUser.PhoneNum,
			&newUser.Avatar,
			&newUser.Gender,
			&newUser.DateOfBirth,
			&newUser.Bio,
			&newUser.Address,
			&newUser.Role,
			&newUser.ServiceRole,
			&newUser.IsVerified,
			&newUser.EmailVerified,
			&newUser.PhoneVerified,
			&newUser.CreatedAt,
			&newUser.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, newUser)
	}
	return users, nil
}

// Update changes the user's data
func (u *userImplementation) Update(user *models.User) error {
	stmt := `
    UPDATE users
    SET
        first_name = $2,
        last_name = $3,
        phone_num = $4,
        bio = $5,
        address = $6,
        is_verified = $7,
        gender = $8,
        date_of_birth = $9,
        updated_at = now()
    WHERE id = $1
    RETURNING
        id,
        first_name,
        last_name,
        email,
        phone_num,
        avatar,
        gender,
        date_of_birth,
        bio,
        address,
        role,
        service_role,
        is_verified,
        email_verified,
        phone_verified,
        created_at,
        updated_at;
    `

	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	err := u.DB.QueryRowContext(
		ctx,
		stmt,
		user.ID,
		user.FirstName,
		user.LastName,
		user.PhoneNum,
		user.Bio,
		user.Address,
		user.IsVerified,
		user.Gender,
		user.DateOfBirth,
	).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.PhoneNum,
		&user.Avatar,
		&user.Gender,
		&user.DateOfBirth,
		&user.Bio,
		&user.Address,
		&user.Role,
		&user.ServiceRole,
		&user.IsVerified,
		&user.EmailVerified,
		&user.PhoneVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), duplicateEmail),
			strings.Contains(err.Error(), duplicatePhone):
			return repository.ErrDuplicateDetails
		default:
			return err
		}
	}
	return nil
}

// UpdateImage updates the user's profile image
func (u *userImplementation) UpdateImage(id string, imageURL string) error {
	stmt := `
    UPDATE users
    SET
        avatar = $2,
        updated_at = now()
    WHERE id = $1
    `

	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	_, err := u.DB.ExecContext(ctx, stmt, id, imageURL)
	if err != nil {
		return err
	}

	return nil
}

// ChangeRole updates the user's role.
func (u *userImplementation) ChangeRole(id, role string) error {
	stmt := `
    UPDATE users
    SET role = $2
    WHERE id = $1
    `

	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	_, err := u.DB.ExecContext(ctx, stmt, id, role)
	if err != nil {
		return err
	}

	return nil
}

// ChangeServiceRole updates the user's service role.
func (u *userImplementation) ChangeServiceRole(id, role string) error {
	stmt := `
    UPDATE users
    SET service_role = $2
    WHERE id = $1
    `

	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	_, err := u.DB.ExecContext(ctx, stmt, id, role)
	if err != nil {
		return err
	}

	return nil
}

// Search searches for user in the database using the specified filters
func (u *userImplementation) Search(filter models.Filter) ([]models.User, error) {
	stmt := `
    SELECT 
        id,
        first_name,
        last_name,
        email,
        phone_num,
        avatar,
        gender,
        date_of_birth,
        bio,
        address,
        role,
        service_role,
        is_verified,
        email_verified,
        phone_verified,
        created_at,
        updated_at
    FROM users
    WHERE
        ($1 = '' OR first_name = $1) AND
        ($2 = '' OR last_name = $2) AND
        ($3 = '' OR phone_num = $3) AND
        ($4 = '' OR role = $4) AND
        ($5 = '' OR service_role = $5) AND
        ($6 = '' OR email = $6) AND
        ($7 = '' OR gender = $7)
    ORDER BY first_name
    LIMIT $8 OFFSET $9
    `

	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	rows, err := u.DB.QueryContext(
		ctx,
		stmt,
		filter.FirstName,
		filter.LastName,
		filter.PhoneNum,
		filter.Role,
		filter.ServiceRole,
		filter.Email,
		filter.Gender,
		filter.Limit,
		filter.Offset(),
	)
	if err != nil {
		return nil, err
	}

	users := []models.User{}
	for rows.Next() {
		user := models.User{}
		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.PhoneNum,
			&user.Avatar,
			&user.Gender,
			&user.DateOfBirth,
			&user.Bio,
			&user.Address,
			&user.Role,
			&user.ServiceRole,
			&user.IsVerified,
			&user.EmailVerified,
			&user.PhoneVerified,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// Verify updates the verification status of a user.
func (u *userImplementation) Verify(id string, status bool) error {
	stmt := `
    UPDATE users
    SET is_verified = $2
    WHERE id = $1;
    `

	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	_, err := u.DB.ExecContext(ctx, stmt, id, status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.ErrRecordNotFound
		}
		return err
	}
	return nil
}

// VerifyUserEmal updates the verification status of user's email or phone number.
func (u *userImplementation) VerifyContact(id, verificationType string, status bool) error {
	stmt := ""
	switch verificationType {
	case "email":
		stmt = `
        UPDATE users
        SET email_verified = $2
        WHERE id = $1;`
	case "phone":
		stmt = `
        UPDATE users
        SET phone_verified = $2
        WHERE id = $1;`
	}

	ctx, cancel := context.WithTimeout(context.Background(), DB_QUERY_TIMEOUT)
	defer cancel()

	_, err := u.DB.ExecContext(ctx, stmt, id, status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.ErrRecordNotFound
		}
		return err
	}
	return nil
}
