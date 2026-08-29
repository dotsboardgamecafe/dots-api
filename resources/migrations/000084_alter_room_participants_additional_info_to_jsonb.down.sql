ALTER TABLE rooms_participants 
ALTER COLUMN additional_info TYPE varchar(50) 
USING additional_info::text;

ALTER TABLE rooms_participants 
ALTER COLUMN additional_info DROP DEFAULT;
