-- citext extension for case-insensitive texts
CREATE
    EXTENSION IF NOT EXISTS CITEXT;

CREATE TABLE IF NOT EXISTS "users" (
  "id"			UUID PRIMARY KEY NOT NULL,
  "first_name"	TEXT NOT NULL,
  "last_name"	TEXT NOT NULL,
  "email"		CITEXT NOT NULL DEFAULT '',
  "phone_num"	TEXT NOT NULL DEFAULT '',
  "password"	TEXT NOT NULL DEFAULT '',
  "gender"		TEXT NOT NULL DEFAULT '',
  "date_of_birth"	DATE NOT NULL DEFAULT '0001-01-01',
  "bio"				TEXT NOT NULL DEFAULT '',
  "address"			TEXT NOT NULL DEFAULT '',
  "avatar"			TEXT NOT NULL DEFAULT '',
  "role"			TEXT NOT NULL DEFAULT 'regular',
  "service_role"	TEXT NOT NULL DEFAULT '',
  "is_verified"		BOOLEAN NOT NULL DEFAULT false,
  "email_verified"	BOOLEAN NOT NULL DEFAULT false,
  "phone_verified"	BOOLEAN NOT NULL DEFAULT false,
  "created_at"		TIMESTAMPTZ NOT NULL DEFAULT 'now()',
  "updated_at"		TIMESTAMPTZ NOT NULL DEFAULT 'now()'
);

CREATE TABLE IF NOT EXISTS "contact_verifications" (
  "id"				UUID PRIMARY KEY NOT NULL,
  "user_id"			UUID NOT NULL,
  "code"			INT NOT NULL,
  "contact_type"	TEXT NOT NULL,
  "created_at"		TIMESTAMPTZ NOT NULL DEFAULT 'now()',
  "expires_at"		TIMESTAMPTZ NOT NULL DEFAULT now() + INTERVAL '10 minutes'
);

CREATE TABLE IF NOT EXISTS "users_bank_info" (
  "user_id"			UUID PRIMARY KEY NOT NULL,
  "bank_name"		TEXT NOT NULL,
  "account_name"	TEXT NOT NULL,
  "account_num"		TEXT UNIQUE NOT NULL,
  "created_at"		TIMESTAMPTZ NOT NULL DEFAULT 'now()',
  "updated_at"		TIMESTAMPTZ NOT NULL DEFAULT 'now()'
);

CREATE TABLE IF NOT EXISTS "users_education_info" (
  "user_id"		UUID PRIMARY KEY NOT NULL,
  "institute"	TEXT NOT NULL,
  "degree"		TEXT NOT NULL,
  "discipline"	TEXT NOT NULL,
  "start"		DATE NOT NULL DEFAULT '0001-01-01',
  "finish"		DATE NOT NULL DEFAULT 'now()',
  "created_at"	TIMESTAMPTZ NOT NULL DEFAULT 'now()',
  "updated_at"	TIMESTAMPTZ NOT NULL DEFAULT 'now()'
);

