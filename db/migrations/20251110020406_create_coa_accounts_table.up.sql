DO $$
BEGIN
    IF NOT EXISTS (
        SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'coa_accounts'
    ) THEN
        CREATE TABLE coa_accounts (
            id BIGSERIAL PRIMARY KEY,
            code VARCHAR(128) NOT NULL,
            name VARCHAR(256) NOT NULL,
            type VARCHAR(16) NOT NULL CHECK (type IN ('ASSET','LIAB','EQUITY','REV','EXP')),
            currency CHAR(8) NOT NULL,
            parent_id BIGINT NULL,
            status VARCHAR(16) DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE','INACTIVE')),
            provider VARCHAR(64) NULL,
            network VARCHAR(32) NULL,
            tags JSONB NULL,
            metadata JSONB NULL,
            created_at TIMESTAMP DEFAULT NOW(),
            updated_at TIMESTAMP DEFAULT NOW(),

            CONSTRAINT uniq_code_currency UNIQUE (code, currency),
            CONSTRAINT fk_parent FOREIGN KEY (parent_id)
                REFERENCES coa_accounts(id)
                ON DELETE SET NULL
                ON UPDATE CASCADE
        );

        CREATE INDEX idx_coa_parent ON coa_accounts(parent_id);
        CREATE INDEX idx_coa_type ON coa_accounts(type);
        CREATE INDEX idx_coa_provider ON coa_accounts(provider);
        CREATE INDEX idx_coa_network ON coa_accounts(network);

        -- ====== Thêm mô tả (comment) ======
        COMMENT ON TABLE coa_accounts IS 'Chart of Accounts (COA) - Danh mục tài khoản kế toán. Dùng cho cấu trúc tài chính.';
        COMMENT ON COLUMN coa_accounts.id IS 'Khóa chính (Primary Key)';
        COMMENT ON COLUMN coa_accounts.code IS 'Mã tài khoản kế toán (duy nhất trong mỗi loại tiền tệ)';
        COMMENT ON COLUMN coa_accounts.name IS 'Tên tài khoản kế toán hiển thị';
        COMMENT ON COLUMN coa_accounts.type IS 'Loại tài khoản kế toán: ASSET (tài sản), LIAB (nợ phải trả), EQUITY (vốn chủ), REV (doanh thu), EXP (chi phí)';
        COMMENT ON COLUMN coa_accounts.currency IS 'Mã loại tiền tệ (ví dụ: USD, VND...)';
        COMMENT ON COLUMN coa_accounts.parent_id IS 'Tài khoản cha (nếu là tài khoản con)';
        COMMENT ON COLUMN coa_accounts.status IS 'Trạng thái: ACTIVE hoặc INACTIVE';
        COMMENT ON COLUMN coa_accounts.provider IS 'Tên nhà cung cấp liên quan (ví dụ: ngân hàng, ví điện tử)';
        COMMENT ON COLUMN coa_accounts.network IS 'Mạng lưới hoặc hệ thống liên quan (ví dụ: blockchain, kênh thanh toán)';
        COMMENT ON COLUMN coa_accounts.tags IS 'Danh sách tag (JSONB) để phân loại thêm';
        COMMENT ON COLUMN coa_accounts.metadata IS 'Dữ liệu phụ bổ sung (JSONB)';
        COMMENT ON COLUMN coa_accounts.created_at IS 'Ngày tạo bản ghi';
        COMMENT ON COLUMN coa_accounts.updated_at IS 'Ngày cập nhật bản ghi cuối cùng';
    END IF;
END
$$;
