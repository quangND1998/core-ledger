DO $$
BEGIN
    IF NOT EXISTS (
        SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'blacklisted_tokens'
    ) THEN
        CREATE TABLE blacklisted_tokens (
            id BIGSERIAL PRIMARY KEY,
            token_hash VARCHAR(64) NOT NULL UNIQUE,
            user_id BIGINT NOT NULL,
            expires_at TIMESTAMP NOT NULL,
            created_at TIMESTAMP DEFAULT NOW(),
            
            CONSTRAINT fk_blacklisted_tokens_user 
                FOREIGN KEY (user_id) 
                REFERENCES users(id) 
                ON DELETE CASCADE
        );

        CREATE INDEX idx_blacklisted_tokens_token_hash ON blacklisted_tokens(token_hash);
        CREATE INDEX idx_blacklisted_tokens_user_id ON blacklisted_tokens(user_id);
        CREATE INDEX idx_blacklisted_tokens_expires_at ON blacklisted_tokens(expires_at);

        COMMENT ON TABLE blacklisted_tokens IS 'Blacklisted JWT tokens - tokens that have been revoked';
        COMMENT ON COLUMN blacklisted_tokens.token_hash IS 'SHA256 hash of the JWT token';
        COMMENT ON COLUMN blacklisted_tokens.expires_at IS 'Token expiration time - used for cleanup';
    END IF;
END
$$;


