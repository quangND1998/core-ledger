DO $$
BEGIN
    IF NOT EXISTS (
        SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'journals'
    ) THEN
        CREATE TABLE journals (
            id BIGSERIAL PRIMARY KEY,
            ts TIMESTAMP DEFAULT NOW() NOT NULL,
            status VARCHAR(16) NOT NULL CHECK (status IN ('DRAFT','POSTED','REVERSED')),
            idempotency_key VARCHAR(191) NOT NULL UNIQUE,
            currency CHAR(8) NOT NULL,
            source VARCHAR(64) NOT NULL,
            memo VARCHAR(256),
            meta JSONB,
            reversal_of BIGINT REFERENCES journals(id) ON DELETE SET NULL ON UPDATE CASCADE,
            posted_by VARCHAR(64),
            posted_at TIMESTAMP,
            lock_version INT DEFAULT 0 NOT NULL,
            created_at TIMESTAMP DEFAULT NOW() NOT NULL,
            updated_at TIMESTAMP DEFAULT NOW() NOT NULL,
            tenant_id VARCHAR(36),
            ledger_code VARCHAR(32),
            batch_id VARCHAR(36)
        );

        CREATE INDEX idx_journals_status ON journals(status);
        CREATE INDEX idx_journals_source ON journals(source);
        CREATE INDEX idx_journals_tenant ON journals(tenant_id);
        CREATE INDEX idx_journals_batch ON journals(batch_id);

        -- Thêm mô tả cho bảng
        COMMENT ON TABLE journals IS 'Bảng lưu trữ các journal entry (giao dịch, sổ kế toán)';

        -- Thêm mô tả cho các cột
        COMMENT ON COLUMN journals.id IS 'ID tự tăng của giao dịch';
        COMMENT ON COLUMN journals.ts IS 'Thời điểm giao dịch được ghi nhận';
        COMMENT ON COLUMN journals.status IS 'Trạng thái giao dịch: DRAFT, POSTED, REVERSED';
        COMMENT ON COLUMN journals.idempotency_key IS 'Khóa idempotency để tránh tạo trùng lặp giao dịch';
        COMMENT ON COLUMN journals.currency IS 'Mã tiền tệ (VD: USD, VND)';
        COMMENT ON COLUMN journals.source IS 'Nguồn phát sinh giao dịch (VD: hệ thống, manual, API)';
        COMMENT ON COLUMN journals.memo IS 'Ghi chú cho giao dịch';
        COMMENT ON COLUMN journals.meta IS 'Thông tin bổ sung dạng JSON';
        COMMENT ON COLUMN journals.reversal_of IS 'Nếu là giao dịch đảo ngược, tham chiếu đến ID của giao dịch gốc';
        COMMENT ON COLUMN journals.posted_by IS 'Người thực hiện ghi sổ';
        COMMENT ON COLUMN journals.posted_at IS 'Thời gian ghi sổ';
        COMMENT ON COLUMN journals.lock_version IS 'Dùng cho optimistic locking';
        COMMENT ON COLUMN journals.created_at IS 'Thời gian tạo bản ghi';
        COMMENT ON COLUMN journals.updated_at IS 'Thời gian cập nhật bản ghi';
        COMMENT ON COLUMN journals.tenant_id IS 'Hỗ trợ đa tenant';
        COMMENT ON COLUMN journals.ledger_code IS 'Mã sổ cái liên quan';
        COMMENT ON COLUMN journals.batch_id IS 'Nhóm giao dịch theo batch';
    END IF;
END
$$;