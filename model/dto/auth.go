package dto

import (
	"github.com/golang-jwt/jwt/v5"
)

type RefreshClaims struct {
	jwt.RegisteredClaims
	ID     int64  `json:"id"`
	Secret string `json:"secret"`
}

type Claims struct {
	ID                  int64           `json:"id"`
	SessionId           string          `json:"session_id"`
	Secret              string          `json:"secret"`
	IsEmployee          bool            `json:"is_employee"`
	TwoFactorEnableFor  interface{}     `json:"two_factor_enable_for"`
	TwoFactorStatus     TwoFactorStatus `json:"two_factor_status"`
	DelegationAccountID string          `json:"delegationAccountId"`
	jwt.RegisteredClaims
}
type TwoFactorEnableFor []string

type LoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterBody struct {
	FullName     string `json:"full_name"`
	AccountType  string `json:"account_type"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	PhoneNumber  string `json:"phone_number"`
	ReferralCode string `json:"referral_code"`
}

type LoginResponseBody struct {
	NeedCompleteRegister bool               `json:"need_complete_register"`
	OneTimeToken         *string            `json:"one_time_token,omitempty"`
	Email                *string            `json:"email,omitempty"`
	IsEmailVerified      *bool              `json:"is_email_verified,omitempty"`
	Provider             *string            `json:"provider,omitempty"`
	ProviderUID          *string            `json:"provider_uid,omitempty"`
	FullName             *string            `json:"full_name,omitempty"`
	AccessToken          string             `json:"access_token,omitempty"`
	AccessExpiredAt      int64              `json:"access_expired_at,omitempty"`
	AccessExpiredIn      int64              `json:"access_expired_in,omitempty"`
	RefreshToken         string             `json:"refresh_token,omitempty"`
	RefreshExpiredAt     *int64             `json:"refresh_expired_at,omitempty"` // pointer nếu có thể null
	TwoFactorStatus      string             `json:"two_factor_status,omitempty"`
	TwoFactorMethod      string             `json:"two_factor_method,omitempty"`
	TwoFactorEnableFor   TwoFactorEnableFor `json:"two_factor_enable_for"`
	Status               *bool              `json:"status,omitempty"`
	NewDevice            *bool              `json:"new_device,omitempty"`
}
type CompleteRegisterReq struct {
	OneTimeToken string `json:"one_time_token"`
	AccountType  string `json:"account_type"`
	PhoneNumber  string `json:"phone_number" binding:"required"`
	FullName     string `json:"full_name,omitempty"`
	ReferralCode string `json:"referral_code"`
	AcceptTnC    bool   `json:"accept_tnc"`
}
