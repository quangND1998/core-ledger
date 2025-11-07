package enum

type Provider string

const (
	ProviderYoobil Provider = "YOOBIL"
	ProviderNeoX   Provider = "NEOX"
	ProviderGTel   Provider = "G_TEL"
	ProviderHPay   Provider = "H_PAY"
)

type SystemPaymentProvider string

const (
	SystemProviderBank      SystemPaymentProvider = "BANK"
	SystemProviderAirwallex SystemPaymentProvider = "AIRWALLEX"
	SystemProviderPayoneer  SystemPaymentProvider = "PAYONEER"
	SystemProviderPingPong  SystemPaymentProvider = "PING_PONG"
	SystemProviderLianLian  SystemPaymentProvider = "LIAN_LIAN"
	SystemProviderTazapay   SystemPaymentProvider = "TAZAPAY"
	SystemProviderYoobil    SystemPaymentProvider = "YOOBIL"
	SystemProviderNeoX      SystemPaymentProvider = "NEOX"
	SystemProvider9Pay      SystemPaymentProvider = "9_PAY"
	SystemProviderGTel      SystemPaymentProvider = "G_TEL"
	SystemProviderHPay      SystemPaymentProvider = "H_PAY"
)

type ProviderInfo struct {
	Name             string   `json:"name"`
	Value            string   `json:"value"`
	Icon             string   `json:"icon"`
	Group            string   `json:"group"`
	TransactionTypes []string `json:"transaction_types"`
}
