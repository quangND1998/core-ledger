package model

import (
	"time"
)

type RuleValue struct {
	ID         uint64        `gorm:"primaryKey;autoIncrement" json:"id"`
	CategoryID uint          `gorm:"not null" json:"category_id"` // FK
	Name       string        `gorm:"size:255" json:"name"`
	Value      string        `gorm:"size:255;not null" json:"value"`
	SortOrder  int           `gorm:"default:0" json:"sort_order"`
	IsDelete   bool          `gorm:"default:false" json:"is_delete"`
	CreatedAt  time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	Category   *RuleCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"` // quan hệ ngược
	Logs       []Log         `gorm:"polymorphic:Loggable;polymorphicValue:rule_values" json:"logs,omitempty"`
}

func (c *RuleValue) TableName() string {
	return "rule_values"
}
func (c *RuleValue) GetLoggableID() uint64   { return c.ID }
func (c *RuleValue) GetLoggableType() string { return "rule_values" }
