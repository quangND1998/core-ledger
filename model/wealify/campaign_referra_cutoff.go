package model

import (
	"time"
)

const TableNameCampaignReferralCutoff = "campaign_referral_cutoff"

type CampaignReferralCutoff struct {
	CurrentRank  int       `gorm:"column:current_rank"   json:"current_rank"`
	FullName     string    `gorm:"column:full_name"      json:"full_name"`
	PhoneNumber  string    `gorm:"column:phone_number"   json:"phone_number"`
	TotalInvite  int       `gorm:"column:total_invite"   json:"total_invite"`
	CustomerID   uint      `gorm:"column:customer_id"    json:"customer_id"`
	Flux         string    `gorm:"column:flux"           json:"flux"`
	KLDUpdatedAt time.Time `gorm:"column:kld_updated_at" json:"kld_updated_at"`
}

func (*CampaignReferralCutoff) TableName() string {
	return TableNameCampaignReferralCutoff
}
