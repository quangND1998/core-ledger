package model

import (
	"time"

	"gorm.io/datatypes"
)

type AccountRuleLayer struct {
	ID          uint64            `gorm:"primaryKey;autoIncrement" json:"id"`
	Code        string            `gorm:"type:varchar(64);unique;not null" json:"code"`
	Name        string            `gorm:"type:varchar(128);not null" json:"name"`
	LayerIndex  int               `gorm:"not null" json:"layer_index"`
	Status      string            `gorm:"type:varchar(16);default:'ACTIVE'" json:"status"`
	Description string            `gorm:"type:text" json:"description,omitempty"`
	Metadata    datatypes.JSONMap `gorm:"type:jsonb;default:'{}'" json:"metadata,omitempty"`
	CreatedAt   time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time         `gorm:"autoUpdateTime" json:"updated_at"`
}

func (AccountRuleLayer) TableName() string {
	return "account_rule_layers"
}


