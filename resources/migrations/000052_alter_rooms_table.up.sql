ALTER TABLE rooms ADD COLUMN IF NOT EXISTS location_id BIGINT references cafes(id) ON DELETE CASCADE ON UPDATE CASCADE;
