DO $$
BEGIN
    IF EXISTS (
        SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'blacklisted_tokens'
    ) THEN
        DROP TABLE IF EXISTS blacklisted_tokens CASCADE;
    END IF;
END
$$;

