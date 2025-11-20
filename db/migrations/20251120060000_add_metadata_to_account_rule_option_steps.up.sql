DO $$
BEGIN
    -- Thêm cột metadata cho account_rule_option_steps nếu chưa có
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_schema = 'public' 
        AND table_name = 'account_rule_option_steps' 
        AND column_name = 'metadata'
    ) THEN
        ALTER TABLE account_rule_option_steps 
        ADD COLUMN metadata JSONB DEFAULT '{}'::jsonb;
        
        COMMENT ON COLUMN account_rule_option_steps.metadata IS 'Cấu hình metadata cho step (separator, description, ...)';
    END IF;
END
$$;

