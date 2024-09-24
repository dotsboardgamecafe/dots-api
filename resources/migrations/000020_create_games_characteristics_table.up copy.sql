CREATE TABLE games_characteristics (
	id serial PRIMARY KEY UNIQUE,
	game_id bigint references cafes(id) ON DELETE CASCADE ON UPDATE CASCADE,
	characteristic_left varchar(100) NOT NULL,
	characteristic_right varchar(100) NOT NULL,
	"value" int NOT NULL
);