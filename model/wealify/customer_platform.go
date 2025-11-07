package model

const TableNameCustomerPlatform = "customer_platforms"

type CustomerPlatform struct {
	CustomerID int64     `gorm:"column:customer_id;not null;index" json:"customer_id"`
	PlatformID string    `gorm:"column:platform_id;not null;index" json:"platform_id"`
	Customer   *Customer `gorm:"foreignKey:CustomerID;references:ID" json:"customer,omitempty"`
	Platform   *Platform `gorm:"foreignKey:PlatformID;references:ID" json:"platform,omitempty"`
}

// TableName Setting's table name
func (CustomerPlatform) TableName() string {
	return TableNameCustomerPlatform
}
