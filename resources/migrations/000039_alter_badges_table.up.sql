-- Menambah kolom badge_category
ALTER TABLE badges
ADD COLUMN badge_category VARCHAR(100) NOT NULL;

-- Menghapus kolom end_date
ALTER TABLE badges
DROP COLUMN IF EXISTS end_date;

-- Menghapus kolom max_claim
ALTER TABLE badges
DROP COLUMN IF EXISTS max_claim;

-- Menghapus kolom badge_type
ALTER TABLE badges
DROP COLUMN IF EXISTS badge_type;

-- Menghapus kolom end_date
ALTER TABLE badges
DROP COLUMN IF EXISTS "start_date";

ALTER TABLE badges 
DROP COLUMN IF EXISTS tier_code;