package model

import (
	"time"

	"gorm.io/datatypes"
)

// CoaAccountRuleStep đại diện cho STEP trong quy tắc tạo mã COA
// Mỗi step có thể là SELECT (từ rule_categories) hoặc TEXT input
// Nếu group_id NULL thì step thuộc về TYPE (như REV, EXP không có group)
type CoaAccountRuleStep struct {
	ID           uint64            `gorm:"primaryKey;autoIncrement" json:"id"`
	TypeID       uint64            `gorm:"not null" json:"type_id"`
	GroupID      *uint64           `gorm:"index" json:"group_id,omitempty"` // NULL nếu step thuộc về TYPE
	StepOrder    int               `gorm:"not null" json:"step_order"`      // Thứ tự step (1, 2, 3...)
	Type         string            `gorm:"type:varchar(16);not null" json:"type"` // SELECT hoặc TEXT
	Label        *string           `gorm:"type:varchar(128)" json:"label,omitempty"`

	// Nếu type = SELECT: dùng category_id
	CategoryID   *uint64           `gorm:"index" json:"category_id,omitempty"`
	CategoryCode *string           `gorm:"type:varchar(50)" json:"category_code,omitempty"` // Cache để query nhanh

	// Nếu type = TEXT: dùng input_code
	InputCode    *string           `gorm:"type:varchar(64)" json:"input_code,omitempty"` // DETAILS, etc.

	Separator    string            `gorm:"type:varchar(8);default:'.'" json:"separator"` // Separator sau step này
	Status       string            `gorm:"type:varchar(16);default:'ACTIVE'" json:"status"`
	Metadata     datatypes.JSONMap `gorm:"type:jsonb;default:'{}'" json:"metadata,omitempty"`
	CreatedAt    time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time         `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	RuleType CoaAccountRuleType  `gorm:"foreignKey:TypeID" json:"rule_type,omitempty"`
	Group    *CoaAccountRuleGroup `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	Category *RuleCategory        `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

func (CoaAccountRuleStep) TableName() string {
	return "coa_account_rule_steps"
}

