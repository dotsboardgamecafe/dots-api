CREATE TABLE admins (
	id serial PRIMARY KEY UNIQUE,
	admin_code varchar(50) NOT NULL UNIQUE,
	email varchar(100) NOT NULL UNIQUE,
	"name" varchar(100) NULL DEFAULT '',
	"password" varchar(100) NOT NULL,
	"status" varchar(50) NOT NULL, --active|inactive
	created_date timestamptz(0) NOT NULL DEFAULT NOW(),
	updated_date timestamptz(0) NULL,
	deleted_date timestamptz(0) NULL,
	CONSTRAINT email_and_deleted_date_in_admins UNIQUE(email,deleted_date) 
);