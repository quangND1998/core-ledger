package dto

import (
	model "core-ledger/model/wealify"
	"core-ledger/pkg/constants"
	"time"

	"gorm.io/datatypes"
)

type CreateTransactionRequest struct {
	BankCode            string                  `json:"bank_code"`
	BankAccountName     string                  `json:"bank_account_name"`
	BankAccountNumber   string                  `json:"bank_account_number"`
	CryptoWalletAddress string                  `json:"crypto_wallet_address"`
	CryptoWalletNetwork string                  `json:"crypto_wallet_network"`
	Type                model.TransactionVcType `json:"type"`
}

type CreateTransactionResponse struct {
}

type ListTransactionFilter struct {
	PaginationDto
	TimeRangeFilter
	Keyword                         string   `json:"keyword" form:"keyword"`
	TransactionTypes                []string `json:"transaction_type[]" form:"transaction_type[]"`
	TransactionStatuses             []string `json:"transaction_status[]" form:"transaction_status[]"`
	VcDetailTransactionTypes        []string `json:"vc_detail_transaction_types" form:"vc_detail_transaction_types"`
	ExcludeVcDetailTransactionTypes []string `json:"exclude_vc_detail_transaction_types" form:"exclude_vc_detail_transaction_types"`
}

type ListTransactionResponse struct {
	Items      []*TransactionListItem `json:"items"`
	TotalItems int64                  `json:"total_items"`
	TotalPage  int64                  `json:"total_page"`
	Limit      int64                  `json:"limit"`
	Page       int64                  `json:"page"`
}

type TransactionListItem struct {
	CreatedAt               time.Time                     `json:"created_at"`
	UpdatedAt               time.Time                     `json:"updated_at"`
	Id                      string                        `json:"id"`
	TransactionId           *string                       `json:"transaction_id"`
	Fee                     *datatypes.JSON               `json:"fee"`
	Rate                    *datatypes.JSON               `json:"rate"`
	Amount                  float64                       `json:"amount"`
	TransactionVcStatus     string                        `json:"transaction_vc_status"`
	TransactionVcType       model.TransactionVcType       `json:"transaction_vc_type"`
	IsIssue                 bool                          `json:"is_issue"`
	VcDetailTransactionType model.VcDetailTransactionType `json:"vc_detail_transaction_type"`
	IsVcTransaction         bool                          `json:"is_vc_transaction"`
	VirtualCardId           *string                       `json:"virtual_card_id"`
	VirtualCard             *model.VirtualCard            `json:"virtual_card"`
	Currency                struct {
		Symbol string `json:"symbol"`
	} `json:"currency"`
	ConfirmTransaction  *model.ConfirmTransaction `json:"confirm_transaction"`
	CryptoWalletId      *string                   `json:"crypto_wallet_id"`
	CryptoWallet        *model.CryptoWallet       `json:"crypto_wallet"`
	TransactionLinkedId *string                   `json:"transaction_linked_id"`
	TransactionLinked   *TransactionListItem      `json:"transaction_linked"`
	ReceivedAmount      float64                   `json:"received_amount"`
	Source              TransactionSource         `json:"source"`
}

type TransactionSource string

const (
	TransactionSourceCard   TransactionSource = "CARD"
	TransactionSourceWallet TransactionSource = "WALLET"
)

type VirtualAccountTopUpDTO struct {
	AccountNo      string
	Provider       constants.VaProvider
	Amount         float64
	PurchaseAmount float64
	Remark         *string
	Currency       string
	TradeNo        *string
}
