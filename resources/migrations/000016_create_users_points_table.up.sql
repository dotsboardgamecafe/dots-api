create table users_points (
    id serial primary key,
    user_id bigint references users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    "data_source" varchar(100) not null default '', --rooms|badges
    source_code varchar(100) null default '',
    "point" bigint null default 0,
    created_date timestamptz(0) not null
);