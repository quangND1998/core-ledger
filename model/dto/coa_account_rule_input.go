package dto

// CoaAccountRuleInput DTO để nhận input từ frontend dựa trên rules
type CoaAccountRuleInput struct {
	TypeCode  string                    `json:"type_code" binding:"required"`  // ASSET, LIAB, REV, EXP
	GroupCode *string                   `json:"group_code,omitempty"`         // FLOAT, BANK, etc. (NULL cho REV, EXP)
	Steps     []CoaAccountRuleStepInput `json:"steps" binding:"required,dive"` // Các giá trị đã chọn cho từng step
}

// CoaAccountRuleStepInput DTO cho input của mỗi step
type CoaAccountRuleStepInput struct {
	StepID      uint64  `json:"step_id" binding:"required"`      // ID của step
	StepOrder   int     `json:"step_order" binding:"required"`   // Thứ tự step
	CategoryCode *string `json:"category_code,omitempty"`         // Nếu step là SELECT
	ValueID     *uint64 `json:"value_id,omitempty"`              // ID của value đã chọn (nếu SELECT)
	Value       *string `json:"value,omitempty"`                 // Value đã chọn (nếu SELECT) hoặc text input (nếu TEXT)
	InputCode   *string `json:"input_code,omitempty"`            // Nếu step là TEXT
}

// CoaAccountDataWithRules DTO kết hợp account data với rule inputs
type CoaAccountDataWithRules struct {
	// Basic fields
	Code        string  `json:"code,omitempty"`
	AccountNo   string  `json:"account_no" binding:"required"`
	Name        string  `json:"name,omitempty"`
	Type        string  `json:"type,omitempty"`
	Currency    string  `json:"currency,omitempty"`
	Status      string  `json:"status" binding:"required"`
	Description *string `json:"description,omitempty"`
	
	// Rule inputs để validate và build account_no
	RuleInput *CoaAccountRuleInput `json:"rule_input,omitempty"`
}


