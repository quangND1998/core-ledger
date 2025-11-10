DO $$
BEGIN
    IF NOT EXISTS (
        SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'snapshots'
    ) THEN
        CREATE TABLE snapshots (
            id BIGSERIAL PRIMARY KEY,
            as_of_date DATE NOT NULL,
            account_id BIGINT NOT NULL REFERENCES coa_accounts(id),
            account_code VARCHAR(128) NOT NULL,
            currency CHAR(8) NOT NULL,
            opening_balance DECIMAL(28,8) NOT NULL,
            debit_total DECIMAL(28,8) NOT NULL,
            credit_total DECIMAL(28,8) NOT NULL,
            movement DECIMAL(28,8) NOT NULL,
            closing_balance DECIMAL(28,8) NOT NULL,
            entry_count INT NOT NULL,
            ledger_code VARCHAR(32),
            tenant_id VARCHAR(36),
            status VARCHAR(16) DEFAULT 'DRAFT' CHECK (status IN ('DRAFT','LOCKED')),
            hash CHAR(64) NOT NULL,
            meta JSONB,
            created_at TIMESTAMP DEFAULT NOW(),
            created_by VARCHAR(64)
        );

        -- Indexes
        CREATE INDEX idx_snapshots_as_of_date ON snapshots(as_of_date);
        CREATE INDEX idx_snapshots_account_id ON snapshots(account_id);
        CREATE INDEX idx_snapshots_account_code ON snapshots(account_code);
        CREATE INDEX idx_snapshots_tenant_id ON snapshots(tenant_id);
        CREATE INDEX idx_snapshots_date_account ON snapshots(as_of_date, account_id);
        CREATE INDEX idx_snapshots_date_code_tenant ON snapshots(as_of_date, account_code, tenant_id);

        -- Comments
        COMMENT ON TABLE snapshots IS 'Snapshot số dư tài khoản hàng ngày (EOD) từ journal lines';
        COMMENT ON COLUMN snapshots.id IS 'Khóa chính';
        COMMENT ON COLUMN snapshots.as_of_date IS 'Ngày chốt (End Of Day)';
        COMMENT ON COLUMN snapshots.account_id IS 'Khóa ngoại tham chiếu tài khoản COA';
        COMMENT ON COLUMN snapshots.account_code IS 'Mã tài khoản cache để đọc nhanh';
        COMMENT ON COLUMN snapshots.currency IS 'Mã tiền tệ';
        COMMENT ON COLUMN snapshots.opening_balance IS 'Số dư đầu ngày';
        COMMENT ON COLUMN snapshots.debit_total IS 'Tổng Nợ trong ngày';
        COMMENT ON COLUMN snapshots.credit_total IS 'Tổng Có trong ngày';
        COMMENT ON COLUMN snapshots.movement IS 'Chênh lệch Nợ - Có';
        COMMENT ON COLUMN snapshots.closing_balance IS 'Số dư cuối ngày';
        COMMENT ON COLUMN snapshots.entry_count IS 'Số dòng journal trong ngày';
        COMMENT ON COLUMN snapshots.ledger_code IS 'Mã sổ (GL/Sub…)';
        COMMENT ON COLUMN snapshots.tenant_id IS 'Tenant nếu đa tenant';
        COMMENT ON COLUMN snapshots.status IS 'Trạng thái snapshot';
        COMMENT ON COLUMN snapshots.hash IS 'Hash kiểm toán';
        COMMENT ON COLUMN snapshots.meta IS 'Dữ liệu bổ sung (JSONB)';
        COMMENT ON COLUMN snapshots.created_at IS 'Ngày tạo snapshot';
        COMMENT ON COLUMN snapshots.created_by IS 'Người hoặc hệ thống tạo snapshot';
    END IF;
END
$$;
