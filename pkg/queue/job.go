package queue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
)

// Job interface: giữ data/metadata của job (không chứa Handle nữa)
type Job interface {
	GetPayload() interface{}
	GetType() string
	GetQueue() string
	GetDelay() time.Duration
	GetRetry() int
	SetQueue(string)
	SetDelay(time.Duration)
	SetRetry(int)
}

// JobHandler interface tách riêng phần xử lý
type JobHandler interface {
	Handle(ctx context.Context, job Job) error
}

// BaseJob struct để embed vào các job cụ thể
type BaseJob struct {
	Queue string        `json:"queue,omitempty"`
	Delay time.Duration `json:"delay,omitempty"`
	Retry int           `json:"retry,omitempty"`
}

// GetQueue trả về tên queue, mặc định là "default"
func (b *BaseJob) GetQueue() string {
	if b.Queue == "" {
		return "default"
	}
	return b.Queue
}

// GetDelay trả về thời gian delay
func (b *BaseJob) GetDelay() time.Duration {
	return b.Delay
}

// GetRetry trả về số lần retry
func (b *BaseJob) GetRetry() int {
	if b.Retry == 0 {
		return 3 // mặc định retry 3 lần
	}
	return b.Retry
}

// SetQueue set queue name
func (b *BaseJob) SetQueue(queue string) {
	b.Queue = queue
}

// SetDelay set delay time
func (b *BaseJob) SetDelay(delay time.Duration) {
	b.Delay = delay
}

// SetRetry set số lần retry
func (b *BaseJob) SetRetry(retry int) {
	b.Retry = retry
}

// JobPayload wraps job data for serialization
type JobPayload struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// CreateTask tạo asynq.Task từ Job
func CreateTask(job Job) (*asynq.Task, error) {
	payload := JobPayload{
		Type: job.GetType(),
		Data: job.GetPayload(),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(job.GetType(), data), nil
}
