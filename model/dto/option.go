package dto

import model "core-ledger/model/core-ledger"

type StepResponse struct {
	StepID    uint64  `json:"step_id"`
	StepOrder int     `json:"step_order"`
	Type      string  `json:"type"` // dropdown | input
	Label     *string `json:"label"`

	// dropdown
	CategoryCode string            `json:"category_code,omitempty"`
	Values       []model.RuleValue `json:"values,omitempty"`

	// input
	InputType string `json:"input_type,omitempty"`
}

type AccountRuleOptionTree struct {
	ID       uint64                   `json:"id"`
	Code     string                   `json:"code"`
	Name     string                   `json:"name"`
	LayerID  uint64                   `json:"layer_id,omitempty"`
	ParentID *uint64                  `json:"parent_id,omitempty"`
	Children []*AccountRuleOptionTree `json:"children"`
}
type RuleValueResp struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type RuleStepResp struct {
	StepID    uint64 `json:"step_id"`
	StepOrder int    `json:"step_order"`
	Type      string `json:"type"` // dropdown | input
	Label     string `json:"label"`

	CategoryCode string          `json:"category_code,omitempty"`
	Values       []RuleValueResp `json:"values,omitempty"`

	InputCode *string `json:"input_code,omitempty"`
	InputType string  `json:"input_type,omitempty"`
}
type RuleGroupResp struct {
	ID        uint64         `json:"id"`
	Code      string         `json:"code"`
	Name      string         `json:"name"`
	InputType string         `json:"input_type"`
	Steps     []RuleStepResp `json:"steps"`
}
type RuleTypeResp struct {
	ID     uint64          `json:"id"`
	Code   string          `json:"code"`
	Name   string          `json:"name"`
	Groups []RuleGroupResp `json:"groups"`
}
