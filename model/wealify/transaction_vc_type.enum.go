package model

type TransactionVcType string

const (
	TransactionVcTypeWithdrawal TransactionVcType = "WITHDRAWAL"
	TransactionVcTypeTopUp      TransactionVcType = "TOP_UP"
	TransactionVcTypePayment    TransactionVcType = "PAYMENT"
	TransactionVcTypeRefund     TransactionVcType = "REFUND"
)

func (v TransactionVcType) String() string {
	return string(v)
}
