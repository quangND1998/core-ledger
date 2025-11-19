package dto

type RuleValueRequest struct {
	ID       uint64 `json:"id,omitempty"`
	Value    string `json:"value" binding:"required"`
	Name     string `json:"name" binding:"required"`
	IsDelete *bool  `json:"is_delete" binding:"boolean,required"`
}

type SaveRuleValueRequest struct {
	Data       []*RuleValueRequest `json:"data" binding:"required,dive"`
	CategoryID uint                `json:"category_id" binding:"required"`
}

type FilterRuleValueRequest struct {
	CategoryID uint `form:"category_id" binding:"required"`
}
