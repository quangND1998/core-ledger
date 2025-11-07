package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

const TableNameVirtualCard = "virtual-cards"

// VirtualCard mapped from table <virtual-cards>
type VirtualCard struct {
	ID          string    `gorm:"column:id;primaryKey" json:"id"`
	CardName    string    `gorm:"column:card_name;not null" json:"card_name"`
	CardNumber  string    `gorm:"column:card_number;not null" json:"card_number"`
	LastFour    string    `gorm:"column:last_four;not null" json:"last_four"`
	Category    string    `gorm:"column:category;not null" json:"category"`
	Cvv         string    `gorm:"column:cvv;not null" json:"cvv"`
	ExpiryDate  string    `gorm:"column:expiry_date;not null" json:"expiry_date"`
	CountryCode string    `gorm:"column:country_code;not null" json:"country_code"`
	ReferentID  string    `gorm:"column:referent_id;not null" json:"referent_id"`
	Email       string    `gorm:"column:email;not null" json:"email"`
	PhoneNumber string    `gorm:"column:phone_number;not null" json:"phone_number"`
	CardPurpose string    `gorm:"column:card_purpose;not null" json:"card_purpose"`
	CardType    string    `gorm:"column:card_type;default:VIRTUAL" json:"card_type"`
	CardStatus  string    `gorm:"column:card_status;default:ACTIVE" json:"card_status"`
	Balance     float64   `gorm:"column:balance;not null" json:"balance"`
	CreatorID   int32     `gorm:"column:creator_id;not null" json:"creator_id"`
	CreatedAt   time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(6)" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP(6)" json:"updated_at"`
	IsDeleted   bool      `gorm:"column:is_deleted;not null" json:"is_deleted"`
	SpendLimit  float64   `gorm:"column:spend_limit;not null" json:"spend_limit"`

	CryptoWallets []*CryptoWallet `gorm:"foreignKey:VirtualCardID" json:"crypto_wallets"`
	Transactions  []*Transaction  `gorm:"foreignKey:VirtualCardID" json:"transactions"`
}

// TableName VirtualCard's table name
func (*VirtualCard) TableName() string {
	return TableNameVirtualCard
}

func (t *VirtualCard) BeforeCreate(tx *gorm.DB) (err error) {
	t.UpdatedAt = time.Now()
	if t.ID == "" {
		t.ID = uuid.NewString()
	}
	return
}

func (t *VirtualCard) BeforeSave(tx *gorm.DB) (err error) {
	t.UpdatedAt = time.Now()
	if t.ID == "" {
		t.ID = uuid.NewString()
	}
	return
}

func (t *VirtualCard) BeforeUpdate(tx *gorm.DB) (err error) {
	t.UpdatedAt = time.Now()
	return
}
