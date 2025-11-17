DO $$
BEGIN
    IF NOT EXISTS (
        SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'rule_categories'
    ) THEN
        CREATE TABLE rule_categories (
            id SERIAL PRIMARY KEY,
            name VARCHAR(100) NOT NULL UNIQUE,   -- ví dụ: Currencies, Providers...
            code VARCHAR(50) NOT NULL UNIQUE,    -- viết tắt: currencies, providers...
            created_at TIMESTAMP DEFAULT now(),
            updated_at TIMESTAMP DEFAULT now()
        );
    END IF;
END
$$;
