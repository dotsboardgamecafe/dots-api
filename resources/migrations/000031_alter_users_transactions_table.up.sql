ALTER TABLE users_transactions 
ADD user_redeem_id BIGINT NULL REFERENCES user_redeem_histories(id) ON DELETE CASCADE ON UPDATE CASCADE;
