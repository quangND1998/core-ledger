package model

import (
	"gorm.io/gorm"
	"time"
)

const TableNameIntegrationPartner = "integration_partners"

type IntegrationPartner struct {
	CreatedAt              time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(6)" json:"created_at"`
	UpdatedAt              time.Time  `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP(6)" json:"updated_at"`
	Status                 bool       `gorm:"column:status;not null;default:1" json:"status"`
	ExpiresAt              time.Time  `gorm:"column:expires_at" json:"expires_at"`
	ID                     int64      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name                   string     `gorm:"column:name" json:"name"`
	APIKey                 string     `gorm:"column:api_key" json:"api_key"`
	APIKeyHash             string     `gorm:"column:key_hash;not null" json:"key_hash"`
	APIKeyLastUsedAt       *time.Time `gorm:"column:last_used_at" json:"last_used_at"`
	UserID                 int64      `gorm:"column:user_id" json:"user_id"`
	PublicKey              string     `gorm:"column:public_key" json:"public_key"`
	WebhookURL             int64      `gorm:"column:webhook_url" json:"webhook_url"`
	NeedMappingTransaction bool       `gorm:"column:need_mapping" json:"need_mapping"`
}

// TableName Customer's table name
func (*IntegrationPartner) TableName() string {
	return TableNameIntegrationPartner
}

func (c *IntegrationPartner) BeforeCreate(tx *gorm.DB) error {
	c.CreatedAt = time.Now()
	return nil
}
