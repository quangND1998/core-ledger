DO $$
BEGIN
    IF NOT EXISTS (
        SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'entries'
    ) THEN

        CREATE TABLE entries (
            id BIGSERIAL PRIMARY KEY,

            journal_id BIGINT NOT NULL,
            line_no INT NOT NULL,
            account_id BIGINT NOT NULL,
            dc VARCHAR(1) NOT NULL CHECK (dc IN ('D','C')),
            amount DECIMAL(28,8) NOT NULL CHECK (amount >= 0),
            amount_atoms BIGINT,
            meta JSON,
            memo VARCHAR(256),

            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

            tenant_id VARCHAR(36),
            ledger_code VARCHAR(32),
            batch_id VARCHAR(36),

            UNIQUE (journal_id, line_no),

            CONSTRAINT fk_entries_journal
                FOREIGN KEY (journal_id)
                REFERENCES journals(id)
                ON DELETE CASCADE
                ON UPDATE CASCADE,

            CONSTRAINT fk_entries_account
                FOREIGN KEY (account_id)
                REFERENCES coa_accounts(id)
                ON DELETE RESTRICT
                ON UPDATE CASCADE
        );

        -- Indexes
        CREATE INDEX idx_entries_journal_id ON entries(journal_id);
        CREATE INDEX idx_entries_account_id ON entries(account_id);
        CREATE INDEX idx_entries_tenant_id ON entries(tenant_id);
        CREATE INDEX idx_entries_ledger_code ON entries(ledger_code);
        CREATE INDEX idx_entries_batch_id ON entries(batch_id);
        CREATE INDEX idx_entries_dc ON entries(dc);

        COMMENT ON TABLE entries IS 'Các dòng chi tiết của journal entry, ghi nợ hoặc ghi có trên từng tài khoản';
        COMMENT ON COLUMN entries.dc IS 'Dòng ghi nợ hay ghi có: D = Debit, C = Credit';
    END IF;
END
$$;
