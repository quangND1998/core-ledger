package model

import (
	"time"

	"gorm.io/datatypes"
)

type Log struct {
	ID uint64 `gorm:"primaryKey"`

	LoggableID   uint64 `gorm:"not null;index:idx_logs_loggable,priority:2"`
	LoggableType string `gorm:"size:128;not null;index:idx_logs_loggable,priority:1"`

	Action   string         `gorm:"size:64;not null"`
	OldValue datatypes.JSON `gorm:"type:jsonb"`
	NewValue datatypes.JSON `gorm:"type:jsonb"`
	Metadata datatypes.JSON `gorm:"type:jsonb;default:'{}'::jsonb"`

	CreatedBy *uint64   `json:"created_by,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (Log) TableName() string {
	return "logs"
}