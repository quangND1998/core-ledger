DO $$
BEGIN
    -- Xóa cột metadata nếu tồn tại
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_schema = 'public' 
        AND table_name = 'account_rule_option_steps' 
        AND column_name = 'metadata'
    ) THEN
        ALTER TABLE account_rule_option_steps DROP COLUMN metadata;
    END IF;
END
$$;

