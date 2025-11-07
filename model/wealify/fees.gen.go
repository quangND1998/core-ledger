package model

import (
	"time"
)

const TableNameFee = "fees"

// Fee mapped from table <fees>
type Fee struct {
	CreatedAt       time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(6)" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP(6)" json:"updated_at"`
	Status          int32     `gorm:"column:status;not null;default:1" json:"status"`
	IsDeleted       int32     `gorm:"column:is_deleted;not null" json:"is_deleted"`
	ID              string    `gorm:"column:id;primaryKey" json:"id"`
	AccountType     string    `gorm:"column:account_type" json:"account_type"`
	AccountLevel    string    `gorm:"column:account_level" json:"account_level"`
	Provider        string    `gorm:"column:provider;not null;default:BANK" json:"provider"`
	ProviderType    string    `gorm:"column:provider_type;not null;default:INDIVIDUAL" json:"provider_type"`
	TransactionType string    `gorm:"column:transaction_type;not null" json:"transaction_type"`
	Description     string    `gorm:"column:description" json:"description"`
	CurrencySymbol  string    `gorm:"column:currency_symbol" json:"currency_symbol"`
	Tier            string    `gorm:"column:tier" json:"tier"`
}

// TableName Fee's table name
func (*Fee) TableName() string {
	return TableNameFee
}
