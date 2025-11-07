package dto

import (
	"time"
)

type ListFeesRequest struct {
	Tier            string `form:"tier"`
	Provider        string `form:"provider"`
	ProviderType    string `form:"provider_type"`
	TransactionType string `form:"transaction_type"`
	CurrencySymbol  string `form:"currency_symbol"`
	Page            int    `form:"page"`
	Limit           int    `form:"limit"`
}
type ListFeesResponse struct {
	Total int64              `json:"total"`
	Page  int                `json:"page"`
	Limit int                `json:"limit"`
	Items []*FeeResponseItem `json:"items"`
}
type FeeResponseItem struct {
	CreatedAt       time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(6)" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP(6)" json:"updated_at"`
	Status          int32     `gorm:"column:status;not null;default:1" json:"status"`
	ID              string    `gorm:"column:id;primaryKey" json:"id"`
	AccountType     string    `gorm:"column:account_type" json:"account_type"`
	AccountLevel    string    `gorm:"column:account_level" json:"account_level"`
	Provider        string    `gorm:"column:provider;not null;default:BANK" json:"provider"`
	ProviderType    string    `gorm:"column:provider_type;not null;default:INDIVIDUAL" json:"provider_type"`
	TransactionType string    `gorm:"column:transaction_type;not null" json:"transaction_type"`
	Description     string    `gorm:"column:description" json:"description"`
	CurrencySymbol  string    `gorm:"column:currency_symbol" json:"currency_symbol"`
	Tier            string    `gorm:"column:tier" json:"tier"`
}

type CreateFeeRequest struct {
	Tier            string              `json:"tier" validate:"omitempty,oneof=STANDARD SILVER GOLD DIAMOND"`
	CurrencySymbol  string              `json:"currency_symbol" validate:"required,max=10"`
	Provider        string              `json:"provider" validate:"required,oneof=BANK WEALIFY PING_PONG LIAN_LIAN PAYONEER WORLD_FIRST TAZAPAY MERCURY YOOBIL NEOX G_TEL"`
	ProviderType    string              `json:"provider_type" validate:"required,oneof=BUSINESS INDIVIDUAL"`
	CustomerIDs     []string            `json:"customer_ids" validate:"dive,required,uuid4"`
	Ranges          []RangesItemRequest `json:"ranges" validate:"required,dive"`
	TransactionType string              `json:"transaction_type" validate:"required,oneof=TOP_UP WITHDRAWAL INTERNAL"`
	Description     string              `json:"description"`
}
type RangesItemRequest struct {
	Min   float64 `json:"min" validate:"required,gte=0"`
	Max   float64 `json:"max" validate:"required,gte=0"`
	Value float64 `json:"value" validate:"required,gte=0"`
	Type  string  `json:"type" validate:"required,oneof=FIXED PERCENT"`
}
