package dto

type UserBalanceResponse struct {
	WalletBalance WalletBalance `json:"wallet_balance"`
	CardBalance   CardBalance   `json:"card_balance"`
}

type WalletBalance struct {
	Balance  float64 `json:"balance"`
	MoneyIn  float64 `json:"money_in"`
	MoneyOut float64 `json:"money_out"`
}

type CardBalance struct {
	Total         float64 `json:"total"`
	TotalTopUp    float64 `json:"total_top_up"`
	TotalWithdraw float64 `json:"total_withdraw"`
}

type UserWalletType string

const (
	UserWalletTypeVc   UserWalletType = "VC_CARD"
	UserWalletTypeMain UserWalletType = "MAIN"
)

type UserBalanceHistoryQuery struct {
	Wallet UserWalletType `json:"wallet"`
}

func (u UserWalletType) String() string {
	return string(u)
}
