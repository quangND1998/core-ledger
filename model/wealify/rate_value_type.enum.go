package model

type RateValueType string

const (
	RateValueTypeManual RateValueType = "MANUAL"
	RateValueTypeLinked RateValueType = "LINKED"
)

func (r RateValueType) String() string {
	return string(r)
}
