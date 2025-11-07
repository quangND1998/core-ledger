package dto

import (
	model "core-ledger/model/wealify"
	"time"
)

type SQSEvent struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Timestamp int64                  `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

type DetailTopUpEvent struct {
	Amount        float64        `json:"amount"`
	TransactionID string         `json:"transaction_id"`
	Payment       Payment        `json:"payment"`
	Customer      Customer       `json:"customer"`
	Fee           model.FeeField `json:"fee"`
	Remark        string         `json:"remark"`
	CreatedAt     time.Time      `json:"created_at"` // parse từ ISO datetime
}

type Payment struct {
	Detail PaymentDetail `json:"detail"`
}

type PaymentDetail struct {
	BankBin       string `json:"bank_bin"`
	BankCode      string `json:"bank_code"`
	BankName      string `json:"bank_name"`
	AccountName   string `json:"account_name"`
	AccountNumber string `json:"account_number"`
}

type Customer struct {
	ID int64 `json:"id"`
}
