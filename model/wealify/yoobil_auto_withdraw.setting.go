package model

type YoobilAutoWithdrawSetting struct {
	Vndw      bool `json:"vndw"`
	Vndy      bool `json:"vndy"`
	Active    bool `json:"active"`
	Withdraw  int  `json:"withdraw"`
	Withdrawn int  `json:"withdrawn"`
}
