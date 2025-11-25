package model

import (
	"time"

	"gorm.io/datatypes"
)

// CoaAccountRuleType đại diện cho TYPE trong quy tắc tạo mã COA
// Ví dụ: ASSET, LIAB, REV, EXP
type CoaAccountRuleType struct {
	ID        uint64            `gorm:"primaryKey;autoIncrement" json:"id"`
	Code      string            `gorm:"type:varchar(16);unique;not null" json:"code"`      // ASSET, LIAB, REV, EXP
	Name      string            `gorm:"type:varchar(128);not null" json:"name"`            // Asset, Liability, Revenue, Expense
	Separator string            `gorm:"type:varchar(8);default:':'" json:"separator"`      // Separator sau type
	Status    string            `gorm:"type:varchar(16);default:'ACTIVE'" json:"status"`
	SortOrder int               `gorm:"default:0" json:"sort_order"`
	Metadata  datatypes.JSONMap `gorm:"type:jsonb;default:'{}'" json:"metadata,omitempty"`
	CreatedAt time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time         `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	Groups []CoaAccountRuleGroup `gorm:"foreignKey:TypeID;constraint:OnDelete:CASCADE" json:"groups,omitempty"`
	Steps  []CoaAccountRuleStep `gorm:"foreignKey:TypeID;constraint:OnDelete:CASCADE" json:"steps,omitempty"`
}

func (CoaAccountRuleType) TableName() string {
	return "coa_account_rule_types"
}



