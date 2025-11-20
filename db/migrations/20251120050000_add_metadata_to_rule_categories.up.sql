DO $$
BEGIN
    -- Thêm cột metadata cho rule_categories nếu chưa có
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_schema = 'public' 
        AND table_name = 'rule_categories' 
        AND column_name = 'metadata'
    ) THEN
        ALTER TABLE rule_categories 
        ADD COLUMN metadata JSONB DEFAULT '{}'::jsonb;
        
        COMMENT ON COLUMN rule_categories.metadata IS 'Cấu hình metadata cho category rule (separator, description, ...)';
    END IF;
END
$$;

