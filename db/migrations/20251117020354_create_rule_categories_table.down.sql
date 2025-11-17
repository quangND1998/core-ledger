DO $$
BEGIN
    IF EXISTS (
        SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'rule_categories'
    ) THEN
        DROP TABLE rule_categories;
    END IF;
END
$$;