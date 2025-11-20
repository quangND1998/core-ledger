DO $$
BEGIN
    IF EXISTS (
        SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'request_coa_accounts'
    ) THEN
        DROP TABLE IF EXISTS request_coa_accounts CASCADE;
    END IF;
END
$$;

