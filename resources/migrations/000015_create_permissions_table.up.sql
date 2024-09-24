CREATE TABLE permissions ( 
id serial PRIMARY KEY UNIQUE, 
permission_code varchar(50) NOT NULL UNIQUE, 
name varchar(100) NOT NULL UNIQUE, 
route_pattern varchar(255) NOT NULL, --/v1/admins/{code}
route_method varchar(50) NOT NULL, --GET,POST,PUT,DELETE
"description" text NULL DEFAULT '', 
status varchar(50) NOT NULL, --active|inactive 
created_date timestamptz(0) NOT NULL DEFAULT NOW(), 
updated_date timestamptz(0) NULL, 
deleted_date timestamptz(0) NULL );