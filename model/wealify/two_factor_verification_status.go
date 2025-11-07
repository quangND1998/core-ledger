package model

type TwoFactorVerificationStatus string

const (
	TwoFactorVerificationStatusUnverified TwoFactorVerificationStatus = "UNVERIFIED"
	TwoFactorVerificationStatusVerified   TwoFactorVerificationStatus = "VERIFIED"
)
