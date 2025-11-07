package model

import (
	"core-ledger/model/enum"
	"strconv"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const TableNameTransaction = "airbyte_raw.wealify_transactions"

// Transaction mapped from table <transactions>
type Transaction struct {
	CreatedAt               time.Time                        `gorm:"column:created_at;not null;autoCreateTime" json:"created_at"`
	UpdatedAt               time.Time                        `gorm:"column:updated_at;not null;autoCreateTime" json:"updated_at"`
	Status                  bool                             `gorm:"column:status;not null;default:1" json:"status"`
	IsDeleted               bool                             `gorm:"column:is_deleted;not null" json:"is_deleted"`
	ID                      string                           `gorm:"column:id;primaryKey" json:"id"`
	TransactionID           string                           `gorm:"column:transaction_id;not null" json:"transaction_id"` // TransactionCode
	Provider                enum.TransactionProvider         `gorm:"column:provider;not null;default:BANK" json:"provider"`
	ProviderType            string                           `gorm:"column:provider_type;not null;default:INDIVIDUAL" json:"provider_type"`
	Fee                     *datatypes.JSON                  `gorm:"column:fee" json:"fee"`
	Rate                    *datatypes.JSON                  `gorm:"column:rate" json:"rate"`
	SystemRate              *datatypes.JSON                  `gorm:"column:system_rate" json:"system_rate"`
	CustomerInfo            *datatypes.JSON                  `gorm:"column:customer_info" json:"customer_info"`
	Amount                  float64                          `gorm:"column:amount;not null" json:"amount"`
	PaymentAmount           float64                          `gorm:"column:payment_amount;not null" json:"payment_amount"`
	Remark                  string                           `gorm:"column:remark" json:"remark"`
	RefMessage              string                           `gorm:"column:ref_message;not null" json:"ref_message"`
	TransactionType         enum.TransactionType             `gorm:"column:transaction_type;not null" json:"transaction_type"`
	TransactionStatus       enum.TransactionStatus           `gorm:"column:transaction_status" json:"transaction_status"`
	Note                    string                           `gorm:"column:note" json:"note"`
	Payment                 *datatypes.JSON                  `gorm:"column:payment" json:"payment"`
	CreatedByEmployee       bool                             `gorm:"column:created_by_employee;not null" json:"created_by_employee"`
	VaTransactionStatus     *enum.VaTransactionStatus        `gorm:"column:va_transaction_status" json:"va_transaction_status"`
	ThirdPartyStatus        enum.TransactionThirdPartyStatus `gorm:"column:third_party_status;default:'CORRECT'" json:"third_party_status"`
	TransactionVcStatus     string                           `gorm:"column:transaction_vc_status;not null" json:"transaction_vc_status"`
	TransactionVcType       TransactionVcType                `gorm:"column:transaction_vc_type" json:"transaction_vc_type"`
	IsVcTransaction         bool                             `gorm:"column:is_vc_transaction;not null" json:"is_vc_transaction"`
	IsIssue                 bool                             `gorm:"column:is_issue;not null" json:"is_issue"`
	VcDetailTransactionType VcDetailTransactionType          `gorm:"column:vc_detail_transaction_type;not null" json:"vc_detail_transaction_type"`
	EffectiveTier           enum.Tier                        `gorm:"column:effective_tier" json:"effective_tier"`
	QrImageUrl              *string                          `gorm:"column:qr_image_url" json:"qr_image_url"`

	CustomerID          int64   `gorm:"column:customer_id" json:"customer_id"`
	CurrencySymbol      string  `gorm:"column:currency_symbol" json:"currency_symbol"`
	CryptoWalletID      *string `gorm:"column:crypto_wallet_id" json:"crypto_wallet_id"`
	ReceiverID          *int64  `gorm:"column:receiver_id" json:"receiver_id"`
	ReceivedWalletID    *string `gorm:"column:received_wallet_id" json:"received_wallet_id"`
	SentWalletID        *string `gorm:"column:sent_wallet_id" json:"sent_wallet_id"`
	SystemPaymentID     *string `gorm:"column:system_payment_id" json:"system_payment_id"`
	TransactionLinkedID *string `gorm:"column:transaction_linked_id" json:"transaction_linked_id"`
	VirtualAccountID    *int64  `gorm:"column:virtual_account_id" json:"virtual_account_id"`
	VirtualCardID       *string `gorm:"column:virtual_card_id" json:"virtual_card_id"`

	Histories          []*TransactionHistorie `gorm:"foreignKey:TransactionID;references:ID" json:"histories"`
	Currency           *Currency              `gorm:"foreignKey:CurrencySymbol;references:Symbol" json:"currency"`
	Customer           *Customer              `gorm:"foreignKey:CustomerID;references:ID" json:"customer"`
	CryptoWallet       *CryptoWallet          `gorm:"foreignKey:CryptoWalletID" json:"crypto_wallet"`
	ConfirmTransaction *ConfirmTransaction    `gorm:"foreignKey:TransactionID;references:ID" json:"confirm_transaction"`
	VirtualAccount     *VirtualAccount        `gorm:"foreignKey:VirtualAccountID;references:ID" json:"virtual_account"`
	VirtualCard        *VirtualCard           `gorm:"foreignKey:VirtualCardID" json:"virtual_card"`
	TransactionLinked  *Transaction           `gorm:"foreignKey:ID;references:TransactionLinkedID" json:"transaction_linked"`
	ReceivedWallet     *Wallet                `gorm:"foreignKey:ReceivedWalletID;references:ID" json:"received_wallet"`
	SentWallet         *Wallet                `gorm:"foreignKey:SentWalletID;references:ID" json:"sent_wallet"`
	SystemPayment      *SystemPayment         `gorm:"foreignKey:SystemPaymentID;references:ID" json:"system_payment"`
	SubTransactions    []*Transaction         `gorm:"foreignKey:TransactionLinkedID;references:ID" json:"sub_transactions"`
}

type TransactionAutoProcessToApprove struct {
	KybStatus         int     `gorm:"column:kyb_status" json:"kyb_status"`
	Amount            float64 `gorm:"column:amount" json:"amount"`
	TransactionID     string  `gorm:"column:id" json:"id"`
	TransactionStatus string  `gorm:"column:transaction_status" json:"transaction_status"`
	PlatformID        string  `gorm:"column:platform_id" json:"platform_id"`
	Remark            string  `gorm:"column:remark" json:"remark"`
}

// TableName Transaction's table name
func (*Transaction) TableName() string {
	return TableNameTransaction
}

func (t *Transaction) BeforeSave(tx *gorm.DB) (err error) {
	if t.ID == "" {
		t.ID = uuid.NewString()
	}
	return
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == "" {
		t.ID = uuid.NewString()
	}
	return
}

func (t *Transaction) BeforeUpdate(tx *gorm.DB) (err error) {
	t.UpdatedAt = time.Now()
	return
}

func (t *Transaction) IsVaTopUp() bool {
	if t.ReceivedWallet != nil {
		if t.ReceivedWallet.Type == WalletTypeVA {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func (t *Transaction) IsMainTopUp() bool {
	if t.ReceivedWallet != nil {
		if t.ReceivedWallet.Type == WalletTypeMain {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func (t *Transaction) IsVaWithdraw() bool {
	if t.SentWalletID != nil {
		if t.SentWallet.Type == WalletTypeVA {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func (t *Transaction) IsMainWithdraw() bool {
	if t.SentWalletID != nil {
		if t.SentWallet.Type == WalletTypeMain {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func (t *Transaction) IsVaTransaction() bool {
	if t.SentWalletID != nil {
		if t.SentWallet.Type == WalletTypeVA {
			return true
		} else {
			return false
		}
	} else if t.ReceivedWalletID != nil {
		if t.ReceivedWallet.Type == WalletTypeVA {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func (t *Transaction) AmountStr() string {
	return humanize.Comma(int64(t.Amount))
}

func (t *Transaction) PaymentAmountStr() string {
	return humanize.Comma(int64(t.PaymentAmount))
}

func (t *Transaction) TopUpType() string {
	if t.QrImageUrl == nil {
		return "NORMAL"
	}
	return "QR"
}

type TransactionView struct {
	Transaction         `json:"-"`
	CreatedAt           time.Time                        `json:"created_at"`
	Status              bool                             `json:"status"`
	ID                  string                           `json:"id"`
	TransactionID       string                           `json:"transaction_id"`
	ProviderType        string                           `json:"provider_type"`
	Fee                 *FeeField                        `json:"fee"`
	Rate                *RateField                       `json:"rate"`
	Amount              float64                          `json:"amount"`
	PaymentAmount       float64                          `json:"payment_amount"`
	Remark              string                           `json:"remark"`
	RefMessage          string                           `json:"ref_message"`
	TransactionType     enum.TransactionType             `json:"transaction_type"`
	TransactionStatus   enum.TransactionStatus           `json:"transaction_status"`
	Note                string                           `json:"note"`
	CreatedByEmployee   bool                             `json:"created_by_employee"`
	VaTransactionStatus *enum.VaTransactionStatus        `json:"va_transaction_status"`
	ThirdPartyStatus    enum.TransactionThirdPartyStatus `json:"third_party_status"`
	QrImageUrl          *string                          `json:"qr_image_url"`
	VirtualAccountID    *int64                           `json:"virtual_account_id"`
	Provider            *enum.TransactionProvider        `json:"-"`

	Payment        *datatypes.JSON `json:"payment"`
	ProviderInfo   *datatypes.JSON `json:"provider"`
	Currency       *datatypes.JSON `json:"currency"`
	Customer       *datatypes.JSON `json:"customer"`
	VirtualAccount *datatypes.JSON `json:"virtual_account"`
	ReceivedWallet *datatypes.JSON `json:"received_wallet"`
	SentWallet     *datatypes.JSON `json:"sent_wallet"`

	// TODO:
	SubTransactions    *datatypes.JSON `json:"sub_transactions"`
	AmountStr          string          `json:"amount_str"`
	PaymentAmountStr   string          `json:"payment_amount_str"`
	ReceivedAmountStr  string          `json:"received_amount_str"`
	ChangedAmountStr   string          `json:"changed_amount_str"`
	HaveSubTransaction bool            `json:"have_sub_transaction"`
}

func (t *TransactionView) ReceivedAmount() float64 {
	switch t.TransactionType {
	case enum.TransactionTypeInternal:
		if t.Fee.Type == FeeTypePercent {
			return -t.Amount * (1 + t.Fee.Amount)
		}
		return -(t.Amount + t.Fee.Amount)

	case enum.TransactionTypeTopUp:
		if t.Fee.Type == FeeTypePercent {
			return t.Amount * (1 - t.Fee.Amount) * t.Rate.Amount
		}
		return (t.Amount - t.Fee.Amount) * t.Rate.Amount

	case enum.TransactionTypeWithdrawal:
		if t.Fee.Type == FeeTypePercent {
			return (t.Amount / t.Rate.Amount) * (1 - t.Fee.Amount)
		}
		return t.Amount/t.Rate.Amount - t.Fee.Amount
	case enum.TransactionTypeAdjustment:
		if t.SentWallet != nil && string(*t.SentWallet) == "{}" {
			return t.Amount
		} else if t.ReceivedWallet != nil && string(*t.ReceivedWallet) == "{}" {
			return -t.Amount
		}
		return t.Amount
	}
	return 0
}
func (t *TransactionView) GetReceivedAmountStr() string {
	return strconv.FormatFloat(t.ReceivedAmount(), 'f', -1, 64)
}
