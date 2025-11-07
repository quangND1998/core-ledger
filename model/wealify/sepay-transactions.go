package model

const TableNameSepayTransaction = "sepay-transactions"

type SepayTransaction struct {
	ID              int64    `gorm:"column:id;primaryKey" json:"id"`
	Gateway         string   `gorm:"column:gateway" json:"gateway"`
	TransactionDate string   `gorm:"column:transactionDate" json:"transactionDate"`
	AccountNumber   string   `gorm:"column:accountNumber" json:"accountNumber"`
	Code            *string  `gorm:"column:code" json:"code"`
	Content         string   `gorm:"column:content" json:"content"`
	TransferType    string   `gorm:"column:transferType" json:"transferType"`
	TransferAmount  float64  `gorm:"column:transferAmount" json:"transferAmount"`
	Accumulated     *float64 `gorm:"column:accumulated" json:"accumulated"`
	SubAccount      *string  `gorm:"column:subAccount" json:"subAccount"`
	ReferenceCode   *string  `gorm:"column:referenceCode" json:"referenceCode"`
	Description     *string  `gorm:"column:description" json:"description"`
}

func (*SepayTransaction) TableName() string {
	return TableNameSepayTransaction
}
