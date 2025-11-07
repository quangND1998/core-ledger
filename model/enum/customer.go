package enum

type Tier string

const (
	TierStandard Tier = "STANDARD"
	TierSilver   Tier = "SILVER"
	TierGold     Tier = "GOLD"
	TierDiamond  Tier = "DIAMOND"
)

func (t Tier) String() string {
	return string(t)
}

type CustomerType string

const (
	CustomerTypeIndividual CustomerType = "INDIVIDUAL"
	CustomerTypeBusiness   CustomerType = "BUSINESS"
)

func (t CustomerType) String() string {
	return string(t)
}
