package model

type FeeType string

const (
	FeeTypeFixed   FeeType = "FIXED"
	FeeTypePercent FeeType = "PERCENT"
)

func (f FeeType) String() string {
	return string(f)
}
