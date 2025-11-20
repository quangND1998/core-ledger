DO $$
BEGIN
    -- Xóa cột metadata nếu tồn tại
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_schema = 'public' 
        AND table_name = 'rule_categories' 
        AND column_name = 'metadata'
    ) THEN
        ALTER TABLE rule_categories DROP COLUMN metadata;
    END IF;
END
$$;

