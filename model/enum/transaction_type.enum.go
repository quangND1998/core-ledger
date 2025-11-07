package enum

type TransactionType string

const (
	TransactionTypeWithdrawal TransactionType = "WITHDRAWAL"
	TransactionTypeTopUp      TransactionType = "TOP_UP"
	TransactionTypeInternal   TransactionType = "INTERNAL"
	TransactionTypeAdjustment TransactionType = "ADJUSTMENT"
)

func (v TransactionType) String() string {
	return string(v)
}
