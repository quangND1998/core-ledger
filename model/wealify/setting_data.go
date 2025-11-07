package model

type VcSettingData struct {
	IssueCard      IssueCard      `json:"ISSUE_CARD"`
	TopUpCard      TopUpCard      `json:"TOP_UP_CARD"`
	TopUpWallet    TopUpWallet    `json:"TOP_UP_WALLET"`
	WithdrawalCard WithdrawalCard `json:"WITHDRAWAL_CARD"`
}

type IssueCard struct {
	Type  FeeType `json:"type"`
	Value float64 `json:"value"`
}

type TopUpCard struct {
	Type  FeeType `json:"type"`
	Value float64 `json:"value"`
}

type TopUpWallet struct {
	Type  FeeType `json:"type"`
	Value float64 `json:"value"`
}

type WithdrawalCard struct {
	Type  FeeType `json:"type"`
	Value float64 `json:"value"`
}
