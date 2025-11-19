DO $$
BEGIN
    IF NOT EXISTS (
        SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'users'
    ) THEN
        CREATE TABLE users (
            id BIGSERIAL PRIMARY KEY,
            email VARCHAR(255) NOT NULL UNIQUE,
            password VARCHAR(255) NOT NULL,
            full_name VARCHAR(255),
            guard_name VARCHAR(50) DEFAULT 'web',
            created_at TIMESTAMP DEFAULT NOW(),
            updated_at TIMESTAMP DEFAULT NOW()
        );

        CREATE INDEX idx_users_email ON users(email);
        CREATE INDEX idx_users_guard_name ON users(guard_name);

        COMMENT ON TABLE users IS 'Users table for authentication and authorization';
        COMMENT ON COLUMN users.guard_name IS 'Guard name for multi-guard support (e.g., web, api)';
    END IF;
END
$$;

