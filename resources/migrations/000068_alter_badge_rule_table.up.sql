-- Step 1: Remove the default value
ALTER TABLE badges_rules
ALTER COLUMN "value" DROP DEFAULT;

-- Step 2: Alter the column type to jsonb
ALTER TABLE badges_rules
ALTER COLUMN "value" TYPE jsonb USING CASE
    WHEN "value" = '' THEN 'null'::jsonb
    ELSE "value"::jsonb
END;

-- Step 3: Set a new default value
ALTER TABLE badges_rules
ALTER COLUMN "value" SET DEFAULT 'null'::jsonb;
