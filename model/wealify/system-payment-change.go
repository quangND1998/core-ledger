package model

import (
	"core-ledger/model/enum"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const TableNameSystemPaymentChange = "system-payment-changes"

type SystemPaymentChange struct {
	CreatedAt       time.Time            `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(6)" json:"created_at"`
	UpdatedAt       time.Time            `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP(6)" json:"updated_at"`
	Status          int32                `gorm:"column:status;not null;default:1" json:"status"`
	IsDeleted       int32                `gorm:"column:is_deleted;not null" json:"is_deleted"`
	ID              string               `gorm:"column:id;primaryKey" json:"id"`
	ChangedAmount   float64              `gorm:"column:changed_amount;not null" json:"changed_amount"`
	ChangeType      enum.TransactionType `gorm:"column:change_type;not null" json:"change_type"`
	SystemPaymentID string               `gorm:"column:system_payment_id" json:"system_payment_id"`
	TransactionID   string               `gorm:"column:transaction_id" json:"transaction_id"`

	SystemPayment *SystemPayment `gorm:"foreignKey:SystemPaymentID;references:ID" json:"system_payment"`
}

// TableName SystemPaymentChange's table name
func (*SystemPaymentChange) TableName() string {
	return TableNameSystemPaymentChange
}

func (t *SystemPaymentChange) BeforeSave(tx *gorm.DB) (err error) {
	if t.ID == "" {
		t.ID = uuid.NewString()
	}
	return
}

func (t *SystemPaymentChange) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == "" {
		t.ID = uuid.NewString()
	}
	return
}

func (t *SystemPaymentChange) BeforeUpdate(tx *gorm.DB) (err error) {
	t.UpdatedAt = time.Now()
	return
}
