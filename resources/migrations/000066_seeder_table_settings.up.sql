-- Seeding settings data
INSERT INTO public.settings (setting_code,set_group,set_key,set_label,set_order,content_type,content_value,is_active,created_date,updated_date) VALUES
	('SET-000001','gender','male','male',0,'string','MALE',true,'2023-05-04 11:57:10+07',NULL),
	('SET-000002','gender','female','female',1,'string','FEMALE',true,'2023-05-04 11:57:10+07',NULL),
    ('SET-000003', 'badge_condition', 'total_spend', 'Total spend for all purchased', 4, 'string', 'total_spend', true, '2024-05-13 00:00:00', NULL),
    ('SET-000004', 'badge_condition', 'spesific_board_game_category', 'Number of plays with spesific game', 5, 'string', 'spesific_board_game_category', true, '2024-05-13 00:00:00', NULL),
    ('SET-000005', 'badge_condition', 'time_limit', 'Number of plays in time limit', 6, 'string', 'time_limit', true, '2024-05-13 00:00:00', NULL),
    ('SET-000006', 'badge_condition', 'life_time', 'Number of plays in life time', 7, 'string', 'life_time', true, '2024-05-13 00:00:00', NULL),
    ('SET-000007', 'badge_condition', 'seasonal', 'Seasonal (holiday)', 8, 'string', 'seasonal', true, '2024-05-13 00:00:00', NULL),
    ('SET-000008', 'badge_condition', 'tournament', 'Winner Tournament', 0, 'string', 'tournament', true, '2024-05-13 00:00:00', NULL);

-- Seeding tiers data
INSERT INTO tiers (tier_code,name,min_point,max_point,description,created_date,updated_date) VALUES
    ('TIER-001','Novice',0,200,'This is Novice Tier',NOW(),NOW()),
    ('TIER-002','Intermediate',201,500,'This is Intermediate Tier',NOW(),NOW()),
    ('TIER-003','Master',501,1000,'This is Master Tier',NOW(),NOW()),
    ('TIER-004','Legend',1001,10000,'This is Legend Tier',NOW(),NOW())