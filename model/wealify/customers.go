package model

import (
	"core-ledger/model/enum"
	"time"

	"gorm.io/gorm"
)

const TableNameCustomer = "customers"

// Customer mapped from table <customers>
type Customer struct {
	CreatedAt                   time.Time                   `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(6)" json:"created_at,omitempty"`
	UpdatedAt                   time.Time                   `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP(6)" json:"updated_at,omitempty"`
	Status                      bool                        `gorm:"column:status;not null;default:1" json:"status,omitempty"`
	IsDeleted                   bool                        `gorm:"column:is_deleted;not null" json:"is_deleted,omitempty"`
	ID                          int64                       `gorm:"column:id;primaryKey;autoIncrement:true" json:"id,omitempty"`
	FullName                    string                      `gorm:"column:full_name;not null" json:"full_name,omitempty"`
	Email                       string                      `gorm:"column:email;not null" json:"email,omitempty"`
	PhoneNumber                 string                      `gorm:"column:phone_number" json:"phone_number,omitempty"`
	Address                     string                      `gorm:"column:address" json:"address,omitempty"`
	DateOfBirth                 time.Time                   `gorm:"column:date_of_birth" json:"date_of_birth,omitempty"`
	Password                    string                      `gorm:"column:password;not null" json:"password,omitempty"`
	TwoFactorStatus             string                      `gorm:"column:two_factor_status;not null;default:DISABLE" json:"two_factor_status,omitempty"`
	TwoFactorVerificationStatus TwoFactorVerificationStatus `gorm:"column:two_factor_verification_status;not null;default:UNVERIFIED" json:"two_factor_verification_status,omitempty"`
	TwoFactorMethod             string                      `gorm:"column:two_factor_method;not null;default:EMAIL" json:"two_factor_method,omitempty"`
	TwoFactorEnableFor          *string                     `gorm:"column:two_factor_enable_for;default:NULL" json:"two_factor_enable_for,omitempty"`
	AuthenticatorAppSecretKey   string                      `gorm:"column:authenticator_app_secret_key" json:"authenticator_app_secret_key,omitempty"`
	AuthenticatorAppDataURL     string                      `gorm:"column:authenticator_app_data_url" json:"authenticator_app_data_url,omitempty"`
	RegisteredAt                time.Time                   `gorm:"column:registered_at;not null" json:"registered_at,omitempty"`
	ChangedPwAt                 time.Time                   `gorm:"column:changed_pw_at;not null" json:"changed_pw_at,omitempty"`
	LastOnlineAt                time.Time                   `gorm:"column:last_online_at;not null" json:"last_online_at,omitempty"`
	Code                        string                      `gorm:"column:customer_id" json:"customer_id,omitempty"`
	AccountLevel                AccountLevel                `gorm:"column:account_level;not null;default:1" json:"account_level,omitempty"`
	AccountType                 string                      `gorm:"column:account_type;not null;default:INDIVIDUAL" json:"account_type,omitempty"`
	KyVerificationStatus        string                      `gorm:"column:ky_verification_status;not null;default:APPROVED" json:"ky_verification_status,omitempty"`
	CallingCodeID               *string                     `gorm:"column:calling_code_id" json:"calling_code_id,omitempty"`
	CountryID                   *string                     `gorm:"column:country_id" json:"country_id,omitempty"`
	LanguageID                  *string                     `gorm:"column:language_id" json:"language_id,omitempty"`
	FileID                      *string                     `gorm:"column:file_id" json:"file_id,omitempty"`
	VaEnable                    bool                        `gorm:"column:va_enable;not null" json:"va_enable,omitempty"`
	ReferralCode                string                      `gorm:"column:referral_code" json:"referral_code,omitempty"`
	ReferralByID                *string                     `gorm:"column:referral_by_id" json:"referral_by_id,omitempty"`
	Tier                        enum.Tier                   `gorm:"column:tier;not null;default:STANDARD" json:"tier,omitempty"`
	InHouse                     bool                        `gorm:"column:in_house;not null" json:"in_house,omitempty"`
	Type                        string                      `gorm:"column:type;not null;default:NORMAL" json:"type,omitempty"`
	WealifyWalletEnable         bool                        `gorm:"column:wealify_wallet_enable;not null" json:"wealify_wallet_enable,omitempty"`
	CountryCode                 string                      `gorm:"column:country_code;not null" json:"country_code,omitempty"`
	IsVcCustomer                bool                        `gorm:"column:is_vc_customer;not null" json:"is_vc_customer,omitempty"`
	IsVerified                  bool                        `gorm:"-" json:"is_verified,omitempty"`
	IsDefaultPassword           bool                        `gorm:"column:is_default_password;not null;default:false" json:"is_default_password,omitempty"`

	// VaWallet           *Wallet             `gorm:"foreignKey:CustomerID" json:"va_wallet,omitempty"`
	// MainWallet         *Wallet             `gorm:"foreignKey:CustomerID" json:"main_wallet,omitempty"`
	Wallets            []Wallet            `gorm:"foreignKey:CustomerID;references:ID" json:"wallets,omitempty"`
	IntegrationPartner *IntegrationPartner `gorm:"foreignKey:UserID" json:"integration_partner,omitempty"`
	CallingCode        *CallingCode        `gorm:"foreignKey:CallingCodeID;references:ID" json:"calling_code,omitempty"`
	Countrie           *Countrie           `gorm:"foreignKey:CountryID;references:ID" json:"countries,omitempty"`
	Language           *Language           `gorm:"foreignKey:LanguageID;references:ID" json:"languages,omitempty"`
	Files              *File               `gorm:"foreignKey:FileID;references:ID" json:"files,omitempty"`
	ReferralBy         *Customer           `gorm:"foreignKey:ReferralByID;references:ID" json:"referral_by,omitempty"`
	RefferedCustomers  []*Customer         `gorm:"foreignKey:ID;references:ReferralByID" json:"referred_cutomers,omitempty"`
}

// TableName Customer's table name
func (*Customer) TableName() string {
	return TableNameCustomer
}

func (c *Customer) BeforeCreate(tx *gorm.DB) error {
	c.CreatedAt = time.Now()
	return nil
}
