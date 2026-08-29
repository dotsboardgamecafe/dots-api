-- Convert additional_info column to JSONB and backfill legacy data
ALTER TABLE rooms_participants 
ALTER COLUMN additional_info TYPE jsonb 
USING CASE 
    WHEN additional_info IS NULL OR additional_info = '' OR additional_info = 'member' 
    THEN '{"registration_type": "self_booking"}'::jsonb
    ELSE additional_info::jsonb
END;

ALTER TABLE rooms_participants 
ALTER COLUMN additional_info SET DEFAULT '{"registration_type": "self_booking"}'::jsonb;
