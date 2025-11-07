package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

const TableNameVcWallet = "vc-wallets"

// VcWallet mapped from table <vc-wallets>
type VcWallet struct {
	ID        string    `gorm:"column:id;primaryKey" json:"id"`
	UserID    int32     `gorm:"column:user_id" json:"user_id"`
	Balance   float64   `gorm:"column:balance;not null" json:"balance"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP(6)" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP(6)" json:"updated_at"`
	IsDeleted bool      `gorm:"column:is_deleted;not null" json:"is_deleted"`
	Version   int32     `gorm:"column:version;not null" json:"version"`
}

// TableName VcWallet's table name
func (*VcWallet) TableName() string {
	return TableNameVcWallet
}

func (v *VcWallet) BeforeSave(tx *gorm.DB) (err error) {
	if v.ID == "" {
		v.ID = uuid.NewString()
	}
	return
}

func (v *VcWallet) BeforeCreate(tx *gorm.DB) (err error) {
	if v.ID == "" {
		v.ID = uuid.NewString()
	}
	return
}
