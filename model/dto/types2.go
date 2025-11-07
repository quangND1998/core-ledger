package dto

import (
	"time"
)

type CustomTime struct {
	*time.Time
}

const layoutNoColonTZ = "2006-01-02T15:04:05-0700"

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}
	t, err := time.Parse(layoutNoColonTZ, s)
	if err != nil {
		return err
	}
	ct.Time = &t
	return nil
}

type LoginResponse struct {
	ExpiresAt CustomTime `json:"expires_at"`
	Token     string     `json:"token"`
}

type FetchBalanceResponse []BalanceItem

type BalanceItem struct {
	AvailableAmount int    `json:"available_amount"`
	Currency        string `json:"currency"`
	PendingAmount   int    `json:"pending_amount"`
	TotalAmount     int    `json:"total_amount"`
}

type KeySet struct {
	ApiKey   string `json:"ApiKey"`
	ClientID string `json:"ClientID"`
}
