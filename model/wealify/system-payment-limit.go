package model

import (
	"time"
)

const TableNameSystemPaymentLimit = "system-payment-limits"

// SystemPaymentLimit mapped from table <system-payment-limits>
type SystemPaymentLimit struct {
	CreatedAt       time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(6)" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP(6)" json:"updated_at"`
	Status          int32     `gorm:"column:status;not null;default:1" json:"status"`
	IsDeleted       int32     `gorm:"column:is_deleted;not null" json:"is_deleted"`
	ID              string    `gorm:"column:id;primaryKey" json:"id"`
	MaxPerDay       float64   `gorm:"column:max_per_day;not null" json:"max_per_day"`
	MaxPerMonth     float64   `gorm:"column:max_per_month" json:"max_per_month"`
	SystemPaymentID string    `gorm:"column:system_payment_id" json:"system_payment_id"`
	CurrencySymbol  string    `gorm:"column:currency_symbol" json:"currency_symbol"`
}

// TableName SystemPaymentLimit's table name
func (*SystemPaymentLimit) TableName() string {
	return TableNameSystemPaymentLimit
}
