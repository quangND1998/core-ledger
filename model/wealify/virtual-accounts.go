package model

import (
	"core-ledger/model/enum"
	"time"
)

const TableNameVirtualAccount = "`virtual-accounts`"

type VABank string

const (
	VABankBIDV  VABank = "BIDV"
	VABankVCB   VABank = "VCB"
	VABankTCB   VABank = "TCB"
	VABankACB   VABank = "ACB"
	VABankMB    VABank = "MB"
	VABankVPB   VABank = "VPB"
	VABankSTB   VABank = "STB"
	VABankTPB   VABank = "TPB"
	VABankOCB   VABank = "OCB"
	VABankMSB   VABank = "MSB"
	VABankHDB   VABank = "HDB"
	VABankSHB   VABank = "SHB"
	VABankVIB   VABank = "VIB"
	VABankEIB   VABank = "EIB"
	VABankBID   VABank = "BID"
	VABankTCB2  VABank = "TCB2"
	VABankTPB2  VABank = "TPB2"
	VABankVPB2  VABank = "VPB2"
	VABankMB2   VABank = "MB2"
	VABankACB2  VABank = "ACB2"
	VABankVCB2  VABank = "VCB2"
	VABankBIDV2 VABank = "BIDV2"
	VABankSTB2  VABank = "STB2"
	VABankHDB2  VABank = "HDB2"
	VABankSHB2  VABank = "SHB2"
	VABankVIB2  VABank = "VIB2"
	VABankEIB2  VABank = "EIB2"
	VABankBID2  VABank = "BID2"
	VABankOCB2  VABank = "OCB2"
	VABankMSB2  VABank = "MSB2"
	VABankTPB3  VABank = "TPB3"
	VABankVPB3  VABank = "VPB3"
	VABankMB3   VABank = "MB3"
	VABankACB3  VABank = "ACB3"
	VABankVCB3  VABank = "VCB3"
	VABankBIDV3 VABank = "BIDV3"
	VABankSTB3  VABank = "STB3"
	VABankHDB3  VABank = "HDB3"
	VABankSHB3  VABank = "SHB3"
	VABankVIB3  VABank = "VIB3"
	VABankEIB3  VABank = "EIB3"
	VABankBID3  VABank = "BID3"
	VABankOCB3  VABank = "OCB3"
	VABankMSB3  VABank = "MSB3"
)

type VAProvider string

const (
	VAProviderYoobil VAProvider = "YOOBIL"
	VAProviderNeox   VAProvider = "NEOX"
	VAProviderHPay   VAProvider = "H_PAY"
	VAProviderGTel   VAProvider = "G_TEL"
	VAProvider9Pay   VAProvider = "9_PAY"
)

func (t VAProvider) String() string {
	return string(t)
}

func VAProviderValues() []VAProvider {
	return []VAProvider{VAProviderYoobil, VAProviderNeox, VAProviderHPay, VAProviderGTel, VAProvider9Pay}
}

type VAStatus string

const (
	VAStatusActive     VAStatus = "ACTIVE"
	VAStatusInactive   VAStatus = "INACTIVE"
	VAStatusRestricted VAStatus = "RESTRICTED"
	VAStatusPending    VAStatus = "PENDING"
	VAStatusProcess    VAStatus = "PROCESS"
	VAStatusRejected   VAStatus = "REJECTED"
)

type VirtualAccount struct {
	ID            int64           `gorm:"primarykey" json:"id"`
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

	// Foreign Key
	CustomerID     int64   `json:"customer_id"`
	CurrencySymbol string  `gorm:"not null;default:'VND'" json:"currency_symbol"`
	DelegatedBy    *string `json:"delegated_by"`
	PlatformID     string  `json:"platform_id"`

	// Relationships
	BankInfo    *BankInfo              `gorm:"foreignKey:Bank;references:Code" json:"bank_info"`
	Customer    *Customer              `gorm:"foreignKey:CustomerID;references:ID" json:"customer,omitempty"`
	Currency    *Currency              `gorm:"foreignKey:CurrencySymbol;references:Symbol" json:"currency,omitempty"`
	Platform    *Platform              `gorm:"foreignKey:PlatformID;references:ID" json:"platform,omitempty"`
	PlatformFee *PlatformFee           `gorm:"foreignKey:PlatformID,Provider;references:PlatformID,Provider" json:"platform_fee,omitempty"`
	Events      []*VirtualAccountEvent `gorm:"foreignKey:VirtualAccountID;references:ID" json:"events,omitempty"`
}

// TableName VirtualAccount's table name
func (*VirtualAccount) TableName() string {
	return TableNameVirtualAccount
}

type StatsNumberVAByCustomer struct {
	TotalCreatedVA      int `gorm:"total_created_"`
	TotalBIDVRestricted int `gorm:"total_bidv_restricted"`
}
