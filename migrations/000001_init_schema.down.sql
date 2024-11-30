DROP TABLE IF EXISTS "waitlist";

DROP TABLE IF EXISTS "payment_recipient_codes";

DROP TABLE IF EXISTS "payment_access_codes";

DROP TABLE IF EXISTS "wallets";

DROP TABLE IF EXISTS "payments";

DROP TABLE IF EXISTS "notifications";

DROP TABLE IF EXISTS "rejected_bookings";

DROP TABLE IF EXISTS "bookings";

DROP TABLE IF EXISTS "service_commission";

DROP TABLE IF EXISTS "services_pricing";

DROP TABLE IF EXISTS "services";

DROP TABLE IF EXISTS "service_images";

DROP TABLE IF EXISTS "users_certifications";

DROP TABLE IF EXISTS "users_education_info";

DROP TABLE IF EXISTS "users_bank_info";

DROP TABLE IF EXISTS "contact_verifications";

DROP TABLE IF EXISTS "users";

DROP EXTENSION IF EXISTS CITEXT;

DROP FUNCTION IF EXISTS "compute_actual_price";
