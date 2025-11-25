package model

import (
	"time"

	"gorm.io/datatypes"
)

// CoaAccountRuleGroup đại diện cho GROUP trong quy tắc tạo mã COA
// Ví dụ: FLOAT, BANK, CUSTODY, CLEARING, DETAILS, KIND
// GROUP có thể NULL nếu TYPE không có group (như REV, EXP có thể bỏ qua KIND)
type CoaAccountRuleGroup struct {
	ID        uint64            `gorm:"primaryKey;autoIncrement" json:"id"`
	TypeID    uint64            `gorm:"not null" json:"type_id"`
	Code      string            `gorm:"type:varchar(64);not null" json:"code"`              // FLOAT, BANK, CUSTODY, etc.
	Name      string            `gorm:"type:varchar(128);not null" json:"name"`            // Float, Bank, Custody, etc.
	Separator string            `gorm:"type:varchar(8);default:':'" json:"separator"`       // Separator sau group
	InputType string            `gorm:"type:varchar(16);default:'SELECT'" json:"input_type"` // SELECT hoặc TEXT
	Status    string            `gorm:"type:varchar(16);default:'ACTIVE'" json:"status"`
	SortOrder int               `gorm:"default:0" json:"sort_order"`
	Metadata  datatypes.JSONMap `gorm:"type:jsonb;default:'{}'" json:"metadata,omitempty"`
	CreatedAt time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time         `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	Type  CoaAccountRuleType  `gorm:"foreignKey:TypeID" json:"type,omitempty"`
	Steps []CoaAccountRuleStep `gorm:"foreignKey:GroupID;constraint:OnDelete:CASCADE" json:"steps,omitempty"`
}

func (CoaAccountRuleGroup) TableName() string {
	return "coa_account_rule_groups"
}



