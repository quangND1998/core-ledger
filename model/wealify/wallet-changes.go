package model

import (
	"core-ledger/model/enum"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const TableNameWalletChange = "wallet-changes"

// WalletChange mapped from table <wallet-changes>
type WalletChange struct {
	CreatedAt     time.Time             `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(6)" json:"created_at"`
	UpdatedAt     time.Time             `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP(6)" json:"updated_at"`
	Status        bool                  `gorm:"column:status;not null;default:1" json:"status"`
	IsDeleted     bool                  `gorm:"column:is_deleted;not null" json:"is_deleted"`
	ID            string                `gorm:"column:id;primaryKey" json:"id"`
	ChangedAmount float64               `gorm:"column:changed_amount;not null" json:"changed_amount"`
	ChangeType    enum.WalletChangeType `gorm:"column:change_type;not null" json:"change_type"`

	WalletID      string `gorm:"column:wallet_id" json:"wallet_id"`
	TransactionID string `gorm:"column:transaction_id" json:"transaction_id"`

	Transaction *Transaction `gorm:"foreignkey:TransactionID" json:"transaction"`
}

// TableName WalletChange's table name
func (*WalletChange) TableName() string {
	return TableNameWalletChange
}

func (w *WalletChange) BeforeCreate(tx *gorm.DB) error {
	if w.ID == "" {
		w.ID = uuid.NewString()
	}
	return nil
}
