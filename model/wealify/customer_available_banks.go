package model

const TableNameCustomerAvailableBank = "customer_available_banks"

type CustomerAvailableBank struct {
	CustomerID int64  `gorm:"column:customer_id" json:"customer_id"`
	BankID     string `gorm:"column:bank_id" json:"bank_code"`

	Bank *BankInfo `gorm:"foreignkey:BankID;references:ID" json:"bank_info"`
}

// TableName Setting's table name
func (*CustomerAvailableBank) TableName() string {
	return TableNameCustomerAvailableBank
}
