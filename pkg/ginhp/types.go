package ginhp

import (
	"core-ledger/model/dto"
	model "core-ledger/model/wealify"
	"time"
)

type VerificationStatus string

const (
	Unverified VerificationStatus = "UNVERIFIED"
	Verified   VerificationStatus = "VERIFIED"
)

type TwoFactorMethod string

const (
	MethodNone TwoFactorMethod = "NONE"
	MethodSMS  TwoFactorMethod = "SMS"
	MethodApp  TwoFactorMethod = "APP"
)

type CallingCodeRelation struct {
	ID     int    `json:"id"`
	Code   string `json:"code"`
	Region string `json:"region"`
}

type AccountRelation struct {
	ID int64 `json:"id"` // hoặc uuid.UUID tùy theo database
}

type AccountRequest struct {
	FullName                    string              `json:"full_name" validate:"required" example:"John Doe"`
	Email                       string              `json:"email" validate:"required,email" example:"john@example.com"`
	CallingCode                 CallingCodeRelation `json:"calling_code"`
	PhoneNumber                 string              `json:"phone_number" validate:"required" example:"123456789"`
	TwoFactorStatus             dto.TwoFactorStatus `json:"two_factor_status" validate:"required,oneof=DISABLED ENABLED"`
	AuthenticatorAppSecretKey   string              `json:"authenticator_app_secret_key" validate:"omitempty" example:"ABCD1234SECRET"`
	TwoFactorVerificationStatus VerificationStatus  `json:"two_factor_verification_status" validate:"required,oneof=UNVERIFIED VERIFIED"`
	TwoFactorMethod             TwoFactorMethod     `json:"two_factor_method" validate:"required,oneof=NONE SMS APP"`
	RegisteredAt                time.Time           `json:"registered_at" example:"2025-07-09T15:04:05Z"`
	Status                      bool                `json:"status" example:"true"`
	IsDeleted                   bool                `json:"is_deleted" example:"false"`
	IsEmployee                  bool                `json:"is_employee" example:"true"`
	Customer                    *model.Customer     `json:"customer"`
	Employee                    *model.Employee     `json:"employee"`
}

type CustomerRequest struct {
	*model.Customer     `json:"customer"`
	UseXApiKey          bool    `json:"use_x_api_key"`
	DelegationAccountID *string `json:"delegation_account_id"`
}

type EmployeeRequest struct {
	ID int64 `json:"id"`
}
