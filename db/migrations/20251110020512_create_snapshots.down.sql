DO $$
BEGIN
    IF EXISTS (
        SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'snapshots'
    ) THEN
        DROP TABLE snapshots;
    END IF;
END
$$;