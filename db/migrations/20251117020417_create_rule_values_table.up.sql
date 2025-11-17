DO $$
BEGIN
    IF NOT EXISTS (
        SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'rule_values'
    ) THEN
        CREATE TABLE rule_values (
            id SERIAL PRIMARY KEY,
            category_id INT NOT NULL REFERENCES rule_categories(id) ON DELETE CASCADE,
            name VARCHAR(255) NULL,
            value VARCHAR(255) NOT NULL,
            is_delete BOOLEAN DEFAULT FALSE,
            sort_order INT DEFAULT 0,
            created_at TIMESTAMP DEFAULT now(),
            updated_at TIMESTAMP DEFAULT now()
        );
    END IF;
END
$$;
