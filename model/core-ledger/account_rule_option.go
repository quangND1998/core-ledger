package model

import (
	"time"

	"gorm.io/datatypes"
)

type AccountRuleOption struct {
	ID             uint64            `gorm:"primaryKey;autoIncrement" json:"id"`
	LayerID        uint64            `gorm:"not null" json:"layer_id"`
	ParentOptionID *uint64           `json:"parent_option_id,omitempty"`
	Code           string            `gorm:"type:varchar(128);not null" json:"code"`
	Name           string            `gorm:"type:varchar(255);not null" json:"name"`
	Status         string            `gorm:"type:varchar(16);default:'ACTIVE'" json:"status"`
	SortOrder      int               `gorm:"default:0" json:"sort_order"`
	InputType      string            `gorm:"type:varchar(16);default:'SELECT'" json:"input_type"`
	Metadata       datatypes.JSONMap `gorm:"type:jsonb;default:'{}'" json:"metadata,omitempty"`
	CreatedAt      time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time         `gorm:"autoUpdateTime" json:"updated_at"`
}

func (AccountRuleOption) TableName() string {
	return "account_rule_options"
}

