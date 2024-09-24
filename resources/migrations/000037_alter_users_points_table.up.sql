CREATE INDEX user_source_code_idx ON users_points(user_id, source_code);
CREATE INDEX user_data_source_code_idx ON users_points(user_id, data_source, source_code);
