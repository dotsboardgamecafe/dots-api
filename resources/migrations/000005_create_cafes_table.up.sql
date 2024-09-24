CREATE TABLE cafes (
	id serial PRIMARY KEY UNIQUE,
	cafe_code varchar(50) NOT NULL UNIQUE,
	"name" varchar(100) NOT NULL DEFAULT '',
	"address" text NULL DEFAULT '',
	"description" text NULL DEFAULT '',
	"status" varchar(50) NOT NULL, --active|inactive
	created_date timestamptz(0) NOT NULL DEFAULT NOW(),
	updated_date timestamptz(0) NULL,
	deleted_date timestamptz(0) NULL
);