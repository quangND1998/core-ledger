package dto

import "time"

type CardMailInfo struct {
	Name     string  `json:"name"`
	LastFour string  `json:"lastFour"`
	Balance  float64 `json:"balance"`
}

type TransactionMailInfo struct {
	CustomerName    string    `json:"customerName"`
	AmountAfterFee  float64   `json:"amountAfterFee"`
	AmountBeforeFee float64   `json:"amountBeforeFee"`
	Currency        string    `json:"currency"`
	Merchant        string    `json:"merchant"`
	Network         string    `json:"network"`
	FeeValue        float64   `json:"feeValue"`
	Fee             float64   `json:"fee"`
	ReasonRefund    string    `json:"reasonRefund"`
	Time            time.Time `json:"time"`
	Hour            int       `json:"hour"`
	Minute          int       `json:"minute"`
	Date            int       `json:"date"`
	Month           int       `json:"month"`
	Year            int       `json:"year"`
}

type UserMailInfo struct {
	Name string `json:"name"`
}

type CryptoWalletMailInfo struct {
	Ethereum CryptoWalletMailInfoItem `json:"ethereum"`
	Solana   CryptoWalletMailInfoItem `json:"solana"`
	Tron     CryptoWalletMailInfoItem `json:"tron"`
}

type WalletMailInfo struct {
	Balance float64 `json:"balance"`
}

type CryptoWalletMailInfoItem struct {
	Address string `json:"address"`
	Image   string `json:"image"`
	Name    string `json:"name"`
}

// main

type CardCancelMail struct {
	Card CardMailInfo `json:"card"`
}

type CardCreateMail struct {
	Card   CardMailInfo
	Crypto CryptoWalletMailInfo
}

type CardFreezeMail struct {
	Card CardMailInfo `json:"card"`
}

type CardPaymentFailedMail struct {
	Card        CardMailInfo        `json:"card"`
	Transaction TransactionMailInfo `json:"transaction"`
}

type CardPaymentSuccessMail struct {
	Transaction TransactionMailInfo  `json:"transaction"`
	Card        CardMailInfo         `json:"card"`
	Crypto      CryptoWalletMailInfo `json:"crypto"`
}

type CardRefundMail struct {
	Card        CardMailInfo        `json:"card"`
	Transaction TransactionMailInfo `json:"transaction"`
}

type CardTopUpCryptoMail struct {
	User        UserMailInfo         `json:"user"`
	Transaction TransactionMailInfo  `json:"transaction"`
	Card        CardMailInfo         `json:"card"`
	Crypto      CryptoWalletMailInfo `json:"crypto"`
}

type CardTopUpWalletMail struct {
	User        UserMailInfo         `json:"user"`
	Card        CardMailInfo         `json:"card"`
	Transaction TransactionMailInfo  `json:"transaction"`
	Crypto      CryptoWalletMailInfo `json:"crypto"`
}

type CardUnfreezeMail struct {
	Card   CardMailInfo         `json:"card"`
	Crypto CryptoWalletMailInfo `json:"crypto"`
}

type SendOtpMail struct {
	User struct {
		FullName string `json:"fullName"`
	} `json:"user"`
	Otp struct {
		Code string `json:"code"`
	} `json:"otp"`
}

type WalletRefundFromCancelCardMail struct {
	User        UserMailInfo        `json:"user"`
	Transaction TransactionMailInfo `json:"transaction"`
	Wallet      WalletMailInfo      `json:"wallet"`
}

type WalletTopUpCryptoMail struct {
	User        UserMailInfo         `json:"user"`
	Transaction TransactionMailInfo  `json:"transaction"`
	Crypto      CryptoWalletMailInfo `json:"crypto"`
	Wallet      WalletMailInfo       `json:"wallet"`
}
