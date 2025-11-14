package dto

type ListEntrytFilter struct {
	BasePaginationQuery
	Search   *string   `json:"search,omitempty" form:"search"`
	Type     *Dc       `json:"type,omitempty" form:"type"`
	Currency []*string `json:"currency,omitempty" form:"currency[]"`
	Sort     *string   `json:"sort,omitempty" form:"sort"`
}

type Dc string

const (
	Debit  Dc = "D"
	Credit Dc = "C"
)
