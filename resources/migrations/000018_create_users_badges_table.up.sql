CREATE TABLE users_badges(
  id bigserial PRIMARY KEY,
  user_id bigint references users(id) ON DELETE CASCADE ON UPDATE CASCADE,
  badge_id bigint REFERENCES badges(id) ON DELETE CASCADE ON UPDATE CASCADE,
  created_date timestamp NULL DEFAULT now(),
  CONSTRAINT user_id_and_badge_id_in_tournaments  UNIQUE(user_id ,badge_id)
);