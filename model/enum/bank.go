package enum

// BankCode core bank code in Vietnam
type BankCode string

const (
	BankCodeUnknown BankCode = ""
	BankCodeBIDC    BankCode = "BIDC"
	BankCodeBIDV    BankCode = "BIDV"
	BankCodeTCB     BankCode = "TCB"
	BankCodeMB      BankCode = "MB"
	BankCodeMSB     BankCode = "MSB"
	BankCodeKLB     BankCode = "KLB"
)

func (b BankCode) String() string {
	return string(b)
}

type VABankCode string

const (
	VABankCodeBIDV         VABankCode = "BIDV"
	VABankCodeBIDVPremium  VABankCode = "BIDV_PREMIUM"
	VABankCodeBIDVPriority VABankCode = "BIDV_PRIORITY"
	VABankCodeKLB          VABankCode = "KLB"
	VABankCodeMB           VABankCode = "MB"
	VABankCodeMBPriority   VABankCode = "MB_PRIORITY"
	VABankCodeMSB          VABankCode = "MSB"
	VABankCodeMSBElite     VABankCode = "MSB_ELITE"
	VABankCodeTCB          VABankCode = "TCB"
	VABankCodeTCBElite     VABankCode = "TCB_ELITE"
)

func (b VABankCode) String() string {
	return string(b)
}

func (b VABankCode) ToBankCode() BankCode {
	switch b {
	case VABankCodeBIDV, VABankCodeBIDVPremium, VABankCodeBIDVPriority:
		return BankCodeBIDV
	case VABankCodeKLB:
		return BankCodeKLB
	case VABankCodeMB, VABankCodeMBPriority:
		return BankCodeMB
	case VABankCodeMSB, VABankCodeMSBElite:
		return BankCodeMSB
	case VABankCodeTCB, VABankCodeTCBElite:
		return BankCodeTCB
	default:
		return BankCodeUnknown
	}
}
