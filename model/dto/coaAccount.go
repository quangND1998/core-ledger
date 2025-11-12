package dto

type ListCoaAccountFilter struct {
	BasePaginationQuery
	Name *string `json:"name,omitempty" form:"name"`
	Code *string `json:"code,omitempty" form:"code"`
}
