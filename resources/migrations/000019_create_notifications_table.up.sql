CREATE TABLE notifications (
  id bigserial NOT NULL PRIMARY KEY,
  receiver_source varchar(100) not null default '', --admins|users  
  receiver_code varchar(100) null,
  notification_type VARCHAR(50) NOT NULL,
  notification_title VARCHAR(100) NOT NULL,
  notification_description TEXT NULL,
  status_read BOOLEAN DEFAULT false,
  created_date timestamp NULL DEFAULT now(),
  updated_date timestamp NULL
);