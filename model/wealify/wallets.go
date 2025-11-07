package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

const TableNameWallet = "wallets"

// Wallet mapped from table <wallets>
type Wallet struct {
	CreatedAt time.Time      `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(6)" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP(6)" json:"updated_at"`
	Status    int32          `gorm:"column:status;not null;default:1" json:"status"`
	IsDeleted int32          `gorm:"column:is_deleted;not null" json:"is_deleted"`
	ID        string         `gorm:"column:id;primaryKey" json:"id"`
	Balance   float64        `gorm:"column:balance;not null" json:"balance"`
	Version   int32          `gorm:"column:version;not null" json:"version"`
	Type      WalletType     `gorm:"column:type" json:"type"`
	Currency  WalletCurrency `gorm:"column:currency" json:"currency"`

	CustomerID int64 `gorm:"column:customer_id" json:"customer_id"`

	Customer *Customer       `gorm:"foreignKey:CustomerID;references:ID" json:"customer"`
	Changes  []*WalletChange `gorm:"foreignKey:WalletID;references:ID" json:"changes"`
}

// TableName Wallet's table name
func (*Wallet) TableName() string {
	return TableNameWallet
}

func (w *Wallet) BeforeCreate(tx *gorm.DB) error {
	if w.ID == "" {
		w.ID = uuid.NewString()
	}
	return nil
}

func (w *Wallet) BeforeSave(tx *gorm.DB) error {
	if w.ID == "" {
		w.ID = uuid.NewString()
	}
	return nil
}
