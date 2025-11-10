DO $$
BEGIN
    IF NOT EXISTS (
        SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'transaction_logs'
    ) THEN
        CREATE TABLE transaction_logs (
            id BIGSERIAL PRIMARY KEY,
            aggregate_type VARCHAR(32) NOT NULL,   -- loại thực thể: JOURNAL / SNAPSHOT
            aggregate_id BIGINT NOT NULL,          -- ID nguồn (journal_id/snapshot_id)
            event_type VARCHAR(64) NOT NULL,       -- loại sự kiện: ledger.posted, ledger.snapshot_locked...
            event_key VARCHAR(191) NOT NULL,       -- khóa idempotent để consumer dedupe
            partition_key VARCHAR(191) NOT NULL,   -- khóa phân vùng (vd: cùng user/ccy)
            payload JSON NOT NULL,                  -- nội dung sự kiện
            headers JSON,                           -- metadata / trace info
            status VARCHAR(16) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING','PUBLISHED','FAILED','DEAD')),
            attempts INT NOT NULL DEFAULT 0,       -- số lần gửi
            next_attempt_at TIMESTAMP,             -- lịch gửi lại
            last_attempt_at TIMESTAMP,             -- lần thử gần nhất
            published_at TIMESTAMP,                -- thời điểm gửi thành công
            error_last TEXT,                        -- lỗi gần nhất nếu có
            tenant_id VARCHAR(36) NOT NULL,
            ledger_code VARCHAR(32),
            seq BIGINT,                             -- số tăng đơn điệu để ordering
            created_at TIMESTAMP NOT NULL DEFAULT NOW()
        );

        -- Indexes
        CREATE INDEX idx_transaction_logs_status ON transaction_logs(status);
        CREATE INDEX idx_transaction_logs_tenant ON transaction_logs(tenant_id);
        CREATE INDEX idx_transaction_logs_ledger ON transaction_logs(ledger_code);
        CREATE INDEX idx_transaction_logs_partition_key ON transaction_logs(partition_key);
        CREATE INDEX idx_transaction_logs_aggregate ON transaction_logs(aggregate_type, aggregate_id);
        CREATE INDEX idx_transaction_logs_next_attempt ON transaction_logs(next_attempt_at);
        CREATE INDEX idx_transaction_logs_seq ON transaction_logs(seq);

        -- Comments
        COMMENT ON TABLE transaction_logs IS 'Bảng lưu trữ sự kiện ledger: journal, snapshot, transaction…';
        COMMENT ON COLUMN transaction_logs.id IS 'Khóa chính';
        COMMENT ON COLUMN transaction_logs.aggregate_type IS 'Loại thực thể nguồn (JOURNAL/SNAPSHOT…)';
        COMMENT ON COLUMN transaction_logs.aggregate_id IS 'ID nguồn (journal_id/snapshot_id)';
        COMMENT ON COLUMN transaction_logs.event_type IS 'Loại sự kiện';
        COMMENT ON COLUMN transaction_logs.event_key IS 'Khóa idempotent để tránh duplicate';
        COMMENT ON COLUMN transaction_logs.partition_key IS 'Khóa phân vùng';
        COMMENT ON COLUMN transaction_logs.payload IS 'Nội dung sự kiện JSON';
        COMMENT ON COLUMN transaction_logs.headers IS 'Metadata/trace info JSON';
        COMMENT ON COLUMN transaction_logs.status IS 'Trạng thái gửi sự kiện';
        COMMENT ON COLUMN transaction_logs.attempts IS 'Số lần gửi';
        COMMENT ON COLUMN transaction_logs.next_attempt_at IS 'Thời điểm dự kiến gửi lại';
        COMMENT ON COLUMN transaction_logs.last_attempt_at IS 'Lần thử gần nhất';
        COMMENT ON COLUMN transaction_logs.published_at IS 'Thời điểm gửi thành công';
        COMMENT ON COLUMN transaction_logs.error_last IS 'Lỗi gần nhất';
        COMMENT ON COLUMN transaction_logs.tenant_id IS 'Tenant nếu đa tenant';
        COMMENT ON COLUMN transaction_logs.ledger_code IS 'Mã sổ';
        COMMENT ON COLUMN transaction_logs.seq IS 'Số tăng đơn điệu để ordering';
        COMMENT ON COLUMN transaction_logs.created_at IS 'Ngày tạo sự kiện';
    END IF;
END
$$;
