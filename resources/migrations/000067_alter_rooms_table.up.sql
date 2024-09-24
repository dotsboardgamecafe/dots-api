-- Renaming first
ALTER TABLE rooms
RENAME COLUMN location_id TO location_city;

-- Altering later (drop index & change type)
ALTER TABLE rooms
DROP CONSTRAINT rooms_location_id_fkey,
ALTER COLUMN location_city TYPE VARCHAR(50);

-- Altering table tournament
ALTER TABLE tournaments
ADD COLUMN IF NOT EXISTS location_city VARCHAR(50);

-- Add Indexing on column location_city
CREATE INDEX rooms_loc_province_idx ON rooms(location_city);
CREATE INDEX tournaments_loc_province_idx ON tournaments(location_city);