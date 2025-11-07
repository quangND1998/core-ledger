package model

import (
	"time"
)

const TableNameCampaignReferral = "campaign-referral"

// CampaignReferral mapped from table <campaign-referral>
type CampaignReferral struct {
	CurrentRank  int64     `gorm:"column:current_rank;not null" json:"current_rank"`
	CustomerID   int32     `gorm:"column:customer_id;not null" json:"customer_id"`
	FullName     string    `gorm:"column:full_name;not null" json:"full_name"`
	PhoneNumber  string    `gorm:"column:phone_number" json:"phone_number"`
	TotalInvite  int64     `gorm:"column:total_invite;not null" json:"total_invite"`
	Flux         string    `gorm:"column:flux;not null" json:"flux"`
	KldUpdatedAt time.Time `gorm:"column:kld_updated_at;default:0000-00-00 00:00:00.000000" json:"kld_updated_at"`
}

// TableName CampaignReferral's table name
func (*CampaignReferral) TableName() string {
	return TableNameCampaignReferral
}
