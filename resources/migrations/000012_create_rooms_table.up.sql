CREATE TABLE rooms (
	id serial PRIMARY KEY UNIQUE,
	game_master_id bigint references admins(id) ON DELETE CASCADE ON UPDATE CASCADE,
	game_id bigint references games(id) ON DELETE CASCADE ON UPDATE CASCADE,
	room_code varchar(50) NOT NULL UNIQUE,
	room_type varchar(50) NOT NULL, --normal|special_event
	"name" varchar(100) NOT NULL DEFAULT '',
	"description" text NULL DEFAULT '',
	difficulty varchar(50) NULL DEFAULT '', --easy|medium|hard
	"start_date" timestamptz(0) NOT NULL,
	end_date timestamptz(0) NOT NULL,
	minimal_participant int NOT NULL default 0,
	maximum_participantint int NOT NULL default 0,
	booking_price int NOT NULL default 0,
	reward_point int NOT NULL default 0,
	instagram_link varchar(100) NULL DEFAULT '',
	"status" varchar(50) NOT NULL, --publish|unpublish
	created_date timestamptz(0) NOT NULL DEFAULT NOW(),
	updated_date timestamptz(0) NULL,
	deleted_date timestamptz(0) NULL
);