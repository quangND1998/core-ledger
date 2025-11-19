DO $$
BEGIN
    IF NOT EXISTS (
        SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'model_has_roles'
    ) THEN
        CREATE TABLE model_has_roles (
            role_id BIGINT NOT NULL,
            model_type VARCHAR(255) NOT NULL,
            model_id BIGINT NOT NULL,
            PRIMARY KEY (role_id, model_id, model_type),
            CONSTRAINT fk_model_has_roles_role 
                FOREIGN KEY (role_id) 
                REFERENCES roles(id) 
                ON DELETE CASCADE
        );

        CREATE INDEX idx_model_has_roles_model ON model_has_roles(model_id, model_type);
        CREATE INDEX idx_model_has_roles_role ON model_has_roles(role_id);

        COMMENT ON TABLE model_has_roles IS 'Polymorphic pivot table - links models (users) to roles';
        COMMENT ON COLUMN model_has_roles.model_type IS 'Model type (e.g., User)';
        COMMENT ON COLUMN model_has_roles.model_id IS 'Model ID (user ID)';
    END IF;
END
$$;

