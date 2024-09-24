ALTER TABLE games ADD duration int NOT NULL DEFAULT 0;
ALTER TABLE games ADD minimal_participant int NOT NULL DEFAULT 0;
ALTER TABLE games ADD maximum_participant int NOT NULL DEFAULT 0;
ALTER TABLE games ADD difficulty decimal NOT NULL DEFAULT 0;
