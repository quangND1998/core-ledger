DO $$
BEGIN
    IF EXISTS (
        SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'journals'
    ) THEN
        DROP TABLE journals;
    END IF;
END
$$;