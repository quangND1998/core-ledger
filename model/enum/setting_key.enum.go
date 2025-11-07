package enum

type SettingKey string

const (
	SettingKeyVaLimit                   SettingKey = "VA_LIMIT"
	SettingKeyVaFee                     SettingKey = "VA_FEE"
	SettingKeyLarkChatId                SettingKey = "LARK_CHAT_ID"
	SettingKeyInactiveVaRestrictedHour  SettingKey = "INACTIVE_VA_RESTRICTED_HOUR"
	SettingKeyYoobilAutoWithdraw        SettingKey = "YOOBIL_AUTO_WITHDRAW"
	SettingKeyYoobilMapBank             SettingKey = "YOOBIL_MAP_BANK"
	SettingKeyGTelSetting               SettingKey = "G_TEL_SETTING"
	SettingKeyGTelWithdraw              SettingKey = "G_TEL_WITHDRAW"
	SettingKeyGTelMapBank               SettingKey = "G_TEL_MAP_BANK"
	SettingKeyNeoXAutoWithdraw          SettingKey = "NEOX_AUTO_WITHDRAW"
	SettingKeyNeoXMapBank               SettingKey = "NEOX_MAP_BANK"
	SettingKeyVcFee                     SettingKey = "VC_FEE"
	SettingKeyVaEarlyApproveTopUp       SettingKey = "VA_EARLY_APPROVE_TOP_UP"
	SettingKeyAutoProcessToApproveTopUp SettingKey = "auto_PROC_APPR"
	SettingKeyBankWhitelist             SettingKey = "BANK_WHITELIST"
)

func (sk SettingKey) String() string {
	return string(sk)
}
