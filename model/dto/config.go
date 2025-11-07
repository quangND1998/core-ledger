package dto

type AirwallexKeySets []AirwallexKeySet
type AirwallexKeySet struct {
	Email           string `json:"email"`
	SystemPaymentID string `json:"id"`
	ApiKey          string `json:"api_key"`
	ClientID        string `json:"client_id"`
}
