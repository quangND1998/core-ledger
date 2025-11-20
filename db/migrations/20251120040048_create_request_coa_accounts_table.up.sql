DO $$
BEGIN
    IF NOT EXISTS (
        SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'request_coa_accounts'
    ) THEN
        CREATE TABLE request_coa_accounts (
            id BIGSERIAL PRIMARY KEY,
            coa_account_id BIGINT NULL,
            request_type VARCHAR(16) NOT NULL CHECK (request_type IN ('CREATE','EDIT')),
            request_status VARCHAR(16) NOT NULL DEFAULT 'PENDING' CHECK (request_status IN ('PENDING','APPROVED','REJECTED')),
            maker_id BIGINT NOT NULL,
            checker_id BIGINT NULL,
            data JSONB NOT NULL,
            comment TEXT NULL,
            reject_reason TEXT NULL,
            created_at TIMESTAMP DEFAULT NOW(),
            updated_at TIMESTAMP DEFAULT NOW(),
            checked_at TIMESTAMP NULL,
            
            CONSTRAINT fk_request_coa_accounts_coa_account 
                FOREIGN KEY (coa_account_id) 
                REFERENCES coa_accounts(id) 
                ON DELETE SET NULL
                ON UPDATE CASCADE,
            CONSTRAINT fk_request_coa_accounts_maker 
                FOREIGN KEY (maker_id) 
                REFERENCES users(id) 
                ON DELETE RESTRICT
                ON UPDATE CASCADE,
            CONSTRAINT fk_request_coa_accounts_checker 
                FOREIGN KEY (checker_id) 
                REFERENCES users(id) 
                ON DELETE SET NULL
                ON UPDATE CASCADE
        );

        CREATE INDEX idx_request_coa_accounts_coa_account_id ON request_coa_accounts(coa_account_id);
        CREATE INDEX idx_request_coa_accounts_maker_id ON request_coa_accounts(maker_id);
        CREATE INDEX idx_request_coa_accounts_checker_id ON request_coa_accounts(checker_id);
        CREATE INDEX idx_request_coa_accounts_request_type ON request_coa_accounts(request_type);
        CREATE INDEX idx_request_coa_accounts_request_status ON request_coa_accounts(request_status);
        CREATE INDEX idx_request_coa_accounts_created_at ON request_coa_accounts(created_at);

        COMMENT ON TABLE request_coa_accounts IS 'Requests for creating or editing COA accounts - requires approval workflow';
        COMMENT ON COLUMN request_coa_accounts.coa_account_id IS 'Reference to coa_accounts (NULL for CREATE requests)';
        COMMENT ON COLUMN request_coa_accounts.request_type IS 'Type of request: CREATE or EDIT';
        COMMENT ON COLUMN request_coa_accounts.request_status IS 'Status: PENDING, APPROVED, or REJECTED';
        COMMENT ON COLUMN request_coa_accounts.maker_id IS 'User who created the request';
        COMMENT ON COLUMN request_coa_accounts.checker_id IS 'User who approved/rejected the request';
        COMMENT ON COLUMN request_coa_accounts.data IS 'JSON data containing all COA account fields';
        COMMENT ON COLUMN request_coa_accounts.comment IS 'Comment from checker';
        COMMENT ON COLUMN request_coa_accounts.reject_reason IS 'Reason for rejection';
        COMMENT ON COLUMN request_coa_accounts.checked_at IS 'Timestamp when request was approved/rejected';
    END IF;
END
$$;

