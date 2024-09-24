ALTER TABLE user_redeem_histories
ADD COLUMN IF NOT EXISTS redeem_information JSONB,
ADD COLUMN IF NOT EXISTS requested_platform VARCHAR(10);
