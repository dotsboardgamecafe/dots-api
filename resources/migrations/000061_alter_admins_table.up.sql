ALTER TABLE admins
ADD COLUMN IF NOT EXISTS username varchar(100) UNIQUE;
