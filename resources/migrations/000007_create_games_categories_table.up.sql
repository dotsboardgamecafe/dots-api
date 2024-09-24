CREATE TABLE games_categories (
	id serial PRIMARY KEY UNIQUE,
	game_id bigint references cafes(id) ON DELETE CASCADE ON UPDATE CASCADE,
	category_name varchar(100) NOT NULL, --Negotiation, Engine Building, Dice Rolling masukin ke table settings
	CONSTRAINT role_id_and_category_name_in_games_categories UNIQUE(game_id , category_name ) 
);