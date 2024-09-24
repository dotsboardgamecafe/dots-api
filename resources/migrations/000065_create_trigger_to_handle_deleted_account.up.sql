ALTER TABLE users DROP CONSTRAINT users_email_key;
ALTER TABLE users DROP CONSTRAINT users_phone_number_key;
ALTER TABLE users DROP CONSTRAINT users_username_key;

-- create trigger 
CREATE OR REPLACE FUNCTION check_unique_deleted_user()
RETURNS TRIGGER AS
$$
BEGIN
    -- Check if the email is being updated
    IF NEW.email IS NOT NULL AND OLD.email <> NEW.email THEN
        -- Prevent updating to a duplicate email unless the existing data is soft-deleted
        IF EXISTS (
            SELECT 1
            FROM public.users
            WHERE email = NEW.email AND (deleted_date IS NULL )
        ) THEN
            RAISE EXCEPTION 'users_email_key, Cannot update to duplicate email: %', NEW.email;
        END IF;
    END IF;

    -- Check if the phone_number is being updated
    IF NEW.phone_number IS NOT NULL AND OLD.phone_number <> NEW.phone_number THEN
        -- Prevent updating to a duplicate phone_number unless the existing data is soft-deleted
        IF EXISTS (
            SELECT 1
            FROM public.users
            WHERE phone_number = NEW.phone_number AND (deleted_date IS NULL )
        ) THEN
            RAISE EXCEPTION 'users_phone_number_key, Cannot update to duplicate phone number: %', NEW.phone_number;
        END IF;
    END IF;

    -- Check if the username is being updated
    IF NEW.username IS NOT NULL AND OLD.username <> NEW.username THEN
        -- Prevent updating to a duplicate username unless the existing data is soft-deleted
        IF EXISTS (
            SELECT 1
            FROM public.users
            WHERE username = NEW.username AND (deleted_date IS NULL )
        ) THEN
            RAISE EXCEPTION 'users_username_key, Cannot update to duplicate username: %', NEW.username;
        END IF;
    END IF;

    -- Check if the insert operation is being performed
    IF TG_OP = 'INSERT' THEN
        -- Prevent inserting duplicate email, phone_number, or username unless the existing data is soft-deleted
       	IF EXISTS (
	        SELECT 1
	        FROM public.users
	        WHERE email = NEW.email AND (deleted_date IS null )
	    ) THEN
	        RAISE EXCEPTION 'users_email_key, Duplicate  email found: %', NEW.email;
	    ELSIF EXISTS (
	        SELECT 1
	        FROM public.users
	        WHERE phone_number = NEW.phone_number AND (deleted_date IS null  )
	    ) THEN
	        RAISE EXCEPTION 'users_phone_number_key, Duplicate phone number found: %', NEW.phone_number;
	    ELSIF EXISTS (
	        SELECT 1
	        FROM public.users
	        WHERE username = NEW.username AND (deleted_date IS null  )
	    ) THEN
	        RAISE EXCEPTION 'users_username_key, Duplicate username found: %', NEW.username;
	    END IF;
    END IF;

    RETURN NEW;
END;
$$
LANGUAGE plpgsql;


CREATE TRIGGER trigger_check_unique_deleted_user
BEFORE INSERT OR UPDATE ON public.users
FOR EACH ROW
EXECUTE FUNCTION check_unique_deleted_user();