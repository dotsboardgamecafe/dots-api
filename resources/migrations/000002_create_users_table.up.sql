CREATE TABLE users (
	id serial PRIMARY KEY UNIQUE,
	user_code varchar(50) NOT NULL UNIQUE,
	email varchar(100) NOT NULL UNIQUE,
	username varchar(100) NOT NULL UNIQUE,
	phone_number varchar(15) NOT NULL UNIQUE,
	fullname varchar(100) NULL DEFAULT '',
	image_url text NULL DEFAULT '',
	latest_point int not null default 0,
	latest_tier varchar(100) not null default '',
	"password" varchar(100) NOT NULL,
	x_player varchar(255) NULL DEFAULT '',
	status_verification BOOLEAN DEFAULT FALSE,
	"status" varchar(50) NOT NULL, --active|inactive
	created_date timestamptz(0) NOT NULL DEFAULT NOW(),
	updated_date timestamptz(0) NULL,
	deleted_date timestamptz(0) NULL
);