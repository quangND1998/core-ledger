package dto

import (
	"core-ledger/model/enum"
	model "core-ledger/model/wealify"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type ListCustomerFilter struct {
	BasePaginationQuery
	IDs         []int64            `form:"ids"`
	FullName    string             `form:"full_name" binding:"required"`
	Email       string             `form:"email" binding:"required"`
	PhoneNumber *string            `form:"phone_number"`
	CountryCode string             `form:"country_code" binding:"required"`
	Type        *enum.CustomerType `form:"type"`
	Tier        *enum.Tier         `form:"tier"`
}

type AutoApproveTopUpCustomerFilter struct {
	Keyword *string            `form:"keyword,omitempty"`
	Type    *enum.CustomerType `form:"type,omitempty"`
	Tier    *enum.Tier         `form:"tier,omitempty"`
	AddedAt *TimeRangeFilter   `form:"added_at"`
	Limit   *int64             `form:"limit,omitempty"`
	Page    *int64             `form:"page,omitempty"`
}

type AutoApproveTopUpCustomer struct {
	model.Customer
	AddedAt time.Time `json:"added_at,omitempty"`
}

type AvailableBankItem struct {
	ID                string `json:"id"`
	InternationalName string `json:"international_name"`
	Name              string `json:"name"`
	ShortName         string `json:"short_name"`
	Code              string `json:"code"`
	SwiftCode         string `json:"swift_code"`
}

type CMSGetCustomerInfoRequest struct {
	Keyword           *string  `form:"keyword"`
	AccountLevels     []int    `form:"account_levels[]"`
	Status            *string  `form:"status"`
	Tier              *string  `form:"tier"`
	VAEnable          *bool    `form:"va_enable"`
	InHouse           *int     `form:"in_house"`
	Page              int      `form:"page"`
	Limit             int      `form:"limit"`
	CustomerIDs       []int64  `form:"customer_ids[]"`
	TransactionRanges []string `form:"transaction_ranges[]"`
	From              string   `json:"-"`
	To                string   `json:"-"`
}
type ExportCMSCustomerInfoRequest struct {
	CMSGetCustomerInfoRequest
	Lang string `form:"lang"`
}
type ExportCMSCustomerInfoResponse struct {
	FileName string `json:"file_name"`
	URL      string `json:"url"`
}

func (r *CMSGetCustomerInfoRequest) Decode(c *gin.Context) error {
	if err := c.ShouldBindQuery(r); err != nil {
		return fmt.Errorf("binding err: %v", err)
	}
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.Page > 1000_000 {
		r.Page = 1000000
	}
	if r.Limit <= 0 {
		r.Limit = 10
	}
	if r.Limit > 1000 {
		r.Limit = 1000
	}
	from, to, err := toTimeRanges(r.TransactionRanges)
	if err != nil {
		return err
	}
	r.From = from
	r.To = to

	return nil
}
func toTimeRanges(ranges []string) (from, to string, err error) {
	if len(ranges) > 0 {
		from = ranges[0]
		if _, err := time.Parse(time.DateOnly, from); err != nil {
			return "", "", fmt.Errorf("parse transaction ranges failed: %v", err)
		}
	}
	if len(ranges) > 1 {
		to = ranges[1]
		if _, err := time.Parse(time.DateOnly, to); err != nil {
			return "", "", fmt.Errorf("parse transaction ranges failed: %v", err)
		}
	}
	return
}

type CMSGetCustomerInfoDataItem struct {
	ID                      int64                    `json:"id"`
	FullName                string                   `json:"full_name"`
	Email                   string                   `json:"email"`
	AccountLevel            AccountLevelInfo         `json:"account_level"`
	Tier                    TierInfo                 `json:"tier"`
	Wallets                 []WalletBalanceInfo      `json:"wallets"`
	VAEnable                bool                     `json:"va_enable"`
	VirtualAccountStatistic VirtualAccountStatistic  `json:"virtual_account_statistic"`
	VirtualAccount          VirtualAccountInfo       `json:"virtual_account"`
	Balance                 *CustomerBalance         `json:"balance"`
	VARealBalance           float64                  `json:"va_real_balance"`
	VARealBalanceStr        string                   `json:"va_real_balance_str"`
	PendingAmount           map[string]PendingAmount `json:"pending_amount"`
}
type TierInfo struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Icon  string `json:"icon"`
}
type AccountLevelInfo struct {
	Name  string             `json:"name"`
	Value model.AccountLevel `json:"value"`
	Icon  string             `json:"icon"`
}
type AccountTypeInfo struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Icon  string `json:"icon"`
}

