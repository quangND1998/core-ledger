package model

import (
	"core-ledger/model/enum"
	"time"

	"gorm.io/datatypes"
)

const TableNameSystemPayment = "`system-payments`"

// SystemPayment mapped from table <system-payments>
type SystemPayment struct {
	CreatedAt        time.Time                  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(6)" json:"created_at"`
	UpdatedAt        time.Time                  `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP(6)" json:"updated_at"`
	Status           int32                      `gorm:"column:status;not null;default:1" json:"status"`
	IsDeleted        int32                      `gorm:"column:is_deleted;not null" json:"is_deleted"`
	ID               string                     `gorm:"column:id;primaryKey" json:"id"`
	Email            string                     `gorm:"column:email" json:"email"`
	AccountLevels    string                     `gorm:"column:account_levels" json:"account_levels"`
	AccountTypes     string                     `gorm:"column:account_types" json:"account_types"`
	TransactionTypes string                     `gorm:"column:transaction_types" json:"transaction_types"`
	ProviderType     string                     `gorm:"column:provider_type;not null;default:BUSINESS" json:"provider_type"`
	PaymentType      string                     `gorm:"column:payment_type;not null;default:BANK" json:"payment_type"`
	Credential       *datatypes.JSON            `gorm:"column:credential" json:"credential"`
	Provider         enum.SystemPaymentProvider `gorm:"column:provider;not null;default:BANK" json:"provider"`
	BankID           *string                    `gorm:"column:bank_id" json:"bank_id"`
	EWalletID        *string                    `gorm:"column:e_wallet_id" json:"e_wallet_id"`
	Code             string                     `gorm:"column:system_payment_id;not null" json:"system_payment_id"`

	BankAccount         *BankAccount           `gorm:"foreignKey:BankID;references:ID" json:"bank"`
	EWallet             *EWallet               `gorm:"foreignKey:EWalletID;references:ID" json:"e_wallet"`
	SystemPaymentChange []*SystemPaymentChange `gorm:"foreignKey:SystemPaymentID;references:ID" json:"changes"`
	SystemPaymentLimit  []*SystemPaymentLimit  `gorm:"foreignKey:SystemPaymentID;references:ID" json:"limits"`
}

// TableName SystemPayment's table name
func (*SystemPayment) TableName() string {
	return TableNameSystemPayment
}
