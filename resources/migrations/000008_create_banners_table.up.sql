CREATE TABLE banners (
	id serial PRIMARY KEY UNIQUE,
	banner_code varchar(50) NOT NULL UNIQUE,
	banner_type varchar(100) NOT NULL, --announcement|promo|game|tournament_prize|beverage
	title varchar(100) NOT NULL DEFAULT '',
	"description" text NULL DEFAULT '',
	image_url text NULL DEFAULT '',
	"status" varchar(50) NOT NULL, --publish|unpublish
	created_date timestamptz(0) NOT NULL DEFAULT NOW(),
	updated_date timestamptz(0) NULL,
	deleted_date timestamptz(0) NULL
);