type StatsByPlatforms struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	TotalCreated    int     `json:"total_created"`
	TotalActive     int     `json:"total_active"`
	TotalInactive   int     `json:"total_inactive"`
	TotalRestricted int     `json:"total_restricted"`
	StandardFee     float64 `json:"standard_fee"`
	SilverFee       float64 `json:"silver_fee"`
	GoldFee         float64 `json:"gold_fee"`
	DiamondFee      float64 `json:"diamond_fee"`
	Type            string  `json:"type"`
	Status          string  `json:"status"`
	CurrentFee      float64 `json:"current_fee"`
}

type VirtualAccountStatistic struct {
	Platforms       []*StatsByPlatforms `json:"platforms,omitempty"`
	TotalCreated    int                 `json:"total_created"`
	TotalActive     int                 `json:"total_active"`
	TotalInactive   int                 `json:"total_inactive"`
	TotalRestricted int                 `json:"total_restricted"`
}
type VirtualAccountInfo struct {
	TotalFee             float64 `json:"total_fee"`
	TotalWithdraw        float64 `json:"total_withdraw"`
	TotalWithdrawSuccess float64 `json:"total_withdraw_success"`
	TotalWithdrawPending float64 `json:"total_withdraw_pending"`
	TotalTopUp           float64 `json:"total_top_up"`
	TotalTopUpPending    float64 `json:"total_top_up_pending"`
	TotalReceived        float64 `json:"total_received"`
	Debt                 float64 `json:"debt"`
	Limit                int     `json:"limit"`
}

type CMSGetCustomerInfoResponse []*CMSGetCustomerInfoDataItem

type WalletBalanceInfo struct {
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
	Type     string  `json:"type"`
}

type CustomerBalance struct {
	WalletType        model.WalletType `json:"wallet_type"`
	TotalTopUpSuccess float64          `json:"total_top_up_success"`
	TotalBalance      float64          `json:"total_balance"`
	PendingBalance    float64          `json:"pending_balance"`
	MonthlyTopUp      float64          `json:"monthly_top_up"`
	MonthlyWithdraw   float64          `json:"monthly_withdraw"`
	ReadyBalance      float64          `json:"ready_balance"`
	ReadyBalanceStr   string           `json:"ready_balance_str"`
	UpdatedAt         time.Time        `json:"updated_at"`
	Currency          struct {
		Symbol string `json:"symbol"`
	} `json:"currency"`
}

type GetCustomerBalanceResponse struct {
	Items []CustomerBalance `json:"items"`
}

type DetailVACustomerInfoResponse struct {
	ID                      int64                    `json:"id"`
	FullName                string                   `json:"full_name"`
	Email                   string                   `json:"email"`
	AccountLevel            AccountLevelInfo         `json:"account_level"`
	Tier                    TierInfo                 `json:"tier"`
	Wallets                 []WalletBalanceInfo      `json:"wallets"`
	VAEnable                bool                     `json:"va_enable"`
	VirtualAccountStatistic VirtualAccountStatistic  `json:"virtual_account_statistic"`
	VirtualAccount          VirtualAccountInfo       `json:"virtual_account"`
	Balance                 *CustomerBalance         `json:"balance"`
	VARealBalance           float64                  `json:"va_real_balance"`
	VARealBalanceStr        string                   `json:"va_real_balance_str"`
	PendingAmount           map[string]PendingAmount `json:"pending_amount"`
}
type PendingAmountDetail struct {
	Amount               float64 `json:"amount"`
	PendingMonthlyAmount float64 `json:"pending_monthly_amount"`
	Currency             string  `json:"currency"`
	Provider             string  `json:"provider"`
	ProviderType         string  `json:"provider_type"`
}
type PendingAmount struct {
	Amount float64               `json:"amount"`
	Detail []PendingAmountDetail `json:"detail"`
}

type DetailVACustomerInfoRequest struct {
	ID int64 `uri:"id"`
}

type GetNotificationsRequest struct {
	NotificationGroup string `json:"notification_group"`
	Page              int64  `json:"page"`
	Limit             int64  `json:"limit"`
}
type GetNotificationsResponse struct {
	NextPage *int64                `json:"next_page"`
	PrevPage *int64                `json:"prev_page"`
	Page     int64                 `json:"page"`
	Limit    int64                 `json:"limit"`
	Total    int64                 `json:"total"`
	Data     []*model.Notification `json:"data"`
}
