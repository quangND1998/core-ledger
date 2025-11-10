package model

import (
	"time"
)

type CoaAccount struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Code      string         `gorm:"type:varchar(128);not null;uniqueIndex:uniq_code_currency" json:"code"`
	Name      string         `gorm:"type:varchar(256);not null" json:"name"`
	Type      string         `gorm:"type:varchar(16);not null;check:type IN ('ASSET','LIAB','EQUITY','REV','EXP')" json:"type"`
	Currency  string         `gorm:"type:char(8);not null;uniqueIndex:uniq_code_currency" json:"currency"`
	ParentID  *uint64        `gorm:"column:parent_id" json:"parent_id,omitempty"`
	Status    string         `gorm:"type:varchar(16);default:'ACTIVE';check:status IN ('ACTIVE','INACTIVE')" json:"status"`
	Provider  *string        `gorm:"type:varchar(64)" json:"provider,omitempty"`
	Network   *string        `gorm:"type:varchar(32)" json:"network,omitempty"`
	Tags      map[string]any `gorm:"type:jsonb" json:"tags,omitempty"`
	Metadata  map[string]any `gorm:"type:jsonb" json:"metadata,omitempty"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	// Quan hệ
	Parent   *CoaAccount  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []CoaAccount `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

// TableName đặt tên bảng rõ ràng
func (CoaAccount) TableName() string {
	return "coa_accounts"
}
