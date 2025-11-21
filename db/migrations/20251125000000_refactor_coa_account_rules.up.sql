-- Drop old tables if they exist (sẽ migrate data sau)
-- Tạo bảng mới với cấu trúc rõ ràng hơn

-- 1. Bảng lưu các TYPE (ASSET, LIAB, REV, EXP)
CREATE TABLE IF NOT EXISTS coa_account_rule_types (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(16) NOT NULL UNIQUE,  -- ASSET, LIAB, REV, EXP
    name VARCHAR(128) NOT NULL,       -- Asset, Liability, Revenue, Expense
    separator VARCHAR(8) DEFAULT ':',  -- Separator sau type (thường là :)
    status VARCHAR(16) DEFAULT 'ACTIVE',
    sort_order INT DEFAULT 0,
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 2. Bảng lưu các GROUP (FLOAT, BANK, CUSTODY, CLEARING, DETAILS, KIND)
-- GROUP có thể NULL nếu TYPE không có group (như REV, EXP có thể bỏ qua KIND)
CREATE TABLE IF NOT EXISTS coa_account_rule_groups (
    id BIGSERIAL PRIMARY KEY,
    type_id BIGINT NOT NULL REFERENCES coa_account_rule_types(id) ON DELETE CASCADE,
    code VARCHAR(64) NOT NULL,         -- FLOAT, BANK, CUSTODY, CLEARING, DETAILS, KIND
    name VARCHAR(128) NOT NULL,        -- Float, Bank, Custody, Clearing, Details, Kind
    separator VARCHAR(8) DEFAULT ':',  -- Separator sau group (thường là :)
    input_type VARCHAR(16) DEFAULT 'SELECT',  -- SELECT hoặc TEXT
    status VARCHAR(16) DEFAULT 'ACTIVE',
    sort_order INT DEFAULT 0,
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(type_id, code)  -- Mỗi type chỉ có 1 group với code này
);

-- 3. Bảng lưu các STEPS trong GROUP hoặc TYPE
-- Nếu group_id NULL thì step thuộc về TYPE (như REV, EXP không có group)
CREATE TABLE IF NOT EXISTS coa_account_rule_steps (
    id BIGSERIAL PRIMARY KEY,
    type_id BIGINT NOT NULL REFERENCES coa_account_rule_types(id) ON DELETE CASCADE,
    group_id BIGINT REFERENCES coa_account_rule_groups(id) ON DELETE CASCADE,
    step_order INT NOT NULL,           -- Thứ tự step (1, 2, 3...)
    type VARCHAR(16) NOT NULL,         -- SELECT hoặc TEXT
    label VARCHAR(128),                 -- Label hiển thị (có thể NULL)
    
    -- Nếu type = SELECT: dùng category_id
    category_id INT REFERENCES rule_categories(id) ON DELETE SET NULL,
    category_code VARCHAR(50),         -- Cache category code để query nhanh
    
    -- Nếu type = TEXT: dùng input_code
    input_code VARCHAR(64),             -- DETAILS, etc.
    
    separator VARCHAR(8) DEFAULT '.',  -- Separator sau step này
    status VARCHAR(16) DEFAULT 'ACTIVE',
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    -- Constraint: phải có category_id HOẶC input_code
    CONSTRAINT chk_step_target CHECK (
        (type = 'SELECT' AND category_id IS NOT NULL AND input_code IS NULL) 
        OR (type = 'TEXT' AND category_id IS NULL AND input_code IS NOT NULL)
    )
);

-- Unique constraint cho steps: mỗi group (hoặc type nếu group_id NULL) chỉ có 1 step với step_order này
-- Dùng partial unique index vì không thể dùng COALESCE trong UNIQUE constraint
CREATE UNIQUE INDEX IF NOT EXISTS idx_coa_rule_steps_unique_group 
    ON coa_account_rule_steps(type_id, group_id, step_order) 
    WHERE group_id IS NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_coa_rule_steps_unique_type 
    ON coa_account_rule_steps(type_id, step_order) 
    WHERE group_id IS NULL;

-- Indexes
CREATE INDEX IF NOT EXISTS idx_coa_rule_groups_type_id ON coa_account_rule_groups(type_id);
CREATE INDEX IF NOT EXISTS idx_coa_rule_steps_type_id ON coa_account_rule_steps(type_id);
CREATE INDEX IF NOT EXISTS idx_coa_rule_steps_group_id ON coa_account_rule_steps(group_id);
CREATE INDEX IF NOT EXISTS idx_coa_rule_steps_category_id ON coa_account_rule_steps(category_id);

-- Comments
COMMENT ON TABLE coa_account_rule_types IS 'Các loại tài khoản COA: ASSET, LIAB, REV, EXP';
COMMENT ON TABLE coa_account_rule_groups IS 'Các nhóm trong mỗi type: FLOAT, BANK, CUSTODY, etc. Có thể NULL nếu type không có group';
COMMENT ON TABLE coa_account_rule_steps IS 'Các bước trong group hoặc type để tạo mã code. Mỗi step có thể là SELECT (từ rule_categories) hoặc TEXT input';

