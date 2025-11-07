package model

import (
	"time"
)

type QrMerchant struct {
	ID                int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name              string     `gorm:"type:varchar(255);not null" json:"name"`
	Code              string     `gorm:"type:varchar(255);not null" json:"code"`
	Note              string     `gorm:"type:varchar(255);not null" json:"note"`
	ImageURL          string     `gorm:"type:text;not null" json:"image_url"`
	CreatedByID       int64      `gorm:"not null" json:"created_by_id"`
	GTelBankID        int64      `gorm:"not null" json:"g_tel_bank_id"`
	BankAccountName   string     `gorm:"type:varchar(255);not null" json:"bank_account_name"`
	BankAccountNumber string     `gorm:"type:varchar(255);not null" json:"bank_account_number"`
	Status            bool       `gorm:"type:tinyint(1);not null" json:"status"`
	CreatedAt         time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt         *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
}

// TableName overrides the table name used by GORM
func (QrMerchant) TableName() string {
	return "qr_merchants"
}
