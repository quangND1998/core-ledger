package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

const TableNameBank = "banks"

type BankAccount struct {
	CreatedAt         time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(6)" json:"created_at"`
	UpdatedAt         time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP(6)" json:"updated_at"`
	Status            int32     `gorm:"column:status;not null;default:1" json:"status"`
	IsDeleted         int32     `gorm:"column:is_deleted;not null" json:"is_deleted"`
	ID                string    `gorm:"column:id;primaryKey" json:"id"`
	BankName          string    `gorm:"column:bank_name" json:"bank_name"`
	BankBin           string    `gorm:"column:bank_bin" json:"bank_bin"`
	AccountName       string    `gorm:"column:account_name" json:"account_name"`
	AchRoutingNumber  string    `gorm:"column:ach_routing_number" json:"ach_routing_number"`
	SwiftBicCode      string    `gorm:"column:swift_bic_code" json:"swift_bic_code"`
	InstitutionNumber string    `gorm:"column:institution_number" json:"institution_number"`
	TransitNumber     string    `gorm:"column:transit_number" json:"transit_number"`
	Iban              string    `gorm:"column:iban" json:"iban"`
	BranchCode        string    `gorm:"column:branch_code" json:"branch_code"`
	SortCode          string    `gorm:"column:sort_code" json:"sort_code"`
	Note              string    `gorm:"column:note" json:"note"`
	BankCode          string    `gorm:"column:bank_code" json:"bank_code"`
	RoutingNumber     string    `gorm:"column:routing_number" json:"routing_number"`
	RoutingBsbCode    string    `gorm:"column:routing_bsb_code" json:"routing_bsb_code"`
	BranchName        string    `gorm:"column:branch_name" json:"branch_name"`
	AccountNumber     string    `gorm:"column:account_number" json:"account_number"`

	CountryID string `gorm:"column:country_id" json:"country_id"`

	Country *Countrie `gorm:"foreignKey:CountryID;references:ID" json:"country"`
}

// TableName Bank's table name
func (*BankAccount) TableName() string {
	return TableNameBank
}

func (ba *BankAccount) BeforeSave(tx *gorm.DB) (err error) {
	if ba.ID == "" {
		ba.ID = uuid.NewString()
	}
	return
}

func (ba *BankAccount) BeforeCreate(tx *gorm.DB) (err error) {
	if ba.ID == "" {
		ba.ID = uuid.NewString()
	}
	return
}

func (ba *BankAccount) ConvertBankPayment() *BankPaymentAccount {
	var bankCountry BankCountry
	if ba.Country != nil {
		bankCountry.Name = ba.Country.Name
		bankCountry.Code = ba.Country.Code
		bankCountry.IsoCode2 = ba.Country.IsoCode2
		bankCountry.IsoCode3 = ba.Country.IsoCode3
	}
	return &BankPaymentAccount{
		Detail: BankPaymentDetail{
			BankCountry:       bankCountry,
			BankName:          ba.BankName,
			BankCode:          ba.BankCode,
			BankBin:           ba.BankBin,
			AccountName:       ba.AccountName,
			AccountNumber:     ba.AccountNumber,
			AchRoutingNumber:  ba.AchRoutingNumber,
			RoutingNumber:     ba.RoutingNumber,
			RoutingBsbCode:    ba.RoutingBsbCode,
			SwiftBicCode:      ba.SwiftBicCode,
			InstitutionNumber: ba.InstitutionNumber,
			TransitNumber:     ba.TransitNumber,
			Iban:              ba.Iban,
			BranchCode:        ba.BranchCode,
			BranchName:        ba.BranchName,
			SortCode:          ba.SortCode,
			Note:              ba.Note,
		},
		PaymentType: PaymentTypeBank,
	}
}
