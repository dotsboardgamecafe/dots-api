CREATE TABLE IF NOT EXISTS tournaments (
	id SERIAL PRIMARY KEY UNIQUE,
	game_id BIGINT references games(id) ON DELETE CASCADE ON UPDATE CASCADE,
	banner_id BIGINT NULL references banners(id) ON DELETE CASCADE ON UPDATE CASCADE,
	participant_badge_id BIGINT REFERENCES badges(id) ON DELETE CASCADE ON UPDATE CASCADE,
	tournament_badge_id BIGINT REFERENCES badges(id) ON DELETE CASCADE ON UPDATE CASCADE,
	tournament_code VARCHAR(50) NOT NULL UNIQUE,
	"name" VARCHAR(100) NOT NULL DEFAULT '',
	tournament_rules TEXT NULL DEFAULT '',
	difficulty VARCHAR(50) NULL DEFAULT '', --easy|medium|hard
	"start_date" timestamptz(0) NOT NULL,
	end_date timestamptz(0) NOT NULL,
	player_slot INT NOT NULL default 0,
	booking_price INT NOT NULL default 0,
	reward_point INT NOT NULL default 0,
	"status" VARCHAR(50) NOT NULL, --publish|unpublish,
	created_date timestamptz(0) NOT NULL DEFAULT NOW(),
	updated_date timestamptz(0) NULL,
	deleted_date timestamptz(0) NULL
);