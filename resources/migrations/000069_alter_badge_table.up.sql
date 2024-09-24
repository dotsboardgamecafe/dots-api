-- Altering table tournament
ALTER TABLE badges
ADD COLUMN IF NOT EXISTS parent_code VARCHAR(50);