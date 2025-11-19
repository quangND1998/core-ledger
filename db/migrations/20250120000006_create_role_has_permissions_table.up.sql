DO $$
BEGIN
    IF NOT EXISTS (
        SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'role_has_permissions'
    ) THEN
        CREATE TABLE role_has_permissions (
            permission_id BIGINT NOT NULL,
            role_id BIGINT NOT NULL,
            PRIMARY KEY (permission_id, role_id),
            CONSTRAINT fk_role_has_permissions_permission 
                FOREIGN KEY (permission_id) 
                REFERENCES permissions(id) 
                ON DELETE CASCADE,
            CONSTRAINT fk_role_has_permissions_role 
                FOREIGN KEY (role_id) 
                REFERENCES roles(id) 
                ON DELETE CASCADE
        );

        CREATE INDEX idx_role_has_permissions_permission ON role_has_permissions(permission_id);
        CREATE INDEX idx_role_has_permissions_role ON role_has_permissions(role_id);

        COMMENT ON TABLE role_has_permissions IS 'Pivot table - links roles to permissions';
    END IF;
END
$$;

