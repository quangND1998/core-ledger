package model

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

const TableNamePayment = "payments"

// Payment mapped from table <payments>
type Payment struct {
	CreatedAt        time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(6)" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP(6)" json:"updated_at"`
	Status           int32     `gorm:"column:status;not null;default:1" json:"status"`
	IsDeleted        int32     `gorm:"column:is_deleted;not null" json:"is_deleted"`
	ID               string    `gorm:"column:id;primaryKey" json:"id"`
	Provider         string    `gorm:"column:provider;not null;default:BANK" json:"provider"`
	ProviderType     string    `gorm:"column:provider_type;not null;default:INDIVIDUAL" json:"provider_type"`
	PaymentStatus    string    `gorm:"column:payment_status;not null;default:PENDING" json:"payment_status"`
	PaymentType      string    `gorm:"column:payment_type;not null;default:BANK" json:"payment_type"`
	IsDefault        int32     `gorm:"column:is_default;not null" json:"is_default"`
	CustomerID       int32     `gorm:"column:customer_id" json:"customer_id"`
	CurrencySymbol   string    `gorm:"column:currency_symbol" json:"currency_symbol"`
	TransferMethodID string    `gorm:"column:transfer_method_id" json:"transfer_method_id"`
	BankID           *string   `gorm:"column:bank_id" json:"bank_id"`
	EWalletID        *string   `gorm:"column:e_wallet_id" json:"e_wallet_id"`
	PaymentID        string    `gorm:"column:payment_id;not null" json:"payment_id"`
}

// TableName Payment's table name
func (*Payment) TableName() string {
	return TableNamePayment
}

func (r *Payment) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.NewString()
	}
	return nil
}

func (r *Payment) BeforeSave(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.NewString()
	}
	return nil
}

type PaymentType string

const (
	PaymentTypeBank    PaymentType = "BANK"
	PaymentTypeEWallet PaymentType = "E_WALLET"
)

type BankCountry struct {
	Name     string `json:"name,omitempty" validate:"omitempty"`
	Code     int64  `json:"code" validate:"required"`
	IsoCode2 string `json:"iso_code2" validate:"required"`
	IsoCode3 string `json:"iso_code3" validate:"required"`
}

type BankPaymentDetail struct {
	BankCountry       BankCountry `json:"bank_country" validate:"required,dive"`
	BankName          string      `json:"bank_name,omitempty" validate:"omitempty"`
	BankCode          string      `json:"bank_code,omitempty" validate:"omitempty"`
	BankBin           string      `json:"bank_bin,omitempty" validate:"omitempty"`
	AccountName       string      `json:"account_name,omitempty" validate:"omitempty"`
	AccountNumber     string      `json:"account_number,omitempty" validate:"omitempty"`
	AchRoutingNumber  string      `json:"ach_routing_number,omitempty" validate:"omitempty"`
	RoutingNumber     string      `json:"routing_number,omitempty" validate:"omitempty"`
	RoutingBsbCode    string      `json:"routing_bsb_code,omitempty" validate:"omitempty"`
	SwiftBicCode      string      `json:"swift_bic_code,omitempty" validate:"omitempty"`
	InstitutionNumber string      `json:"institution_number,omitempty" validate:"omitempty"`
	TransitNumber     string      `json:"transit_number,omitempty" validate:"omitempty"`
	Iban              string      `json:"iban,omitempty" validate:"omitempty"`
	BranchCode        string      `json:"branch_code,omitempty" validate:"omitempty"`
	BranchName        string      `json:"branch_name,omitempty" validate:"omitempty"`
	SortCode          string      `json:"sort_code,omitempty" validate:"omitempty"`
	Note              string      `json:"note,omitempty" validate:"omitempty"`
}

// BankPaymentAccount is a strong type for transaction's payment field
type BankPaymentAccount struct {
	Detail      BankPaymentDetail `json:"detail" validate:"required,dive"`
	PaymentType PaymentType       `json:"payment_type" validate:"required,oneof=BANK E_WALLET"`
}

func (bp *BankPaymentAccount) IsValid() bool {
	return bp.Detail.BankCountry.Code == 0 &&
		bp.Detail.BankName != "" &&
		bp.Detail.AccountName != "" &&
		bp.Detail.AccountNumber != ""
}

type EWalletPaymentDetail struct {
	AccountName string `json:"account_name,omitempty" validate:"omitempty"`
	Email       string `json:"email" validate:"required,email"`
}

type EWalletPayment struct {
	Detail      EWalletPaymentDetail `json:"detail" validate:"required,dive"`
	PaymentType PaymentType          `json:"payment_type" validate:"required,oneof=BANK E_WALLET"`
}

func (ew *EWalletPayment) IsValid() bool {
	validate := validator.New()
	return validate.Var(ew.Detail.Email, "required,email") == nil
}
