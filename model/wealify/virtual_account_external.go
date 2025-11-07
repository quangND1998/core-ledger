package model

import (
	"core-ledger/model/enum"
	"time"
)

const TableNameVirtualAccountExternal = "virtual_account_externals"

type VirtualAccountExternal struct {
	ID            string          `gorm:"primaryKey" json:"id"`
	CreatedAt     time.Time       `gorm:"created_at" json:"created_at"`
	UpdatedAt     time.Time       `gorm:"updated_at" json:"updated_at"`
	CardNumber    string          `json:"card_number"`
	CardHolder    string          `gorm:"not null" json:"card_holder"`
	Bank          enum.VABankCode `gorm:"type:varchar(20);not null;default:'BIDV'" json:"bank"`
	Provider      VAProvider      `gorm:"type:varchar(20);not null;default:'YOOBIL'" json:"provider"`
	TotalReceived float64         `gorm:"not null;default:0" json:"total_received"`
	TotalFee      float64         `gorm:"not null;default:0" json:"total_fee"`
	Status        VAStatus        `gorm:"type:varchar(20);not null;default:'PENDING'" json:"status"`
	OrderNo       string          `gorm:"not null;unique" json:"order_no"`
	ApprovedAt    *time.Time      `json:"approved_at"`
	RejectedAt    *time.Time      `json:"rejected_at"`
	FeeName       *string         `json:"fee_name"`
	RestrictedAt  *time.Time      `json:"restricted_at"`
	InactiveAt    *time.Time      `json:"inactive_at"`
	IsDeleted     bool            `gorm:"default:false" json:"is_deleted"`

	DelegatedBy *string `json:"delegated_by"`
	PlatformID  string  `json:"platform_id"`
}

// TableName VirtualAccount's table name
func (*VirtualAccountExternal) TableName() string {
	return TableNameVirtualAccountExternal
}
