package dto

// CoaAccountRuleTypeResp response cho TYPE
type CoaAccountRuleTypeResp struct {
	ID        uint64                  `json:"id"`
	Code      string                  `json:"code"`
	Name      string                  `json:"name"`
	Separator string                  `json:"separator"`
	Groups    []CoaAccountRuleGroupResp `json:"groups"`
}

// CoaAccountRuleGroupResp response cho GROUP
type CoaAccountRuleGroupResp struct {
	ID        uint64                 `json:"id"`
	Code      string                 `json:"code"`
	Name      string                 `json:"name"`
	InputType string                 `json:"input_type"`
	Separator string                 `json:"separator"`
	Steps     []CoaAccountRuleStepResp `json:"steps"`
}

// CoaAccountRuleStepResp response cho STEP
type CoaAccountRuleStepResp struct {
	StepID      uint64              `json:"step_id"`
	StepOrder   int                 `json:"step_order"`
	Type        string              `json:"type"` // SELECT hoặc TEXT
	Label       string              `json:"label"`
	CategoryCode string             `json:"category_code,omitempty"`
	Values      []RuleValueResp     `json:"values,omitempty"`
	InputCode   string              `json:"input_code,omitempty"`
	InputType   string              `json:"input_type,omitempty"`
	Separator   string              `json:"separator"`
}

