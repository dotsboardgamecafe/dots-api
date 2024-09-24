INSERT INTO permissions
(permission_code, "name", route_pattern, route_method, description, status)
VALUES
('PRMS-20240430ITYRBLUEUP','get-tournament-badges','/v1/tournament-badges/*','GET','get-tournament-badges','active'),
('PRMS-20240430AQZOAEFXAR','update-tournament-badges','/v1/tournament-badges/*','PUT','update-tournament-badges','active'),
('PRMS-20240430QQEVOSVPNI','add-tournament-badges','/v1/tournament-badges','POST','add-tournament-badges','active'),
('PRMS-20240430PMUEECHACJ','tournament-setwinner','/v1/tournaments/*/close','PUT','tournament-setwinner','active'),
('PRMS-20240430HIKIOSVPIN','invoice-history-cms','/v1/invoices/*/history','GET','invoice-history-cms','active'),
('PRMS-20240430AISYUNYUON','claim-invoice-cms','/v1/invoices/*/claim','POST','claim-invoice-cms','active'),
('PRMS-20240430SJOPQOONIC', 'member-update-badges', '/v1/users/*/badges/*', 'PUT', 'member-update-badges', 'active');
