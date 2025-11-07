package dto

type DirectionType string

const (
	DirectionTypeAsc  DirectionType = "ASC"
	DirectionTypeDesc DirectionType = "DESC"
)

func (d DirectionType) String() string {
	return string(d)
}

type PaginationDto struct {
	Page     *int64 `json:"page" form:"page"`
	Limit    *int64 `json:"limit" form:"limit"`
	DataSort *int64 `json:"data_sort" form:"data_sort"`
}

type SortType struct {
	Direction DirectionType
	Field     string
}
