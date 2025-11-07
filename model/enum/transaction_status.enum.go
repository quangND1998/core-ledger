package enum

type TransactionStatus string

const (
	TransactionStatusUnknown   TransactionStatus = "UNKNOWN"
	TransactionStatusPending   TransactionStatus = "PENDING"
	TransactionStatusProcess   TransactionStatus = "PROCESS"
	TransactionStatusApproved  TransactionStatus = "APPROVED"
	TransactionStatusRejected  TransactionStatus = "REJECTED"
	TransactionStatusWaiting   TransactionStatus = "WAITING"
	TransactionStatusCancelled TransactionStatus = "CANCELLED"
	TransactionStatusExpired   TransactionStatus = "EXPIRED"
	TransactionStatusOnHold    TransactionStatus = "ON_HOLD"
)

func (s TransactionStatus) String() string {
	return string(s)
}

func (s TransactionStatus) Values() []TransactionStatus {
	return []TransactionStatus{
		TransactionStatusPending,
		TransactionStatusProcess,
		TransactionStatusApproved,
		TransactionStatusRejected,
		TransactionStatusWaiting,
		TransactionStatusCancelled,
		TransactionStatusExpired,
		TransactionStatusOnHold,
	}
}
