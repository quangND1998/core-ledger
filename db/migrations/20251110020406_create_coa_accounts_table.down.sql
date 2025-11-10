DO $$
BEGIN
    IF EXISTS (
        SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'coa_accounts'
    ) THEN
        DROP TABLE coa_accounts;
    END IF;
END
$$;