
DO $$
BEGIN
    -- Kiểm tra bảng 'groups' đã tồn tại trong schema 'public' chưa
    IF NOT EXISTS (
        SELECT 1
        FROM pg_tables 
        WHERE schemaname = 'public' 
          AND tablename = 'groups'
    ) THEN
        -- Tạo bảng 'groups'
        CREATE TABLE groups (
            id BIGSERIAL PRIMARY KEY,                  -- ID tự tăng
            parent_id BIGINT REFERENCES groups(id) ON DELETE NO ACTION, -- Node cha, cho phép NULL nếu là Type
            type VARCHAR(20) NOT NULL,                -- Level của node: type, group, subgroup
            code VARCHAR(50) NOT NULL,                -- Mã node, unique trong cùng parent
            name VARCHAR(100) NOT NULL,               -- Tên hiển thị
            description TEXT,                         -- Mô tả chi tiết
            display_order INT DEFAULT 0,              -- Thứ tự hiển thị trong cùng parent
            created_at TIMESTAMP DEFAULT now(),       -- Thời gian tạo
            updated_at TIMESTAMP DEFAULT now(),       -- Thời gian cập nhật
            UNIQUE(parent_id, code)                   -- Code phải unique trong cùng parent
        );

        -- Comment bảng
        COMMENT ON TABLE groups IS 'Bảng lưu cấu trúc Code Rule Tree: type, group, subgroup';
        -- Comment các cột quan trọng
        COMMENT ON COLUMN groups.parent_id IS 'Node cha (NULL nếu là Type)';
        COMMENT ON COLUMN groups.type IS 'Level của node: type, group, subgroup';
        COMMENT ON COLUMN groups.code IS 'Mã node, duy nhất trong cùng parent';
        COMMENT ON COLUMN groups.display_order IS 'Thứ tự hiển thị khi render tree';
    END IF;
END
$$;