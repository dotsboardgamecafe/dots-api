CREATE TABLE games (
	id serial PRIMARY KEY UNIQUE,
	cafe_id bigint references cafes(id) ON DELETE CASCADE ON UPDATE CASCADE,
	game_code varchar(50) NOT NULL UNIQUE,
	game_type varchar(100) NOT NULL, --Board Games, Card Games, Party Games
	"name" varchar(100) NOT NULL,
	image_url text NULL DEFAULT '',
	"description" text NULL DEFAULT '',
	collection_url text NULL DEFAULT '',
	"status" varchar(50) NOT NULL, --active|inactive
	created_date timestamptz(0) NOT NULL DEFAULT NOW(),
	updated_date timestamptz(0) NULL,
	deleted_date timestamptz(0) NULL
);