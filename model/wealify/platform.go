package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

const TableNamePlatform = "platforms"

type Platform struct {
	ID          string    `gorm:"column:id;primaryKey" json:"id,omitempty"`
	Name        string    `gorm:"column:name;not null" json:"name,omitempty" `
	Type        string    `gorm:"column:type;default:CUSTOM" json:"type,omitempty"`
	Status      string    `gorm:"column:status;default:ACTIVE" json:"status,omitempty"`
	StandardFee *float64  `gorm:"column:standard_fee" json:"standard_fee,omitempty"`
	SilverFee   *float64  `gorm:"column:silver_fee" json:"silver_fee,omitempty"`
	GoldFee     *float64  `gorm:"column:gold_fee" json:"gold_fee,omitempty"`
	DiamondFee  *float64  `gorm:"column:diamond_fee" json:"diamond_fee,omitempty"`
	CreatedAt   time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(6)" json:"created_at,omitempty"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP(6);autoCreateTime;autoUpdateTime" json:"updated_at,omitempty"`
	IsDeleted   bool      `gorm:"column:is_deleted;not null" json:"is_deleted,omitempty"`

	Customers         []*Customer         `gorm:"many2many:customer_platforms" json:"customers,omitempty"`
	CustomerPlatforms []*CustomerPlatform `gorm:"foreignKey:PlatformID;references:ID" json:"customer_platforms,omitempty"`
	PlatformFee       []*PlatformFee      `gorm:"foreignKey:PlatformID;references:ID" json:"platform_fee,omitempty"`
}

// TableName Platform's table name
func (*Platform) TableName() string {
	return TableNamePlatform
}

func (p *Platform) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	return nil
}

func (p *Platform) AfterCreate(tx *gorm.DB) error {
	var fees []*PlatformFee

	for _, provider := range VAProviderValues() {
		fees = append(fees, &PlatformFee{
			ID:         uuid.New().String(),
			Provider:   provider,
			PlatformID: p.ID,
		})
	}
	if err := tx.Model(&PlatformFee{}).Save(fees).Error; err != nil {
		return err
	}
	return nil
}

func (p *Platform) BeforeSave(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Platform) BeforeDelete(tx *gorm.DB) error {
	if err := tx.Where("platform_id = ?", p.ID).Delete(&PlatformFee{}).Error; err != nil {
		return err
	}
	if err := tx.Where("platform_id = ?", p.ID).Delete(&CustomerPlatform{}).Error; err != nil {
		return err
	}
	if err := tx.Where("platform_id = ?", p.ID).Delete(&BankPlatform{}).Error; err != nil {
		return err
	}
	return nil
}
