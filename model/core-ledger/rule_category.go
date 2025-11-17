package model

import (
	"time"
)

type RuleCategory struct {
	ID         uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string      `gorm:"size:100;unique;not null" json:"name"`
	Code       string      `gorm:"size:50;unique;not null" json:"code"`
	CreatedAt  time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
	RuleValues []RuleValue `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE" json:"rule_values"` // quan hệ 1-n với RuleValue
}

func (c *RuleCategory) TableName() string {
	return "rule_categories"
}
