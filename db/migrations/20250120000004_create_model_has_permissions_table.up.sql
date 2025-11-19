DO $$
BEGIN
    IF NOT EXISTS (
        SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'model_has_permissions'
    ) THEN
        CREATE TABLE model_has_permissions (
            permission_id BIGINT NOT NULL,
            model_type VARCHAR(255) NOT NULL,
            model_id BIGINT NOT NULL,
            PRIMARY KEY (permission_id, model_id, model_type),
            CONSTRAINT fk_model_has_permissions_permission 
                FOREIGN KEY (permission_id) 
                REFERENCES permissions(id) 
                ON DELETE CASCADE
        );

        CREATE INDEX idx_model_has_permissions_model ON model_has_permissions(model_id, model_type);
        CREATE INDEX idx_model_has_permissions_permission ON model_has_permissions(permission_id);

        COMMENT ON TABLE model_has_permissions IS 'Polymorphic pivot table - links models (users) to permissions';
        COMMENT ON COLUMN model_has_permissions.model_type IS 'Model type (e.g., User)';
        COMMENT ON COLUMN model_has_permissions.model_id IS 'Model ID (user ID)';
    END IF;
END
$$;

