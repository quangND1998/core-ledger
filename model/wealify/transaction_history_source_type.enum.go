package model

type TransactionHistorySourceType string

const (
	TransactionHistorySourceTypeVaWallet    TransactionHistorySourceType = "VA_WALLET"
	TransactionHistorySourceTypeVcWallet    TransactionHistorySourceType = "VC_WALLET"
	TransactionHistorySourceTypeVirtualCard TransactionHistorySourceType = "VIRTUAL_CARD"
)

func (t TransactionHistorySourceType) String() string {
	return string(t)
}
