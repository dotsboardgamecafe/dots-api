CREATE TABLE rooms_participants(
    id serial primary key,
    room_id bigint references rooms(id) ON DELETE CASCADE ON UPDATE CASCADE,
    user_id bigint references users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    status_winner boolean default false,
    reward_point int  NOT NULL default 0,
    "status" varchar(50) not null, --active|pending, 
    created_date timestamptz(0) NOT NULL DEFAULT NOW(),
    updated_date timestamptz(0) NULL,
    CONSTRAINT room_id_and_user_id_in_rooms_participants UNIQUE(room_id ,user_id )
);