CREATE TABLE IF NOT EXISTS "users_certifications" (
  "id"		UUID PRIMARY KEY NOT NULL,
  "user_id"	UUID NOT NULL,
  "url"		TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS "service_images" (
  "id"				UUID PRIMARY KEY NOT NULL,
  "user_id"			UUID NOT NULL,
  "service_type"	TEXT NOT NULL,
  "url"				TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS "services" (
  "id"				UUID PRIMARY KEY NOT NULL,
  "user_id"			UUID NOT NULL,
  "service_type"	TEXT NOT NULL,
  "service_desc"	TEXT NOT NULL,
  "rate_per_hour"	DECIMAL(10,2) NOT NULL,
  "experience_years"	INT NOT NULL,
  "availability"		JSONB NOT NULL,
  "address"				TEXT NOT NULL,
  "created_at"			TIMESTAMPTZ NOT NULL DEFAULT 'now()',
  "updated_at"			TIMESTAMPTZ NOT NULL DEFAULT 'now()'
);

CREATE TABLE IF NOT EXISTS "services_pricing" (
  "id"				UUID PRIMARY KEY NOT NULL,
  "service_type"	TEXT NOT NULL,
  "rate_per_hour"	DECIMAL(10,2) NOT NULL,
  "created_at"		TIMESTAMPTZ NOT NULL DEFAULT 'now()',
  "updated_at"		TIMESTAMPTZ NOT NULL DEFAULT 'now()'
);

CREATE TABLE IF NOT EXISTS "service_commission" (
  "id"				UUID PRIMARY KEY NOT NULL,
  "percentage"		DECIMAL NOT NULL,
  "created_at"		TIMESTAMPTZ NOT NULL DEFAULT 'now()',
  "updated_at"		TIMESTAMPTZ NOT NULL DEFAULT 'now()'
);

-- function to get current commission percentage
CREATE OR REPLACE FUNCTION compute_actual_price(total_price DECIMAL(10,2))
RETURNS DECIMAL AS $$
BEGIN
    RETURN total_price * ( 1 - (
        SELECT percentage / 100
        FROM service_commission
        LIMIT 1
    ));
END;
$$ LANGUAGE plpgsql IMMUTABLE;

CREATE TABLE IF NOT EXISTS "bookings" (
  "id"				UUID PRIMARY KEY NOT NULL,
  "requester_id"	UUID NOT NULL,
  "provider_id"		UUID, -- made nullable for future update.
  "requester_loc"	TEXT NOT NULL,
  "service_type"	TEXT NOT NULL,
  "booking_type"	TEXT NOT NULL,
  "service_desc"	TEXT NOT NULL,
  "service_begin"	TIMESTAMPTZ NOT NULL DEFAULT 'now()',
  "service_end"		TIMESTAMPTZ NOT NULL,
  "total_price"		DECIMAL(10,2) NOT NULL,
  "actual_price"	DECIMAL(10,2) GENERATED ALWAYS AS (
  	compute_actual_price(total_price)
  ) STORED,
  "status"			TEXT NOT NULL DEFAULT 'open',
  "created_at"		TIMESTAMPTZ NOT NULL DEFAULT 'now()',
  "updated_at"		TIMESTAMPTZ NOT NULL DEFAULT 'now()'
);

CREATE UNIQUE INDEX IF NOT EXISTS unique_user_id_contact_type
	ON "contact_verifications" ("user_id", "contact_type");

CREATE UNIQUE INDEX IF NOT EXISTS unique_user_id_service_type
    ON "services" ("user_id", "service_type");

CREATE INDEX IF NOT EXISTS idx_services_service_type
    ON services(service_type);

CREATE UNIQUE INDEX IF NOT EXISTS unique_nonempty_phone_num 
	ON "users" ((phone_num)) 
	WHERE phone_num != '';

CREATE UNIQUE INDEX IF NOT EXISTS unique_nonempty_email 
	ON "users" ((email)) 
	WHERE email != '';

CREATE INDEX IF NOT EXISTS idx_services_availability
	ON services
	USING GIN (availability);

ALTER TABLE IF EXISTS "contact_verifications"
	ADD FOREIGN KEY ("user_id")
	REFERENCES "users" ("id");

ALTER TABLE IF EXISTS "users_bank_info"
	ADD FOREIGN KEY ("user_id")
	REFERENCES "users" ("id");

ALTER TABLE IF EXISTS "users_education_info"
	ADD FOREIGN KEY ("user_id")
	REFERENCES "users" ("id");

ALTER TABLE IF EXISTS "users_certifications"
	ADD FOREIGN KEY ("user_id")
	REFERENCES "users" ("id");

ALTER TABLE IF EXISTS "service_images"
	ADD FOREIGN KEY ("user_id")
	REFERENCES "users" ("id");

ALTER TABLE IF EXISTS "services"
	ADD FOREIGN KEY ("user_id")
	REFERENCES "users" ("id");

ALTER TABLE IF EXISTS "bookings"
	ADD FOREIGN KEY ("requester_id")
	REFERENCES "users" ("id");

ALTER TABLE IF EXISTS "bookings"
	ADD FOREIGN KEY ("provider_id")
	REFERENCES "users" ("id");
