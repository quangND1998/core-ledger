package jsonfield

import (
	"core-ledger/model/enum"
	model "core-ledger/model/wealify"
	"fmt"
)

type RateField struct {
	Value float64             `json:"value"`
	Type  model.RateValueType `json:"type"`
}

type SettingDataVALimit struct {
	Gold               int                          `json:"gold"`
	Silver             int                          `json:"silver"`
	Diamond            int                          `json:"diamond"`
	Standard           int                          `json:"standard"`
	CapLimitByCustomer []*VirtualAccountCustomerCap `json:"cap_limit_by_customer,omitempty"`
}
type VirtualAccountCustomerCap struct {
	CustomerID *int64           `json:"customer_id,omitempty"` //in case have customer_id, we apply for this, ignore tier
	Tier       enum.Tier        `json:"tier"`
	Bank       enum.VABankCode  `json:"bank"`
	Provider   model.VAProvider `json:"provider"`
	DailyCap   *int64           `json:"daily_cap"`
	WeeklyCap  *int64           `json:"weekly_cap"`
	MonthlyCap *int64           `json:"monthly_cap"`
	AllTimeCap *int64           `json:"all_time_cap"`
}

func (s *SettingDataVALimit) Validate() error {
	if s.Gold < 0 {
		return fmt.Errorf("field 'gold' can not less than 0")
	}
	if len(s.CapLimitByCustomer) > 0 {
	}
	return nil
}
