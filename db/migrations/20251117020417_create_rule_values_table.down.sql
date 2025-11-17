DO $$
BEGIN
    IF EXISTS (
        SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'rule_values'
    ) THEN
        DROP TABLE rule_values;
    END IF;
END
$$;