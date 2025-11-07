package dto

type VerifyOtpRequest struct {
	Code string `json:"code"`
}

type SendOtpRequest struct {
	Code         string `json:"code"`
	ExpireMinute int    `json:"expire_minute"`
}
