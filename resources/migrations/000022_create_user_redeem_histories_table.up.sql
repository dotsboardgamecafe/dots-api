CREATE TABLE IF NOT EXISTS user_redeem_histories (
  id SERIAL PRIMARY KEY,
  user_id BIGINT references users(id) ON DELETE CASCADE ON UPDATE CASCADE,
  reward_id BIGINT references rewards(id) ON DELETE CASCADE ON UPDATE CASCADE,
  users_transactions_id BIGINT NULL references users_transactions(id) ON DELETE CASCADE ON UPDATE CASCADE,
  is_used SMALLINT NOT NULL DEFAULT 0, -- 1: Used; 0: Not used
  created_date timestamptz(0) NOT NULL DEFAULT NOW(),
  updated_date timestamptz(0) NOT NULL -- Could be used for marking when the voucher be used
);