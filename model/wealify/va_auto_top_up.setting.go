package model

import "time"

type VaAutoTopUpSettingDataItem struct {
	ID      int64     `json:"id"`
	AddedAt time.Time `json:"added_at"`
}

type BankWhitelistSettingDataItem struct {
	Name        string `json:"name"`
	ShortName   string `json:"short_name"`
	SwiftCode   string `json:"swift_code"`
	CustomerIDs int64  `json:"customer_ids"`
}
