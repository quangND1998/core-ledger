package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

const TableNameCryptoWallet = "crypto-wallets"

// CryptoWallet mapped from table <crypto-wallets>
type CryptoWallet struct {
	ID            string    `gorm:"column:id;primaryKey" json:"id"`
	Network       string    `gorm:"column:network;default:SOL" json:"network"`
	Address       string    `gorm:"column:address" json:"address"`
	UserID        int32     `gorm:"column:user_id" json:"user_id"`
	CreatedAt     time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(6)" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP(6)" json:"updated_at"`
	IsDeleted     bool      `gorm:"column:is_deleted;not null" json:"is_deleted"`
	VirtualCardID *string   `gorm:"column:virtual_card_id" json:"virtual_card_id"`
	ImageUrl      string    `gorm:"column:image_url;not null" json:"image_url"`

	VirtualCard  *VirtualCard   `gorm:"foreignKey:id;references:VirtualCardID" json:"virtual_card"`
	Transactions []*Transaction `gorm:"foreignKey:CryptoWalletID" json:"transactions"`
}

// TableName CryptoWallet's table name
func (*CryptoWallet) TableName() string {
	return TableNameCryptoWallet
}

func (cw *CryptoWallet) BeforeSave(tx *gorm.DB) (err error) {
	if cw.ID == "" {
		cw.ID = uuid.NewString()
	}
	return
}

func (cw *CryptoWallet) BeforeCreate(tx *gorm.DB) (err error) {
	if cw.ID == "" {
		cw.ID = uuid.NewString()
	}
	return
}

func (cw *CryptoWallet) IsCard() bool {
	return cw.VirtualCardID != nil
}
