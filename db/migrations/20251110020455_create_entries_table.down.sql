DO $$
BEGIN
    IF EXISTS (
        SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'entries'
    ) THEN
        DROP TABLE entries;
    END IF;
END
$$;