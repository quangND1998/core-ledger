package model

import "time"

type SocialAccount struct {
	ID           uint64     `gorm:"id"`
	CustomerID   int64      `gorm:"user_id"`
	Provider     string     `gorm:"provider"` // ENUM in DB, string in Go
	ProviderUID  string     `gorm:"provider_uid"`
	Email        string     `gorm:"email"`
	AccessToken  *string    `gorm:"access_token"`
	RefreshToken *string    `gorm:"refresh_token"`
	ExpiresAt    *time.Time `gorm:"expires_at"`
	LinkedAt     time.Time  `gorm:"linked_at"`
}

type FindMatchedSocialAccount struct {
	CustomerID             int64   `json:"customer_id"`
	CustomerTwoFAStatus    string  `json:"customer_two_fa_status"`
	CustomerTwoFAMethod    string  `json:"customer_two_fa_method"`
	CustomerTwoFAEnableFor *string `json:"customer_two_fa_enable_for"`
	CustomerStatus         bool    `json:"customer_status"`
	Email                  string  `json:"email"`
	Provider               string  `json:"provider"`
	IsEmailVerified        *bool   `json:"is_email_verified"`
	IsNewUser              bool    `json:"-"`
}

func (t *SocialAccount) TableName() string {
	return "social_accounts"
}
