package model

import (
	"core-ledger/model/enum"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const TableNameTransactionLedger = "transaction_ledgers"

type TransactionLedger struct {
	ID            string         `gorm:"primaryKey"`
	TransactionID string         `gorm:"not null;index"`
	WalletID      string         `gorm:"not null;index"`
	EntryType     enum.EntryType `gorm:"type:varchar(16);not null"` // DEBIT or CREDIT
	Amount        float64        `gorm:"type:decimal(20,8);not null"`
	CreatedAt     time.Time      `gorm:"autoCreateTime;default:CURRENT_TIMESTAMP"`
}

func (t *TransactionLedger) TableName() string {
	return TableNameTransactionLedger
}

func (t *TransactionLedger) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.NewString()
	}
	return nil
}
