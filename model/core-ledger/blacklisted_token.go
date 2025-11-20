package model

import (
	"time"

	"gorm.io/gorm"
)

// BlacklistedToken represents a blacklisted JWT token
type BlacklistedToken struct {
	Entity
	ID        uint64    `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	TokenHash string    `gorm:"type:varchar(64);not null;uniqueIndex" json:"token_hash"` // SHA256 hash of token
	UserID    uint64    `gorm:"not null;index" json:"user_id"`
	ExpiresAt time.Time `gorm:"not null;index" json:"expires_at"` // Token expiration time
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (b *BlacklistedToken) TableName() string {
	return "blacklisted_tokens"
}

// BeforeCreate hook to ensure token hash is set
func (b *BlacklistedToken) BeforeCreate(tx *gorm.DB) error {
	if b.TokenHash == "" {
		return gorm.ErrRecordNotFound
	}
	return nil
}


