package queue

import (
	"fmt"
	"log"
	"time"

	config "core-ledger/configs"

	"github.com/hibiken/asynq"
)

var (
	Client *asynq.Client
	Server *asynq.Server
	Config *config.QueueConfig
)

// InitQueue khởi tạo hệ thống queue với config
func InitQueue() error {
	// Lấy cấu hình queue
	fmt.Println("Initializing queue system...")
	queueConfig, err := config.GetQueueConfigWithValidation()
	fmt.Println("Queue config:", queueConfig)
	if err != nil {
		return fmt.Errorf("invalid queue config: %v", err)
	}

	Config = queueConfig

	// Tạo Redis client cho Asynq
	redisOpt := asynq.RedisClientOpt{
		Addr:     queueConfig.RedisAddr,
		Password: queueConfig.RedisPassword,
		DB:       queueConfig.RedisDB,
	}

	// Khởi tạo client để gửi jobs
	Client = asynq.NewClient(redisOpt)
	if Client == nil {
		return fmt.Errorf("failed to create Asynq client")
	}

	// Khởi tạo server để xử lý jobs
	Server = asynq.NewServer(
		redisOpt,
		asynq.Config{
			// Số lượng worker đồng thời
			Concurrency: queueConfig.Concurrency,
			// Các queue ưu tiên
			Queues: queueConfig.Queues,
			// Retry policy
			RetryDelayFunc: func(n int, err error, t *asynq.Task) time.Duration {
				return time.Duration(n) * time.Second
			},
		},
	)

	log.Printf("Queue system initialized successfully with config: Redis=%s, Concurrency=%d, Queues=%v",
		queueConfig.RedisAddr, queueConfig.Concurrency, queueConfig.Queues)
	return nil
}

// CloseQueue đóng kết nối queue
func CloseQueue() {
	if Client != nil {
		Client.Close()
	}
	if Server != nil {
		Server.Shutdown()
	}
}

// Dispatch helper function giống Laravel để dispatch job
func Dispatch(job Job, options ...DispatchOption) error {
	fmt.Println("Dispatching job:", job.GetType())

	if Client == nil {
		return fmt.Errorf("queue client not initialized")
	}

	// Tạo task từ job
	task, err := CreateTask(job)
	if err != nil {
		return fmt.Errorf("failed to create task: %v", err)
	}

	// Tạo asynq options
	var asynqOpts []asynq.Option

	// Set queue name từ job
	queueName := job.GetQueue()
	asynqOpts = append(asynqOpts, asynq.Queue(queueName))

	// Set retry từ job
	retry := job.GetRetry()
	if retry > 0 {
		asynqOpts = append(asynqOpts, asynq.MaxRetry(retry))
	}

	// Set delay từ job
	delay := job.GetDelay()
	if delay > 0 {
		asynqOpts = append(asynqOpts, asynq.ProcessIn(delay))
	}

	// Apply custom options
	for _, option := range options {
		asynqOpts = append(asynqOpts, option(task))
	}

	// Gửi task
	info, err := Client.Enqueue(task, asynqOpts...)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %v", err)
	}

	log.Printf("Job dispatched: %s, queue: %s, id: %s", job.GetType(), queueName, info.ID)
	return nil
}

// DispatchLater dispatch job với delay
func DispatchLater(job Job, delay time.Duration, options ...DispatchOption) error {
	job.SetDelay(delay)
	return Dispatch(job, options...)
}

// DispatchAt dispatch job tại thời điểm cụ thể
func DispatchAt(job Job, processAt time.Time, options ...DispatchOption) error {
	options = append(options, ProcessAt(processAt))
	return Dispatch(job, options...)
}

// DispatchOnQueue dispatch job vào queue cụ thể
func DispatchOnQueue(job Job, queueName string, options ...DispatchOption) error {
	job.SetQueue(queueName)
	return Dispatch(job, options...)
}

// DispatchOption function type cho các option
type DispatchOption func(*asynq.Task) asynq.Option

// ProcessIn option để set delay
func ProcessIn(delay time.Duration) DispatchOption {
	return func(task *asynq.Task) asynq.Option {
		return asynq.ProcessIn(delay)
	}
}

// ProcessAt option để set thời gian xử lý
func ProcessAt(processAt time.Time) DispatchOption {
	return func(task *asynq.Task) asynq.Option {
		return asynq.ProcessAt(processAt)
	}
}

// Queue option để set queue name
func Queue(queueName string) DispatchOption {
	return func(task *asynq.Task) asynq.Option {
		return asynq.Queue(queueName)
	}
}

// Retry option để set số lần retry
func Retry(retries int) DispatchOption {
	return func(task *asynq.Task) asynq.Option {
		return asynq.MaxRetry(retries)
	}
}

// Timeout option để set timeout
func Timeout(timeout time.Duration) DispatchOption {
	return func(task *asynq.Task) asynq.Option {
		return asynq.Timeout(timeout)
	}
}
