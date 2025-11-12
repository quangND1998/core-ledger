DO $$
BEGIN
    IF EXISTS (
        SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'groups'
    ) THEN
        DROP TABLE groups;
    END IF;
END
$$;