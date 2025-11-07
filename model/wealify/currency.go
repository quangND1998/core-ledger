package model

import (
	"time"
)

const TableNameCurrency = "currencies"

type Currency struct {
	CreatedAt        time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(6)" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP(6)" json:"updated_at"`
	Status           bool      `gorm:"column:status;not null;default:1" json:"status"`
	IsDeleted        bool      `gorm:"column:is_deleted;not null" json:"is_deleted"`
	ID               string    `gorm:"column:id;not null" json:"id"`
	Name             string    `gorm:"column:name" json:"name"`
	Symbol           string    `gorm:"column:symbol;primaryKey" json:"symbol"`
	Code             string    `gorm:"column:code;not null" json:"code"`
	ProviderTypes    string    `gorm:"column:provider_types;not null" json:"provider_types"`
	TransactionTypes string    `gorm:"column:transaction_types;not null" json:"transaction_types"`
	FileID           string    `gorm:"column:file_id" json:"file_id"`
	Providers        string    `gorm:"column:providers;not null" json:"providers"`
}

// TableName Currency's table name
func (*Currency) TableName() string {
	return TableNameCurrency
}
