package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

const TableNameSystemPaymentBalance = "system_payment_balances"

type SystemPaymentBalance struct {
	CreatedAt       time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(6)" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP(6)" json:"updated_at"`
	CurrencySymbol  string     `gorm:"column:currency_symbol;not null;default:VND" json:"currency_symbol"`
	DeletedAt       *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	ID              string     `gorm:"column:id;primaryKey" json:"id"`
	Balance         float64    `gorm:"column:balance;not null" json:"balance"`
	SystemPaymentID string     `gorm:"column:system_payment_id" json:"system_payment_id"`
}

func (*SystemPaymentBalance) TableName() string {
	return TableNameSystemPaymentBalance
}

func (s *SystemPaymentBalance) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.NewString()
	}
	return nil
}

func (s *SystemPaymentBalance) BeforeSave(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.NewString()
	}
	return nil
}
