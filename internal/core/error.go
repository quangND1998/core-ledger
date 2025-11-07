package core

import (
	"errors"
)

var ErrInvalidPlatform = errors.New("invalid platform")

// error code module + service + category + sequence
// example
// + module VA: 02
// + service user: 001
// + category: 01
// + sequence: 000
// error code 0200101000
//
// category :
// + 01 validate
// + 02 business
// + 03 third party

type AppErrorCode string

const (
	ErrCodeUnknown AppErrorCode = "-1"

	ErrCodeWealifyWalletFeatureInactive AppErrorCode = "0100102403"

	ErrCodeVALimitExceed                        AppErrorCode = "0200201001"
	ErrCodeVACardHolder                         AppErrorCode = "0200201002"
	ErrCodeVAFeatureInactive                    AppErrorCode = "0200102403"
	ErrCodeVAPayoutProviderNotEnoughBalance     AppErrorCode = "0200402001"
	ErrCodeVAPayoutProviderCallThirdPartyFailed AppErrorCode = "0200403001"
	ErrCodeVAPayoutCannotSplitTransaction       AppErrorCode = "0200403001"
)

type AppError struct {
	Code        AppErrorCode `json:"code"`
	Scope       string       `json:"scope,omitempty"`
	Message     string       `json:"message"`
	Description string       `json:"description,omitempty"`
}

func (e *AppError) Error() string {
	//if e. != nil {
	//	return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	//}
	return e.Description
}

func (e *AppError) ErrorCode() AppErrorCode {
	return e.Code
}

func (e *AppError) ErrorMessage() string {
	return e.Message
}

var MapCodeToScope = map[AppErrorCode]string{
	ErrCodeVALimitExceed:                        "VA.VA.CREATE.LIMIT_EXCEEDED",
	ErrCodeVACardHolder:                         "VA.CREATE.VALIDATE.CARD_HOLDER",
	ErrCodeVAPayoutProviderNotEnoughBalance:     "VA.TRANSACTION.WITHDRAW.NOT_ENOUGH_BALANCE",
	ErrCodeVAPayoutProviderCallThirdPartyFailed: "VA.TRANSACTION.UNKNOWN_ERROR",
}

var MapCodeToMessage = map[AppErrorCode]string{
	ErrCodeUnknown:                      "Có gì đó bất thường, vui lòng kiểm tra lại",
	ErrCodeWealifyWalletFeatureInactive: "Bạn chưa được active tính năng Wealify Wallet",

	ErrCodeVALimitExceed:                        "Số lượng VA đã đạt giới hạn",
	ErrCodeVACardHolder:                         "Tên VA không hợp lệ",
	ErrCodeVAFeatureInactive:                    "Bạn chưa được active tính năng VA",
	ErrCodeVAPayoutProviderNotEnoughBalance:     "Số dư không đủ",
	ErrCodeVAPayoutProviderCallThirdPartyFailed: "Có gì đó không đúng",
}

var MapCodeToDescription = map[AppErrorCode]string{
	ErrCodeVALimitExceed:                        "Số lượng VA đã đạt giới hạn",
	ErrCodeVACardHolder:                         "Tên VA không hợp lệ",
	ErrCodeVAPayoutProviderNotEnoughBalance:     "Số dư không đủ",
	ErrCodeVAPayoutProviderCallThirdPartyFailed: "Có gì đó không đúng",
}

func NewError(code AppErrorCode, customDescription ...string) *AppError {
	err := &AppError{
		Code:        code,
		Scope:       MapCodeToScope[code],
		Message:     MapCodeToMessage[code],
		Description: MapCodeToDescription[code],
	}
	if len(customDescription) > 0 && customDescription[0] != "" {
		err.Description = customDescription[0]
	}
	return err
}
