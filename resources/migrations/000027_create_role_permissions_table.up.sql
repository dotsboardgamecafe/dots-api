CREATE TABLE role_permissions ( 
id serial PRIMARY KEY UNIQUE, 
role_id bigint references roles(id) ON DELETE CASCADE ON UPDATE CASCADE, 
permission_id bigint references permissions(id) ON DELETE CASCADE ON UPDATE CASCADE, 
created_date timestamptz(0) NOT NULL DEFAULT NOW());