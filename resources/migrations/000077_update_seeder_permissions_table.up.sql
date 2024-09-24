INSERT INTO permissions
(permission_code, "name", route_pattern, route_method, description, status)
VALUES
('PRMS-20240728ITYRBLUEUP','get-game-mechanics','/v1/game-mechanics','GET','get-game-mechanics','active'),
('PRMS-20240728AQZOAEFXAR','get-detail-game-mechanics','/v1/game-mechanics/*','GET','get-detail-game-mechanics','active'),
('PRMS-20240728QQEVOSVPNI','add-game-mechanics','/v1/game-mechanics','POST','add-game-mechanics','active'),
('PRMS-20240728PMUEECHACJ','update-game-mechanics','/v1/game-mechanics/*','PUT','update-game-mechanics','active'),
('PRMS-20240728HIKIOSVPIN','delete-game-mechanics','/v1/game-mechanics/*','DELETE','delete-game-mechanics','active'),
('PRMS-20240728ITYRBLUEUA','get-game-types','/v1/game-types','GET','get-game-types','active'),
('PRMS-20240728AQZOAEFXAB','get-detail-game-types','/v1/game-types/*','GET','get-detail-game-types','active'),
('PRMS-20240728QQEVOSVPNC','add-game-types','/v1/game-types','POST','add-game-types','active'),
('PRMS-20240728PMUEECHACD','update-game-types','/v1/game-types/*','PUT','update-game-types','active'),
('PRMS-20240728HIKIOSVPIE','delete-game-types','/v1/game-types/*','DELETE','delete-game-types','active'),
('PRMS-20240728HIKIOSVPIG','update-users','/v1/users/*','PUT','update-users','active');
