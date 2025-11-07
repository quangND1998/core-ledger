package model

const TableNameBankInfo = "bank_info"

type BankInfo struct {
	ID                string `gorm:"column:id;primaryKey" json:"id"`
	Name              string `gorm:"column:name" json:"name"`
	InternationalName string `gorm:"column:international_name" json:"international_name"`
	ShortName         string `gorm:"column:short_name" json:"short_name"`
	Code              string `gorm:"column:code" json:"code"`
	SwiftCode         string `gorm:"column:swift_code" json:"swift_code"`
	Bin               string `gorm:"column:bin" json:"bin"`
	CountryCode       string `gorm:"column:country_code" json:"country_code"`
}

// TableName Bank's table name
func (*BankInfo) TableName() string {
	return TableNameBankInfo
}
