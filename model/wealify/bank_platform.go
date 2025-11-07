package model

const TableNameBankPlatform = "bank_platforms"

type BankPlatform struct {
	BankID     string `json:"bank_id"`
	PlatformID string `json:"name"`

	Bank     *BankInfo `gorm:"foreignKey:BankID;references:ID" json:"bank"`
	Platform *Platform `gorm:"foreignKey:PlatformID;references:ID" json:"platform"`
}

// TableName Bank's table name
func (*BankPlatform) TableName() string {
	return TableNameBankPlatform
}
