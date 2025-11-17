package ruleValue

type RuleValueRequest struct {
	ID       uint64 `json:"id"`
	Code     string `json:"code" binding:"required,unique,unique_ruleCode"`
	Name     string `json:"name" binding:"required"`
	IsDelete bool   `json:"is_delete" binding:"required"`
}

type SaveRuleValueRequest struct {
	Data       []RuleValueRequest `json:"data"  binding:"required,dive,required,min=2"`
	CategoryID uint64             `json:"category_id" binding:"required"`
}
