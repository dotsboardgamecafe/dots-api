ALTER TABLE rooms 
  ADD maximum_participant int NOT NULL DEFAULT 0, -- "Used for player slot"
  ADD image_url VARCHAR(100)  NULL DEFAULT '';
