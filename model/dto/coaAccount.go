package dto

import model "core-ledger/model/core-ledger"

type ListCoaAccountFilter struct {
	BasePaginationQuery
	Search    *string   `json:"search,omitempty" form:"search"`
	Status    []string  `json:"status,omitempty" form:"status[]"`
	Types     []*string `json:"types,omitempty" form:"types[]"`
	Currency  []*string `json:"currency,omitempty" form:"currency[]"`
	Networks  []*string `json:"networks,omitempty" form:"networks[]"`
	Providers []*string `json:"providers,omitempty" form:"providers[]"`
	Sort      *string   `json:"sort,omitempty" form:"sort"`
}

type CoaAccountDetailResponse struct {
	CoaAccount *model.CoaAccount `json:"coa_account,omitempty"`
	Entries    []model.Entry     `json:"entries"`
	Snapshots  []model.Snapshot  `json:"snapshots"`
}

type CoaAccountCreateRequest struct {
	Name        string `json:"name,omitempty" binding:"required"`
	AccountNo   string `json:"account_no,omitempty" binding:"required"`
	Code        string `json:"code,omitempty" binding:"required"`
	Type        string `json:"type,omitempty" binding:"required"`
	Status      string `json:"status,omitempty" binding:"required"`
	Description string `json:"description,omitempty"`
	Provider    string `json:"provider,omitempty" `
	Network     string `json:"network,omitempty" `
	Currency    string `json:"currency,omitempty" `
}
type CoaAccountExistAccountNo struct {
	AccountNo string `json:"account_no,omitempty" binding:"required"`
}

type UpdateCoaAccountRequest struct {
	Name string `json:"name,omitempty" binding:"required"`
}

type UpdateStatusCoaAccountRequest struct {
	ID     int64  `json:"id,omitempty" binding:"required"`
	Status string `json:"status,omitempty" binding:"required"`
}
