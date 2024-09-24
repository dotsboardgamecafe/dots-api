create table settings (
    id serial primary key,
    setting_code varchar(100) not null unique,
    set_group varchar(50) null default '',
    set_key varchar(100) null default '',
    set_label varchar(100) null default '',
    set_order int not null default 0,
    content_type varchar(10) not null default '', -- string, json, bool, 
    content_value text not null default '',
    is_active boolean default false,
    created_date timestamptz(0) not null,
    updated_date timestamptz(0) null
);
