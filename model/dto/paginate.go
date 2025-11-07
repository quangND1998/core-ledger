package dto

type BasePaginationQuery struct {
	IDs       []int64  `form:"ids" json:"ids,omitempty"`
	StringIDs []string `form:"string_ids" json:"string_ids,omitempty"`
	Keyword   *string  `form:"keyword" json:"keyword,omitempty"`
	Page      *int64   `form:"page" json:"page,omitempty"`
	Limit     *int64   `form:"limit" json:"limit,omitempty"`
	IsDeleted *bool    `form:"is_deleted" json:"is_deleted,omitempty"`
	Order     *any     `form:"order" json:"order,omitempty"`
	Cursor    *string  `form:"cursor" json:"cursor,omitempty"`
	StartDate *string  `form:"start_date" json:"start_date,omitempty"`
	EndDate   *string  `form:"end_date" json:"end_date,omitempty"`
}

type PaginationResponse[T any] struct {
	Items     []T    `json:"items"`
	Total     int64  `json:"total"`
	Limit     int64  `json:"limit"`
	Page      int64  `json:"page"`
	TotalPage int64  `json:"total_page"`
	NextPage  *int64 `json:"next_page"`
	PrevPage  *int64 `json:"prev_page"`
}
