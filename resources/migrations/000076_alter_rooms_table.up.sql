-- Changing data type of columns start_date and end_date to DATE
ALTER TABLE rooms
ALTER COLUMN "start_date" TYPE DATE,
ALTER COLUMN end_date TYPE DATE;

ALTER TABLE rooms
ADD COLUMN start_time TIME NULL,
ADD COLUMN end_time TIME NULL;