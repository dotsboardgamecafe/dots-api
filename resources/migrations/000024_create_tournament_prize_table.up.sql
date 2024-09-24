CREATE TABLE IF NOT EXISTS tournament_prizes (
	id SERIAL PRIMARY KEY UNIQUE,
	tournament_id BIGINT REFERENCES tournaments(id) ON DELETE CASCADE ON UPDATE CASCADE,
	prize_name VARCHAR(50) NOT NULL, -- 1st Place VP, etc.
	prize_point INT NOT NULL DEFAULT 0,
	created_date timestamptz(0) NOT NULL DEFAULT NOW(),
	updated_date timestamptz(0) NULL
);