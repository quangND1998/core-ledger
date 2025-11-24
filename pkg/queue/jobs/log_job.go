package jobs

import (
	"time"

	"core-ledger/pkg/queue"
	"gorm.io/datatypes"
)

// LogJob job để log database changes
type LogJob struct {
	queue.BaseJob
	LoggableID   uint64         `json:"loggable_id"`
	LoggableType string         `json:"loggable_type"`
	Action       string         `json:"action"`
	OldValue     datatypes.JSON `json:"old_value,omitempty"`
	NewValue     datatypes.JSON `json:"new_value,omitempty"`
	Metadata     datatypes.JSON `json:"metadata,omitempty"`
	CreatedBy    *uint64        `json:"created_by,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
}

// GetPayload trả về payload của job
func (j *LogJob) GetPayload() interface{} {
	return j
}

// GetType trả về loại job
func (j *LogJob) GetType() string {
	return "log:create"
}

// NewLogJob tạo log job mới
func NewLogJob(loggableID uint64, loggableType, action string, oldValue, newValue, metadata datatypes.JSON, createdBy *uint64, createdAt time.Time) *LogJob {
	return &LogJob{
		BaseJob: queue.BaseJob{
			Queue: "logs", // Queue riêng cho logs
			Retry: 3,      // Retry 3 lần nếu fail
		},
		LoggableID:   loggableID,
		LoggableType: loggableType,
		Action:       action,
		OldValue:     oldValue,
		NewValue:     newValue,
		Metadata:     metadata,
		CreatedBy:    createdBy,
		CreatedAt:    createdAt,
	}
}




