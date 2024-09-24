-- delete existing table
DROP TABLE IF EXISTS tournament_participants;

-- create new table
CREATE TABLE tournament_participants (
    id SERIAL PRIMARY KEY,
    tournament_id BIGINT REFERENCES tournaments(id) ON DELETE CASCADE ON UPDATE CASCADE,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    status_winner BOOLEAN DEFAULT FALSE,
    reward_point INT NOT NULL DEFAULT 0,
    "status" VARCHAR(50) NOT NULL, -- active|pending
    additional_info VARCHAR(50) NOT NULL,
    created_date TIMESTAMPTZ(0) NOT NULL DEFAULT NOW(),
    updated_date TIMESTAMPTZ(0) NULL,
    CONSTRAINT tournament_id_and_user_id_in_tournaments_participants UNIQUE (tournament_id, user_id)
);