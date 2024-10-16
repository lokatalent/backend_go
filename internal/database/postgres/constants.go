package postgres

import (
	"time"
)

const DB_QUERY_TIMEOUT = 15 * time.Second

// users table constraints
const (
	duplicateEmail       = "unique_nonempty_email"
	duplicatePhone       = "unique_nonempty_phone_num"
	duplicateService     = "unique_user_id_service_type"
	duplicateBankAcctNum = "users_bank_info_account_num_key"
)
