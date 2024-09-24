-- Menambah kolom is_claim
ALTER TABLE users_badges
ADD COLUMN is_claim boolean default false;
