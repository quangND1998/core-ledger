package queue

import (
	"time"

	"github.com/hibiken/asynq"
)

// DispatchOption: option pattern cho enqueue
type DispatchOption func(*asynq.Task) asynq.Option

// ProcessIn: set delay
func ProcessIn(delay time.Duration) DispatchOption {
	return func(task *asynq.Task) asynq.Option {
		return asynq.ProcessIn(delay)
	}
}

// ProcessAt: set thời điểm xử lý
func ProcessAt(processAt time.Time) DispatchOption {
	return func(task *asynq.Task) asynq.Option {
		return asynq.ProcessAt(processAt)
	}
}

// Queue: set queue name
func Queue(queueName string) DispatchOption {
	return func(task *asynq.Task) asynq.Option {
		return asynq.Queue(queueName)
	}
}

// Retry: set số lần retry
func Retry(retries int) DispatchOption {
	return func(task *asynq.Task) asynq.Option {
		return asynq.MaxRetry(retries)
	}
}

// Timeout: set timeout
func Timeout(timeout time.Duration) DispatchOption {
	return func(task *asynq.Task) asynq.Option {
		return asynq.Timeout(timeout)
	}
}


