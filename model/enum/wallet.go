package enum

type WalletChangeType string

const (
	WalletChangeTypeTopUp      WalletChangeType = "TOP_UP"
	WalletChangeTypeWithdrawal WalletChangeType = "WITHDRAWAL"
	WalletChangeTypeInternal   WalletChangeType = "INTERNAL" // Deprecated
)
