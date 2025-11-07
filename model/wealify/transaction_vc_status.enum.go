package model

type TransactionVcStatus string

const (
	TransactionVcStatusProcessing TransactionVcStatus = "PROCESSING"
	TransactionVcStatusSuccess    TransactionVcStatus = "SUCCESS"
	TransactionVcStatusFailure    TransactionVcStatus = "FAILURE"
)

func (s TransactionVcStatus) String() string {
	return string(s)
}
