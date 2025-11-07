package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"

	"gorm.io/datatypes"
)

type SystemPaymentTransaction struct {
	ID             string         `gorm:"primaryKey;column:id"`
	Code           string         `gorm:"column:code;uniqueIndex:uq_system_payment_transactions_code_system_payment_id"`
	Amount         float64        `gorm:"column:amount;type:decimal(20,8);not null"`
	CurrencySymbol string         `gorm:"column:currency_symbol;size:20;not null"`
	Direction      string         `gorm:"column:direction;size:10"` // "IN" / "OUT"
	ProviderType   string         `gorm:"column:provider_type;size:100"`
	ProviderStatus string         `gorm:"column:provider_status;size:100"`
	Status         string         `gorm:"column:status;size:50"`
	Metadata       datatypes.JSON `gorm:"column:metadata;type:jsonb"`
	CreatedAt      time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;autoUpdateTime"`

	SystemPaymentID string `gorm:"column:system_payment_id;not null;uniqueIndex:uq_system_payment_transactions_code_system_payment_id"`
}

func (*SystemPaymentTransaction) TableName() string {
	return "system_payment_transactions"
}

func (item *SystemPaymentTransaction) BeforeCreate(tx *gorm.DB) error {
	if item.ID == "" {
		item.ID = uuid.NewString()
	}
	return nil
}

func (item *SystemPaymentTransaction) BeforeSave(tx *gorm.DB) error {
	if item.ID == "" {
		item.ID = uuid.NewString()
	}
	return nil
}
