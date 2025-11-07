package model

import (
	"time"
)

const TableNameEWallet = "e-wallets"

// EWallet mapped from table <e-wallets>
type EWallet struct {
	CreatedAt   time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(6)" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP(6)" json:"updated_at"`
	Status      bool      `gorm:"column:status;not null;default:1" json:"status"`
	IsDeleted   bool      `gorm:"column:is_deleted;not null" json:"is_deleted"`
	ID          string    `gorm:"column:id;primaryKey" json:"id"`
	AccountName string    `gorm:"column:account_name" json:"account_name"`
	Email       string    `gorm:"column:email" json:"email"`
	Detail      string    `gorm:"column:detail" json:"detail"`
}

// TableName EWallet's table name
func (*EWallet) TableName() string {
	return TableNameEWallet
}
