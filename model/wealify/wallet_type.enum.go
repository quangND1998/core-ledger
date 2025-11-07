package model

type WalletType string

const (
	WalletTypeMain   = "MAIN"
	WalletTypeVA     = "VA"
	WalletTypeSystem = "SYSTEM"
)

func (w WalletType) String() string {
	return string(w)
}
