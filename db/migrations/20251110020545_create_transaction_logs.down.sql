DO $$
BEGIN
    IF EXISTS (
        SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'transaction_logs'
    ) THEN
        DROP TABLE transaction_logs;
    END IF;
END
$$;