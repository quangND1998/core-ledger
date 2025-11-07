package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

const TableNamePlatformFee = "platform_fees"

type PlatformFee struct {
	ID                  string     `gorm:"column:id;primaryKey" json:"id,omitempty"`
	Status              string     `gorm:"column:status;default:ACTIVE" json:"status,omitempty"`
	TopUpStandardFee    *float64   `gorm:"column:top_up_standard_fee" json:"top_up_standard_fee,omitempty"`
	TopUpSilverFee      *float64   `gorm:"column:top_up_silver_fee" json:"top_up_silver_fee,omitempty"`
	TopUpGoldFee        *float64   `gorm:"column:top_up_gold_fee" json:"top_up_gold_fee,omitempty"`
	TopUpDiamondFee     *float64   `gorm:"column:top_up_diamond_fee" json:"top_up_diamond_fee,omitempty"`
	WithdrawStandardFee *float64   `gorm:"column:withdraw_standard_fee" json:"withdraw_standard_fee,omitempty"`
	WithdrawSilverFee   *float64   `gorm:"column:withdraw_silver_fee" json:"withdraw_silver_fee,omitempty"`
	WithdrawGoldFee     *float64   `gorm:"column:withdraw_gold_fee" json:"withdraw_gold_fee,omitempty"`
	WithdrawDiamondFee  *float64   `gorm:"column:withdraw_diamond_fee" json:"withdraw_diamond_fee,omitempty"`
	Provider            VAProvider `gorm:"column:provider;not null" json:"provider,omitempty"`
	CreatedAt           time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(6)" json:"created_at,omitempty"`
	UpdatedAt           time.Time  `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP(6)" json:"updated_at,omitempty"`

	PlatformID string `gorm:"column:platform_id;not null" json:"platform_id"`

	Platform *Platform `gorm:"foreignKey:PlatformID;references:id" json:"platform"`
}

// TableName Platform's table name
func (*PlatformFee) TableName() string {
	return TableNamePlatformFee
}

func (pf *PlatformFee) BeforeCreate(tx *gorm.DB) error {
	if pf.ID == "" {
		pf.ID = uuid.NewString()
	}
	return nil
}
