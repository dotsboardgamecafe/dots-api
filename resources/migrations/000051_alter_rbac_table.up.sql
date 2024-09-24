ALTER TABLE users ADD role_id bigint references roles(id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE admins ADD role_id bigint references roles(id) ON DELETE CASCADE ON UPDATE CASCADE;