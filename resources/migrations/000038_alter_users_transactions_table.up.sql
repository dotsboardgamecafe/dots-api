CREATE INDEX user_trx_data_source_idx ON users_transactions(user_id, data_source);
CREATE INDEX user_trx_data_source_code_idx ON users_transactions(user_id, data_source, source_code);
