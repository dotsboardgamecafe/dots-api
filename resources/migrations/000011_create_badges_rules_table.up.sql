CREATE TABLE badges_rules(
	id serial PRIMARY KEY UNIQUE,
	badge_id bigint references badges(id) ON DELETE CASCADE ON UPDATE CASCADE,
	key_condition varchar(100) NOT NULL DEFAULT '', --playwithgm|spend_amount|
	value_type varchar(100) NOT NULL DEFAULT '', --boolean|integer
	"value" text NULL DEFAULT '' --true/false||0-100
);