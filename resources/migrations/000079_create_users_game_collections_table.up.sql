CREATE TABLE IF NOT EXISTS users_game_collections(
  user_id bigint references users(id) ON DELETE CASCADE ON UPDATE CASCADE,
  game_id bigint REFERENCES games(id) ON DELETE CASCADE ON UPDATE CASCADE,
  created_date timestamp NULL DEFAULT now(),
  PRIMARY KEY(user_id, game_id)
);