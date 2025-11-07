package dto

type TwoFactorStatus string

const (
	TFSDisabled TwoFactorStatus = "DISABLED"
	TFSEnabled  TwoFactorStatus = "ENABLED"
)
