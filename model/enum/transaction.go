package enum

import (
	"core-ledger/internal/module/constants"
)

type TransactionThirdPartyStatus string

const (
	TransactionThirdPartyStatusCorrect   TransactionThirdPartyStatus = "CORRECT"
	TransactionThirdPartyStatusIncorrect TransactionThirdPartyStatus = "INCORRECT"
)

type TransactionEvent string

const (
	TransactionEventNew     TransactionEvent = "NEW" // just created
	TransactionEventApprove TransactionEvent = "APPROVE"
	TransactionEventReject  TransactionEvent = "REJECT"
)

type TransactionProvider string

const (
	TransactionProviderBank       TransactionProvider = "BANK"
	TransactionProviderAirwallex  TransactionProvider = "AIRWALLEX"
	TransactionProviderWealify    TransactionProvider = "WEALIFY"
	TransactionProviderPingPong   TransactionProvider = "PING_PONG"
	TransactionProviderLianLian   TransactionProvider = "LIAN_LIAN"
	TransactionProviderPayoneer   TransactionProvider = "PAYONEER"
	TransactionProviderWorldFirst TransactionProvider = "World_FIRST"
	TransactionProviderTaZaPay    TransactionProvider = "TAZAPAY"
	TransactionProviderMercury    TransactionProvider = "MERCURY"
	TransactionProviderYoobil     TransactionProvider = "YOOBIL"
	TransactionProviderNeoX       TransactionProvider = "NEOX"
	TransactionProviderGTel       TransactionProvider = "G_TEL"
	TransactionProviderHPay       TransactionProvider = "H_PAY"
	TransactionProvider9Pay       TransactionProvider = "9_PAY"
)

func (tp TransactionProvider) String() string {
	return string(tp)
}

type VaTransactionStatus string

// VA Transaction Status
const (
	VaTransactionStatusFailure    VaTransactionStatus = "FAILURE"
	VaTransactionStatusVerifying  VaTransactionStatus = "VERIFYING"
	VaTransactionStatusWaiting    VaTransactionStatus = "WAITING"
	VaTransactionStatusProcessing VaTransactionStatus = "PROCESSING"
	VaTransactionStatusSuccess    VaTransactionStatus = "SUCCESS"
)

type TransactionProviderInfo map[TransactionProvider]ProviderInfo

var TRANSACTION_PROVIDER_INFO = TransactionProviderInfo{
	TransactionProviderBank: {
		Name:  "Bank transfer",
		Value: TransactionProviderBank.String(),
		Icon:  "",
		Group: constants.PROVIDER_GROUP_BANK,
		TransactionTypes: []string{
			constants.TRANSACTION_TYPE_TOP_UP,
			constants.TRANSACTION_TYPE_WITHDRAWAL,
		},
	},
	TransactionProviderWealify: {
		Name:  "Wealify",
		Value: TransactionProviderWealify.String(),
		Icon:  "",
		Group: constants.PROVIDER_GROUP_WEALIFY,
		TransactionTypes: []string{
			constants.TRANSACTION_TYPE_INTERNAL,
		},
	},
	TransactionProviderPingPong: {
		Name:  "PingPong",
		Value: TransactionProviderPingPong.String(),
		Icon:  "",
		Group: constants.PROVIDER_GROUP_E_WALLET,
		TransactionTypes: []string{
			constants.TRANSACTION_TYPE_TOP_UP,
			constants.TRANSACTION_TYPE_WITHDRAWAL,
		},
	},
	TransactionProviderLianLian: {
		Name:  "LianLian",
		Value: TransactionProviderLianLian.String(),
		Icon:  "",
		Group: constants.PROVIDER_GROUP_E_WALLET,
		TransactionTypes: []string{
			constants.TRANSACTION_TYPE_TOP_UP,
			constants.TRANSACTION_TYPE_WITHDRAWAL,
		},
	},
	TransactionProviderPayoneer: {
		Name:  "Payoneer",
		Value: TransactionProviderPayoneer.String(),
		Icon:  "",
		Group: constants.PROVIDER_GROUP_E_WALLET,
		TransactionTypes: []string{
			constants.TRANSACTION_TYPE_TOP_UP,
			constants.TRANSACTION_TYPE_WITHDRAWAL,
		},
	},
	TransactionProviderYoobil: {
		Name:  "Yoobil",
		Value: TransactionProviderYoobil.String(),
		Icon:  "",
		Group: constants.PROVIDER_GROUP_VA,
		TransactionTypes: []string{
			constants.TRANSACTION_TYPE_TOP_UP,
			constants.TRANSACTION_TYPE_WITHDRAWAL,
		},
	},
	TransactionProviderNeoX: {
		Name:  "NeoX",
		Value: TransactionProviderNeoX.String(),
		Icon:  "",
		Group: constants.PROVIDER_GROUP_VA,
		TransactionTypes: []string{
			constants.TRANSACTION_TYPE_TOP_UP,
			constants.TRANSACTION_TYPE_WITHDRAWAL,
		},
	},
	TransactionProvider9Pay: {
		Name:  "NeoX",
		Value: TransactionProvider9Pay.String(),
		Icon:  "",
		Group: constants.PROVIDER_GROUP_VA,
		TransactionTypes: []string{
			constants.TRANSACTION_TYPE_TOP_UP,
			constants.TRANSACTION_TYPE_WITHDRAWAL,
		},
	},
	TransactionProviderWorldFirst: {
		Name:  "WorldFirst",
		Value: TransactionProviderWorldFirst.String(),
		Icon:  "",
		Group: constants.PROVIDER_GROUP_E_WALLET,
		TransactionTypes: []string{
			constants.TRANSACTION_TYPE_TOP_UP,
		},
	},
	TransactionProviderTaZaPay: {
		Name:  "Tazapay",
		Value: TransactionProviderTaZaPay.String(),
		Icon:  "",
		Group: constants.PROVIDER_GROUP_E_WALLET,
		TransactionTypes: []string{
			constants.TRANSACTION_TYPE_TOP_UP,
		},
	},
	TransactionProviderMercury: {
		Name:  "Mercury",
		Value: TransactionProviderMercury.String(),
		Icon:  "",
		Group: constants.PROVIDER_GROUP_E_WALLET,
		TransactionTypes: []string{
			constants.TRANSACTION_TYPE_TOP_UP,
		},
	},
	TransactionProviderGTel: {
		Name:  "G Tel",
		Value: TransactionProviderGTel.String(),
		Icon:  "",
		Group: constants.PROVIDER_GROUP_VA,
		TransactionTypes: []string{
			constants.TRANSACTION_TYPE_TOP_UP,
			constants.TRANSACTION_TYPE_WITHDRAWAL,
		},
	},
	TransactionProviderHPay: {
		Name:  "H Pay",
		Value: TransactionProviderHPay.String(),
		Icon:  "",
		Group: constants.PROVIDER_GROUP_VA,
		TransactionTypes: []string{
			constants.TRANSACTION_TYPE_TOP_UP,
			constants.TRANSACTION_TYPE_WITHDRAWAL,
		},
	},
}

type ServiceType string

const (
	ServiceTypeWallet ServiceType = "WALLET"
	ServiceTypeVA     ServiceType = "VA"
	ServiceTypeVC     ServiceType = "VC"
)
