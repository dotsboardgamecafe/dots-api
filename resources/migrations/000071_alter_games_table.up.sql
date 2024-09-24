ALTER TABLE games RENAME COLUMN difficulty TO level;
ALTER TABLE games ADD COLUMN difficulty VARCHAR(50) NULL;
ALTER TABLE games ADD COLUMN admin_id int NULL;
ALTER TABLE games ADD COLUMN admin_id int NULL;
