package model

type VcDetailTransactionType string

const (
	VcDetailTransactionTypeCardTopUpCrypto     VcDetailTransactionType = "CARD_TOP_UP_CRYPTO" // top up card, receive from crypto wallet
	VcDetailTransactionTypeCardTopUp           VcDetailTransactionType = "CARD_TOP_UP"        // top up card, withdraw from vc wallet
	VcDetailTransactionTypeCardIssueTopUp      VcDetailTransactionType = "CARD_ISSUE_TOP_UP"  // top up card, first tran when issue card, withdraw from vc wallet
	VcDetailTransactionTypeCardWithDraw        VcDetailTransactionType = "CARD_WITHDRAW"      // withdraw card, cancel card, refund to wallet
	VcDetailTransactionTypeCardPayment         VcDetailTransactionType = "CARD_PAYMENT"
	VcDetailTransactionTypeCardRefund          VcDetailTransactionType = "CARD_REFUND"           // withdraw card, aspire spend
	VcDetailTransactionTypeWalletTopUp         VcDetailTransactionType = "WALLET_TOP_UP"         // top up wallet, receive from crypto
	VcDetailTransactionTypeWalletWithdraw      VcDetailTransactionType = "WALLET_WITHDRAW"       // withdraw wallet, top up to card
	VcDetailTransactionTypeWalletWithdrawBank  VcDetailTransactionType = "WALLET_WITHDRAW_BANK"  // withdraw wallet, transfer to bank account
	VcDetailTransactionTypeWalletIssueWithdraw VcDetailTransactionType = "WALLET_ISSUE_WITHDRAW" // withdraw wallet, to issue then top up card
	VcDetailTransactionTypeWalletRefund        VcDetailTransactionType = "WALLET_REFUND"         // top up wallet, cancel card, with draw from card
	// top up wallet, cancel card, with draw from card
)

func (v VcDetailTransactionType) String() string {
	return string(v)
}
