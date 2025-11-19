DO $$
BEGIN
    IF NOT EXISTS (
        SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'roles'
    ) THEN
        CREATE TABLE roles (
            id BIGSERIAL PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
            guard_name VARCHAR(50) NOT NULL DEFAULT 'web',
            created_at TIMESTAMP DEFAULT NOW(),
            updated_at TIMESTAMP DEFAULT NOW(),
            CONSTRAINT unique_role_name_guard UNIQUE (name, guard_name)
        );

        CREATE INDEX idx_roles_name ON roles(name);
        CREATE INDEX idx_roles_guard_name ON roles(guard_name);

        COMMENT ON TABLE roles IS 'Roles table - stores all available roles';
        COMMENT ON COLUMN roles.name IS 'Role name (e.g., admin, editor, writer)';
        COMMENT ON COLUMN roles.guard_name IS 'Guard name for multi-guard support';
    END IF;
END
$$;

