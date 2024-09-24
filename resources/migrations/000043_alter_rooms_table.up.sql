ALTER TABLE rooms DROP minimal_participant;
ALTER TABLE rooms DROP maximum_participant;

ALTER TABLE rooms ADD instruction text NULL; -- "Used for filling 'special instruction' on UI side (special event only)";

