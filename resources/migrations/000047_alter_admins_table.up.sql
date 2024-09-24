ALTER TABLE admins ADD phone_number varchar(15) NULL;
ALTER TABLE admins ADD CONSTRAINT admins_phone_number_unique UNIQUE (phone_number);
