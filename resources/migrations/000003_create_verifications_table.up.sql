CREATE TABLE verifications(
	id SERIAL PRIMARY KEY,
	actor_type VARCHAR(10) NOT NULL, --admin|user
	verification_type VARCHAR(50) NOT NULL, --verify_registration||verify_forgotpassword
	email VARCHAR(100) NOT NULL,
	token VARCHAR(255) NOT NULL,
	is_used BOOLEAN DEFAULT FALSE,
	expired_date timestamptz(0) NOT NULL,
	created_date timestamptz(0) NOT NULL DEFAULT NOW(),
	updated_date timestamptz(0) NULL
);