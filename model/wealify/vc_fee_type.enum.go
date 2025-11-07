package model

type VcFeeType string

const (
	VcFeeTypeTopUpCard      VcFeeType = "TOP_UP_CARD"
	VcFeeTypeIssueCard      VcFeeType = "ISSUE_CARD"
	VcFeeTypeWithdrawalCard VcFeeType = "WITHDRAWAL_CARD"
	VcFeeTypeTopUpWallet    VcFeeType = "TOP_UP_WALLET"
)

func (v VcFeeType) String() string {
	return string(v)
}
