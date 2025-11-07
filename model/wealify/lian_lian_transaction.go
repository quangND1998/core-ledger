package model

const TableNameLianLianTransaction = "lian-lian-transactions"

type LianLianTransaction struct {
	ID              string                    `gorm:"column:id;primaryKey" json:"id"`
	RequestID       string                    `gorm:"column:request_id" json:"request_id"`
	BusinessOrderID string                    `gorm:"column:business_order_id" json:"business_order_id"`
	Amount          string                    `gorm:"column:amount" json:"amount"`
	Currency        string                    `gorm:"column:currency" json:"currency"`
	Name            string                    `gorm:"column:name" json:"name"`
	Type            LianLianTransactionType   `gorm:"column:type" json:"type"`
	UserID          string                    `gorm:"column:user_id" json:"user_id"`
	Status          LianLianTransactionStatus `gorm:"column:status" json:"status"`
	CreateTime      int64                     `gorm:"column:create_time" json:"create_time"`
}

// TableName LianLianTransaction's table name
func (*LianLianTransaction) TableName() string {
	return TableNameLianLianTransaction
}

type LianLianTransactionType string
type LianLianTransactionStatus string

const (
	LianLianTransactionTypePayout          LianLianTransactionType = "PAYOUT"
	LianLianTransactionTypeReceipt         LianLianTransactionType = "RECEIPT"
	LianLianTransactionTypeConversion      LianLianTransactionType = "CONVERSION"
	LianLianTransactionTypeWithdrawal      LianLianTransactionType = "WITHDRAWAL"
	LianLianTransactionTypeAddFunds        LianLianTransactionType = "ADD_FUNDS"
	LianLianTransactionTypeReceiving       LianLianTransactionType = "RECEIVING"
	LianLianTransactionTypeRefund          LianLianTransactionType = "REFUND"
	LianLianTransactionTypeCardPayout      LianLianTransactionType = "CARD_PAYOUT"
	LianLianTransactionTypePromoCommission LianLianTransactionType = "PROMO_COMMISSION"
	LianLianTransactionTypeStandingPayment LianLianTransactionType = "STANDING_PAYMENT"

	LianLianTransactionStatusProcessing LianLianTransactionStatus = "PROCESSING"
	LianLianTransactionStatusCompleted  LianLianTransactionStatus = "COMPLETED"
	LianLianTransactionStatusCanceled   LianLianTransactionStatus = "CANCELED"
	LianLianTransactionStatusRefunding  LianLianTransactionStatus = "REFUNDING"
	LianLianTransactionStatusRefunded   LianLianTransactionStatus = "REFUNDED"
	LianLianTransactionStatusFailed     LianLianTransactionStatus = "FAILED"
)
