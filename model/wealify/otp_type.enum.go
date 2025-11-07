package model

type OtpType string

const (
	OtpTypeMail OtpType = "MAIL"
)

func (otp OtpType) String() string {
	return string(otp)
}
