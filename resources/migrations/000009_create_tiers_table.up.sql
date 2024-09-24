CREATE TABLE tiers(
	id serial PRIMARY KEY UNIQUE,
	tier_code varchar(50) NOT NULL UNIQUE,
	"name" varchar(100) NOT NULL DEFAULT '',
	min_point int NOT NULL default 0,
	max_point int NOT NULL default 0,
	"description" text NULL DEFAULT '',
	created_date timestamptz(0) NOT NULL DEFAULT NOW(),
	updated_date timestamptz(0) NULL,
	deleted_date timestamptz(0) NULL
);