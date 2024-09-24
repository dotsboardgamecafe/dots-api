CREATE TABLE IF NOT EXISTS tournament_participants (
  id SERIAL PRIMARY KEY,
  user_id BIGINT references users(id) ON DELETE CASCADE ON UPDATE CASCADE,
  tournament_prize_id BIGINT NULL references tournament_prizes(id) ON DELETE CASCADE ON UPDATE CASCADE,
  created_date timestamptz(0) NOT NULL DEFAULT NOW(),
  updated_date timestamptz(0) NULL,
  CONSTRAINT tournament_prize_id_and_user_id_in_tournaments_participants UNIQUE(tournament_prize_id ,user_id )
);