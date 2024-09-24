-- Change type data
ALTER TABLE notifications 
ADD COLUMN image_url VARCHAR(250) DEFAULT '';

-- Change column name
ALTER TABLE notifications
RENAME COLUMN notification_type TO "type";

-- Change column name
ALTER TABLE notifications
RENAME COLUMN notification_title TO "title";

-- Change column name
ALTER TABLE notifications
RENAME COLUMN notification_description TO "description";


-- Add transaction code
ALTER TABLE notifications 
ADD COLUMN transaction_code VARCHAR(50) DEFAULT '';

-- Add notification code
ALTER TABLE notifications 
ADD COLUMN notification_code VARCHAR(50) DEFAULT '';
