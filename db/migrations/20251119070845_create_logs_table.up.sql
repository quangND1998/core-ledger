CREATE TABLE IF NOT EXISTS logs (
    id BIGSERIAL PRIMARY KEY,
    loggable_id BIGINT NOT NULL,
    loggable_type VARCHAR(128) NOT NULL,

    action VARCHAR(64) NOT NULL,
    old_value JSONB,
    new_value JSONB,
    metadata JSONB DEFAULT '{}'::jsonb,

    created_by BIGINT, -- optional: user id who triggered
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Index để truy vấn theo object nhanh (type + id)
CREATE INDEX IF NOT EXISTS idx_logs_loggable ON logs (loggable_type, loggable_id);
-- Index trên created_at nếu cần truy vấn theo thời gian
CREATE INDEX IF NOT EXISTS idx_logs_created_at ON logs (created_at);