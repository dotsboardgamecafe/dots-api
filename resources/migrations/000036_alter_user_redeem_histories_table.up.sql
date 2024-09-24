ALTER TABLE users_transactions DROP COLUMN IF EXISTS user_redeem_id; 
DROP TABLE IF EXISTS user_redeem_histories;

CREATE TABLE IF NOT EXISTS user_redeem_histories (
  id SERIAL PRIMARY KEY,
  user_id BIGINT references users(id) ON DELETE CASCADE ON UPDATE CASCADE,
  "custom_id" VARCHAR(100) UNIQUE,
  invoice_code VARCHAR(100) NOT NULL,
  invoice_amount BIGINT NOT NULL DEFAULT 0,
  "description" TEXT NOT NULL, -- Saving all invoice items
  created_date timestamptz(0) NOT NULL DEFAULT NOW(),
  updated_date timestamptz(0) NOT NULL
);

CREATE INDEX user_invoice_code_idx ON user_redeem_histories(user_id, invoice_code);
