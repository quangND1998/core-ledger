package model

import (
	"strings"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CoaAccount struct {
	Entity
	ID        uint64          `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Code      string          `gorm:"type:varchar(128);not null;uniqueIndex:uniq_code_currency" json:"code"`
	AccountNo string          `gorm:"type:varchar(64);uniqueIndex" json:"account_no"`
	Name      string          `gorm:"type:varchar(256);not null" json:"name"`
	Type      string          `gorm:"type:varchar(16);not null;check:type IN ('ASSET','LIAB','EQUITY','REV','EXP')" json:"type"`
	Currency  string          `gorm:"type:char(8);not null;uniqueIndex:uniq_code_currency" json:"currency"`
	ParentID  *uint64         `gorm:"column:parent_id" json:"parent_id,omitempty"`
	Status    string          `gorm:"type:varchar(16);default:'ACTIVE';check:status IN ('ACTIVE','INACTIVE')" json:"status"`
	Provider  *string         `gorm:"type:varchar(64)" json:"provider,omitempty"`
	Network   *string         `gorm:"type:varchar(32)" json:"network,omitempty"`
	Tags      map[string]any  `gorm:"type:jsonb" json:"tags,omitempty"`
	Metadata  *datatypes.JSON `gorm:"type:jsonb" json:"metadata,omitempty"`
	CreatedAt time.Time       `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time       `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	// Quan hệ
	Parent   *CoaAccount  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []CoaAccount `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Entries  []Entry      `gorm:"foreignKey:AccountID" json:"entries"`
	Journals []Journal    `gorm:"-" json:"journals"`
}

// TableName đặt tên bảng rõ ràng
func (c *CoaAccount) TableName() string {
	return "coa_accounts"
}

func (c *CoaAccount) ScopeSearch(search string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if strings.TrimSpace(search) == "" {
			return db
		}
		likeQuery := "%" + search + "%"
		return db.Where(db.
			Where("name LIKE ?", likeQuery).
			Or("code LIKE ?", likeQuery).
			Or("account_no LIKE ?", likeQuery))
	}
}

func (c *CoaAccount) ScopeStatus(status []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(status) == 0 {
			return db
		}
		return db.Where("status IN ?", status)
	}
}

func (c *CoaAccount) ScopeTypes(types []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(types) == 0 {
			return db
		}
		return db.Where("type IN ?", types) // Postgres array overlap example
	}
}

func (c *CoaAccount) ScopeProviders(providers []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(providers) == 0 {
			return db
		}
		return db.Where("provider IN ?", providers) // Postgres array overlap example
	}
}

func (c *CoaAccount) ScopeSort(sortStr string) func(db *gorm.DB) *gorm.DB {
	return c.Entity.ScopeSort(sortStr, CoaAccount{})
}
