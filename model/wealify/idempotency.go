package model

import "time"

type IdempotencyRecordStatus string

const (
	IdempotencyInProgress IdempotencyRecordStatus = "IN_PROGRESS"
	IdempotencyCompleted  IdempotencyRecordStatus = "COMPLETED"
	IdempotencyFailed     IdempotencyRecordStatus = "FAILED"
)

// IdempotencyRecord is our GORM model for the table
type IdempotencyRecord struct {
	IdempotencyKey string                  `gorm:"primaryKey" json:"idempotency_key,omitempty"`
	EmployeeID     int64                   `gorm:"employee_id" json:"employee_id"`
	RequestHash    string                  `gorm:"request_hash" json:"request_hash"`
	Status         IdempotencyRecordStatus `gorm:"status" json:"status"` // IN_PROGRESS, COMPLETED, FAILED
	ResponseCode   int                     `gorm:"response_code" json:"response_code"`
	ResponseBody   []byte                  `gorm:"type:jsonb" json:"response_body"`
	ExpiresAt      time.Time               `gorm:"expires_at" json:"expires_at"`
	CreatedAt      time.Time               `gorm:"created_at" json:"created_at"`
}
