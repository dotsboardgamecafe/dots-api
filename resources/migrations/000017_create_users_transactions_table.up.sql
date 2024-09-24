CREATE TABLE users_transactions(
  id bigserial PRIMARY KEY,
  user_id bigint references users(id) ON DELETE CASCADE ON UPDATE CASCADE,
  "data_source" varchar(100) not null default '', --rooms
  source_code varchar(100) null default '',
  transaction_code varchar(50) NOT NULL UNIQUE,
  aggregator_code varchar(100) NULL UNIQUE,
  price bigint NOT NULL default 0,
  payment_method varchar(100) NOT NULL DEFAULT '',
  payment_link varchar(255) NOT NULL DEFAULT '',
  "status" varchar(20) NOT NULL DEFAULT '', --pending
  resp_payload text NULL DEFAULT '',
  created_date timestamp NULL DEFAULT now(),
  updated_date timestamp NULL
);