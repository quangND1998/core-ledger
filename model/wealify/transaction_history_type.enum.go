package model

type TransactionHistoryType string

const (
	TransactionHistoryTypeCreated       TransactionHistoryType = "CREATED"
	TransactionHistoryTypeUpdated       TransactionHistoryType = "UPDATED"
	TransactionHistoryTypeReviewed      TransactionHistoryType = "REVIEWED"
	TransactionHistoryTypeHold          TransactionHistoryType = "HOLD"
	TransactionHistoryTypeResolved      TransactionHistoryType = "RESOLVED"
	TransactionHistoryTypeChangeStatus  TransactionHistoryType = "CHANGE_STATUS"
	TransactionHistoryTypeCancelled     TransactionHistoryType = "CANCELLED"
	TransactionHistoryTypeMappedPending TransactionHistoryType = "MAPPED_PENDING"
)

func (t TransactionHistoryType) String() string {
	return string(t)
}
