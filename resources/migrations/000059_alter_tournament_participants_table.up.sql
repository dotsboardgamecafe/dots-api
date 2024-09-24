ALTER TABLE tournament_participants ADD COLUMN IF NOT EXISTS position int default 0;
ALTER TABLE tournament_participants ADD COLUMN IF NOT EXISTS transaction_code varchar(50);