CREATE TABLE users (
    id UUID PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    nickname VARCHAR(100) NOT NULL,
    password_hash TEXT NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    country CHAR(2) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION notify_user_changes()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'DELETE' THEN
        PERFORM pg_notify(
            'user_changes',
            json_build_object(
                'event', TG_OP,
                'id', OLD.id,
                'firstName', OLD.first_name,
                'lastName', OLD.last_name,
                'nickname', OLD.nickname,
                'email', OLD.email,
                'country', OLD.country,
                'createdAt', OLD.created_at,
                'updatedAt', OLD.updated_at
            )::text
        );
    ELSE
        PERFORM pg_notify(
            'user_changes',
            json_build_object(
                'event', TG_OP,
                'id', NEW.id,
                'firstName', NEW.first_name,
                'lastName', NEW.last_name,
                'nickname', NEW.nickname,
                'email', NEW.email,
                'country', NEW.country,
                'createdAt', NEW.created_at,
                'updatedAt', NEW.updated_at
            )::text
        );
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER user_changes_trigger
AFTER INSERT OR UPDATE OR DELETE ON users
FOR EACH ROW EXECUTE FUNCTION notify_user_changes();
