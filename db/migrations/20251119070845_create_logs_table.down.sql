DO $$
BEGIN
    IF EXISTS (
        SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'logs'
    ) THEN
        DROP TABLE logs;
    END IF;
END
$$;