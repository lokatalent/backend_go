package postgres

import (
	"time"
)

const DB_QUERY_TIMEOUT = 15 * time.Second

// users table constraints
const (
	duplicateEmail = "users_email_key"
)
