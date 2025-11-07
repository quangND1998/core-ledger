package constants

type VaProvider string

const (
	VaProviderYoobil VaProvider = "YOOBIL"
	VaProviderNeoX   VaProvider = "NEOX"
	VaProviderGTel   VaProvider = "G_TEl"
	VaProviderHPay   VaProvider = "H_PAY"
)

func (p VaProvider) String() string {
	return string(p)
}
