DO $$
BEGIN
    IF NOT EXISTS (
        SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'permissions'
    ) THEN
        CREATE TABLE permissions (
            id BIGSERIAL PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
            guard_name VARCHAR(50) NOT NULL DEFAULT 'web',
            created_at TIMESTAMP DEFAULT NOW(),
            updated_at TIMESTAMP DEFAULT NOW(),
            CONSTRAINT unique_permission_name_guard UNIQUE (name, guard_name)
        );

        CREATE INDEX idx_permissions_name ON permissions(name);
        CREATE INDEX idx_permissions_guard_name ON permissions(guard_name);

        COMMENT ON TABLE permissions IS 'Permissions table - stores all available permissions';
        COMMENT ON COLUMN permissions.name IS 'Permission name (e.g., edit articles, delete users)';
        COMMENT ON COLUMN permissions.guard_name IS 'Guard name for multi-guard support';
    END IF;
END
$$;

