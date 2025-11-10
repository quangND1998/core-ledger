package jobs

import (
	"time"

	"core-ledger/pkg/queue"
)

// DataProcessJob job xử lý dữ liệu theo pattern Laravel
type DataProcessJob struct {
	queue.BaseJob
	ProcessType string                 `json:"process_type"`
	Action      string                 `json:"action"`
	Data        map[string]interface{} `json:"data"`
	Options     map[string]interface{} `json:"options,omitempty"`
}

// GetPayload trả về payload của job
func (j *DataProcessJob) GetPayload() interface{} {
	return j
}

// GetType trả về loại job
func (j *DataProcessJob) GetType() string {
	return "data:process"
}

// NewDataProcessJob tạo data process job mới
func NewDataProcessJob(processType, action string, data map[string]interface{}) *DataProcessJob {
	return &DataProcessJob{
		BaseJob: queue.BaseJob{
			Queue: "default",
			Retry: 3,
		},
		ProcessType: processType,
		Action:      action,
		Data:        data,
	}
}

// SetOptions set processing options
func (j *DataProcessJob) SetOptions(options map[string]interface{}) {
	j.Options = options
}

// SetQueue set queue name
func (j *DataProcessJob) SetQueue(queue string) {
	j.Queue = queue
}

// SetDelay set delay time
func (j *DataProcessJob) SetDelay(delay time.Duration) {
	j.Delay = delay
}

// SetRetry set số lần retry
func (j *DataProcessJob) SetRetry(retry int) {
	j.Retry = retry
}
