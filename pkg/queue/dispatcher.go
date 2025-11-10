package queue

import (
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

// Dispatcher interface để gửi job (thay thế globals)
type Dispatcher interface {
	Dispatch(job Job, options ...DispatchOption) error
	DispatchLater(job Job, delay time.Duration, options ...DispatchOption) error
	DispatchAt(job Job, processAt time.Time, options ...DispatchOption) error
	DispatchOnQueue(job Job, queueName string, options ...DispatchOption) error
}

type asynqDispatcher struct {
	client *asynq.Client
}

func NewDispatcher(client *asynq.Client) Dispatcher {
	return &asynqDispatcher{client: client}
}

func (d *asynqDispatcher) Dispatch(job Job, options ...DispatchOption) error {
	if d.client == nil {
		return fmt.Errorf("queue client not initialized")
	}

	task, err := CreateTask(job)
	if err != nil {
		return fmt.Errorf("failed to create task: %v", err)
	}

	var asynqOpts []asynq.Option

	queueName := job.GetQueue()
	asynqOpts = append(asynqOpts, asynq.Queue(queueName))

	retry := job.GetRetry()
	if retry > 0 {
		asynqOpts = append(asynqOpts, asynq.MaxRetry(retry))
	}

	delay := job.GetDelay()
	if delay > 0 {
		asynqOpts = append(asynqOpts, asynq.ProcessIn(delay))
	}

	for _, option := range options {
		asynqOpts = append(asynqOpts, option(task))
	}

	_, err = d.client.Enqueue(task, asynqOpts...)
	return err
}

func (d *asynqDispatcher) DispatchLater(job Job, delay time.Duration, options ...DispatchOption) error {
	job.SetDelay(delay)
	return d.Dispatch(job, options...)
}

func (d *asynqDispatcher) DispatchAt(job Job, processAt time.Time, options ...DispatchOption) error {
	options = append(options, ProcessAt(processAt))
	return d.Dispatch(job, options...)
}

func (d *asynqDispatcher) DispatchOnQueue(job Job, queueName string, options ...DispatchOption) error {
	job.SetQueue(queueName)
	return d.Dispatch(job, options...)
}


