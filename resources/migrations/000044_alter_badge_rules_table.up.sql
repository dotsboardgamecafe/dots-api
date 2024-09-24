-- Add column badge_rule_code
ALTER TABLE badges_rules
ADD COLUMN badge_rule_code VARCHAR(50) NOT NULL UNIQUE;