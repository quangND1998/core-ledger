package jobs

import (
	"time"

	"core-ledger/pkg/queue"
)

// DataProcessJob job xử lý dữ liệu theo pattern Laravel
type MyJob struct {
	queue.BaseJob
	ProcessType string                 `json:"process_type"`
	Action      string                 `json:"action"`
	Data        map[string]interface{} `json:"data"`
	Options     map[string]interface{} `json:"options,omitempty"`
}

// GetPayload trả về payload của job
func (j *MyJob) GetPayload() interface{} {
	return j
}

// GetType trả về loại job
func (j *MyJob) GetType() string {
	return "my:job"
}

// NewDataProcessJob tạo data process job mới
func NewMyJob(processType, action string, data map[string]interface{}) *MyJob {
	return &MyJob{
		BaseJob: queue.BaseJob{
			Queue: "default",
			Retry: 1,
		},
		ProcessType: processType,
		Action:      action,
		Data:        data,
	}
}

// SetOptions set processing options
func (j *MyJob) SetOptions(options map[string]interface{}) {
	j.Options = options
}

// SetQueue set queue name
func (j *MyJob) SetQueue(queue string) {
	j.Queue = queue
}

// SetDelay set delay time
func (j *MyJob) SetDelay(delay time.Duration) {
	j.Delay = delay
}

// SetRetry set số lần retry
func (j *MyJob) SetRetry(retry int) {
	j.Retry = retry
}
