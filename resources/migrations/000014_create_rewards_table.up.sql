CREATE TABLE IF NOT EXISTS rewards (
  id SERIAL PRIMARY KEY,
  tier_id BIGINT references tiers(id) ON DELETE CASCADE ON UPDATE CASCADE,
  "name" VARCHAR(100) NOT NULL,
  image_url VARCHAR(100) NULL DEFAULT '',
  "description" TEXT NULL,
  redeemable_point INT NOT NULL, -- How many points/VP needed by user to claim the selected reward
  category_type VARCHAR(25) NOT NULL, -- fnb | game | tournament (?) 
  reward_code VARCHAR(50) NOT NULL UNIQUE,
  voucher_code VARCHAR(25) NOT NULL UNIQUE, -- Used for manual code & able to be redeemed
  "status" SMALLINT NOT NULL DEFAULT 0, -- 1: Active; 0: Inactive
  expired_date timestamptz(0) NOT NULL,
  created_date timestamptz(0) NOT NULL DEFAULT NOW(),
  updated_date timestamptz(0) NULL
);
