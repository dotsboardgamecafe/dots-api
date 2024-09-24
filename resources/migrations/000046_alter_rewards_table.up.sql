-- Drop redeemable_point column
ALTER TABLE rewards DROP COLUMN redeemable_point;

-- Drop voucher_code column
ALTER TABLE rewards DROP COLUMN voucher_code;

-- Alter status column to VARCHAR(50)
ALTER TABLE rewards ALTER COLUMN status SET DATA TYPE VARCHAR(50);

-- Add deleted_date column
ALTER TABLE rewards ADD COLUMN deleted_date timestamptz(0) NULL;