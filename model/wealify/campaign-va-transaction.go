package model

import (
	"time"
)

const TableNameCampaignVaTransaction = "campaign-va-transaction"

// CampaignVaTransaction mapped from table <campaign-va-transaction>
type CampaignVaTransaction struct {
	CurrentRank          int64     `gorm:"column:current_rank;not null" json:"current_rank"`
	CustomerID           int32     `gorm:"column:customer_id;not null" json:"customer_id"`
	FullName             string    `gorm:"column:full_name;not null" json:"full_name"`
	PhoneNumber          string    `gorm:"column:phone_number" json:"phone_number"`
	Volume               float64   `gorm:"column:volume" json:"volume"`
	Flux                 string    `gorm:"column:flux;not null" json:"flux"`
	TransactionUpdatedAt time.Time `gorm:"column:transaction_updated_at;default:0000-00-00 00:00:00.000000" json:"transaction_updated_at"`
}

// TableName CampaignVaTransaction's table name
func (*CampaignVaTransaction) TableName() string {
	return TableNameCampaignVaTransaction
}
