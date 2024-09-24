ALTER TABLE user_redeem_histories 
  ADD COLUMN IF NOT EXISTS "custom_id" VARCHAR(100) UNIQUE;
