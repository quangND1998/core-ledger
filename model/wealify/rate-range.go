package model

import (
	"time"
)

const TableNameRateRange = "rate-ranges"

// RateRange mapped from table <rate-ranges>
type RateRange struct {
	CreatedAt time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(6)" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP(6)" json:"updated_at"`
	Status    int32     `gorm:"column:status;not null;default:1" json:"status"`
	IsDeleted int32     `gorm:"column:is_deleted;not null" json:"is_deleted"`
	ID        string    `gorm:"column:id;primaryKey" json:"id"`
	Min       float64   `gorm:"column:min;not null" json:"min"`
	Max       float64   `gorm:"column:max;not null;default:999999999999999" json:"max"`
	Value     float64   `gorm:"column:value;not null" json:"value"`
	Type      string    `gorm:"column:type;not null;default:MANUAL" json:"type"`

	RateID string `gorm:"column:rate_id" json:"rate_id"`

	Rate *Rate `gorm:"foreignKey:RateID;references:ID" json:"rate"`
}

// TableName RateRange's table name
func (*RateRange) TableName() string {
	return TableNameRateRange
}
