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
