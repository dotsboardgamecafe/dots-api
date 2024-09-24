CREATE TABLE badges(
	id serial PRIMARY KEY UNIQUE,
	tier_code varchar(50) NOT NULL UNIQUE,
	badge_type varchar(100) NOT NULL, --retail|f&b|play,
	"name" varchar(100) NOT NULL DEFAULT '',
	image_url text NULL DEFAULT '',
	"start_date" timestamptz(0) NOT NULL,
	end_date timestamptz(0) NOT NULL,
	max_claim int NOT NULL default 0,
	reward_point int NOT NULL default 0,
	"status" varchar(50) NOT NULL, --active|inactive
	created_date timestamptz(0) NOT NULL DEFAULT NOW(),
	updated_date timestamptz(0) NULL,
	deleted_date timestamptz(0) NULL
);