-- Adding columns image_url, prizes_img_url, description, start_time, end_time, 
ALTER TABLE tournaments
ADD COLUMN image_url VARCHAR(255) NULL,
ADD COLUMN prizes_img_url VARCHAR(255) NULL,
ADD COLUMN start_time TIME NULL,
ADD COLUMN end_time TIME NULL;

-- Changing column difficulty to level
ALTER TABLE tournaments
RENAME COLUMN difficulty TO level;

ALTER TABLE tournaments
RENAME COLUMN reward_point TO participant_vp;


-- Dropping column in table tournaments
ALTER TABLE tournaments
DROP COLUMN banner_id,
DROP COLUMN participant_badge_id,
DROP COLUMN tournament_badge_id;


-- Changing data type of columns start_date and end_date to DATE
ALTER TABLE tournaments
ALTER COLUMN "start_date" TYPE DATE,
ALTER COLUMN end_date TYPE DATE;
