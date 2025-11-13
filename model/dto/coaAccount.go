package dto

import model "core-ledger/model/core-ledger"

type ListCoaAccountFilter struct {
	BasePaginationQuery
	Name      *string `json:"name,omitempty" form:"name"`
	Code      *string `json:"code,omitempty" form:"code"`
	Status    *string `json:"status,omitempty" form:"status"`
	Type      *string `json:"type,omitempty" form:"type"`
	AccountNo *string `json:"account_no,omitempty" form:"account_no"`
}

type CoaAccountDetailResponse struct {
	CoaAccount *model.CoaAccount `json:"coa_account,omitempty"`
	Entries    []model.Entry     `json:"entries"`
	Snapshots  []model.Snapshot  `json:"snapshots"`
}
