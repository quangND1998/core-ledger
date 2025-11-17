CREATE TABLE IF NOT EXISTS account_rule_layers (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(64) NOT NULL UNIQUE,
    name VARCHAR(128) NOT NULL,
    layer_index INT NOT NULL,
    status VARCHAR(16) NOT NULL DEFAULT 'ACTIVE',
    description TEXT,
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS account_rule_options (
    id BIGSERIAL PRIMARY KEY,
    layer_id BIGINT NOT NULL REFERENCES account_rule_layers(id) ON DELETE CASCADE,
    parent_option_id BIGINT REFERENCES account_rule_options(id) ON DELETE CASCADE,
    code VARCHAR(128) NOT NULL,
    name VARCHAR(255) NOT NULL,
    status VARCHAR(16) NOT NULL DEFAULT 'ACTIVE',
    sort_order INT NOT NULL DEFAULT 0,
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_account_rule_options_unique
    ON account_rule_options (layer_id, COALESCE(parent_option_id, 0), code);

CREATE TABLE IF NOT EXISTS account_rule_option_steps (
    id BIGSERIAL PRIMARY KEY,
    option_id BIGINT NOT NULL REFERENCES account_rule_options(id) ON DELETE CASCADE,
    step_order INT NOT NULL,
    category_id BIGINT REFERENCES rule_categories(id) ON DELETE SET NULL,
    input_code VARCHAR(64),
    input_label VARCHAR(128),
    input_type VARCHAR(16) NOT NULL DEFAULT 'SELECT',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(option_id, step_order),
    CONSTRAINT chk_account_rule_option_steps_target CHECK (
        (category_id IS NOT NULL AND input_code IS NULL)
        OR (category_id IS NULL AND input_code IS NOT NULL)
    )
);

CREATE INDEX IF NOT EXISTS idx_account_rule_options_layer_id 
    ON account_rule_options(layer_id);

CREATE INDEX IF NOT EXISTS idx_account_rule_options_parent_id 
    ON account_rule_options(parent_option_id);

CREATE INDEX IF NOT EXISTS idx_account_rule_option_steps_option_id 
    ON account_rule_option_steps(option_id);

CREATE INDEX IF NOT EXISTS idx_account_rule_option_steps_category_id 
    ON account_rule_option_steps(category_id);
