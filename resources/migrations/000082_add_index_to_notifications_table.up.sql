CREATE INDEX IF NOT EXISTS notifications_code_idx ON notifications (notification_code);
CREATE INDEX IF NOT EXISTS notifications_transaction_code_and_type_idx
    ON notifications (transaction_code, "type");