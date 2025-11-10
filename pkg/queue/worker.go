package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"reflect"
	"time"

	"github.com/hibiken/asynq"
)

// JobFactory function để tạo job instance mới
type JobFactory func() Job

// Worker quản lý việc xử lý jobs
type Worker struct {
	server    *asynq.Server
	mux       *asynq.ServeMux
	factories map[string]JobFactory
	handlers  map[string]JobHandler
}

// defaultRetryDelayFunc tạo exponential backoff cho retry
// Retry delay: 1s, 2s, 4s, 8s, 16s, 30s (max)
func defaultRetryDelayFunc(n int, e error, t *asynq.Task) time.Duration {
	// Exponential backoff: 2^(n-1) seconds, max 30 seconds
	delay := time.Duration(math.Pow(2, float64(n-1))) * time.Second
	if delay > 30*time.Second {
		delay = 30 * time.Second
	}
	return delay
}

// retryDelayFuncWithJobBackoff đọc backoff từ job payload nếu có, nếu không dùng default
func retryDelayFuncWithJobBackoff(n int, e error, t *asynq.Task) time.Duration {
	// Parse payload để lấy backoff
	var payload JobPayload
	if err := json.Unmarshal(t.Payload(), &payload); err == nil {
		// Nếu job có backoff được định nghĩa
		if len(payload.Backoff) > 0 {
			// n là số lần retry (0-based: 0, 1, 2, ...)
			// Retry #0 → array[0], Retry #1 → array[1], Retry #2 → array[2], ...
			index := n
			if index < 0 {
				index = 0
			}
			if index >= len(payload.Backoff) {
				// Nếu vượt quá mảng, dùng giá trị cuối cùng
				index = len(payload.Backoff) - 1
			}
			delay := time.Duration(payload.Backoff[index]) * time.Second
			return delay
		}
	}
	// Fallback về default exponential backoff
	return defaultRetryDelayFunc(n, e, t)
}

// NewWorker tạo worker instance mới với config
func NewWorker(redisAddr string, concurrency int, queues map[string]int) *Worker {
	var w *Worker
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency:   concurrency,
			Queues:        queues,
			RetryDelayFunc: retryDelayFuncWithJobBackoff,
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, t *asynq.Task, err error) {
				if w == nil {
					return
				}
				jobType := t.Type()
				handler, ok := w.handlers[jobType]
				if !ok {
					return
				}
				// chỉ gọi Failed khi hết retry
				retryCount, okRetry := asynq.GetRetryCount(ctx)
				maxRetry, okMax := asynq.GetMaxRetry(ctx)
				// Nếu không lấy được metadata retry, bỏ qua để tránh gọi Failed sai thời điểm
				if !(okRetry && okMax) {
					return
				}
				if maxRetry > 0 && (retryCount+1) < maxRetry {
					return
				}
				// dựng lại job
				factory, exists := w.factories[jobType]
				if !exists {
					return
				}
				job := factory()
				if job == nil {
					return
				}
				if errPopulate := w.populateJobData(job, t.Payload()); errPopulate != nil {
					log.Printf("failed to populate job for Failed hook: %v", errPopulate)
				}
				if fh, ok := handler.(failableHandler); ok {
					fh.Failed(ctx, job, err)
				}
			}),
		},
	)

	w = &Worker{
		server:    srv,
		mux:       asynq.NewServeMux(),
		factories: make(map[string]JobFactory),
		handlers:  make(map[string]JobHandler),
	}
	return w
}

// NewWorkerWithRedis sử dụng đầy đủ RedisClientOpt (addr/password/db/..)
func NewWorkerWithRedis(opt asynq.RedisClientOpt, concurrency int, queues map[string]int) *Worker {
	var w *Worker
	srv := asynq.NewServer(
		opt,
		asynq.Config{
			Concurrency:   concurrency,
			Queues:        queues,
			RetryDelayFunc: retryDelayFuncWithJobBackoff,
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, t *asynq.Task, err error) {
				if w == nil {
					return
				}
				jobType := t.Type()
				handler, ok := w.handlers[jobType]
				if !ok {
					return
				}
				retryCount, okRetry := asynq.GetRetryCount(ctx)
				maxRetry, okMax := asynq.GetMaxRetry(ctx)
				if !(okRetry && okMax) {
					return
				}
				if maxRetry > 0 && (retryCount+1) < maxRetry {
					return
				}
				factory, exists := w.factories[jobType]
				if !exists {
					return
				}
				job := factory()
				if job == nil {
					return
				}
				if errPopulate := w.populateJobData(job, t.Payload()); errPopulate != nil {
					log.Printf("failed to populate job for Failed hook: %v", errPopulate)
				}
				if fh, ok := handler.(failableHandler); ok {
					fh.Failed(ctx, job, err)
				}
			}),
		},
	)
	w = &Worker{
		server:    srv,
		mux:       asynq.NewServeMux(),
		factories: make(map[string]JobFactory),
		handlers:  make(map[string]JobHandler),
	}
	return w
}

// RegisterJob đăng ký job factory và handler
// jobTemplate: struct mẫu để tạo instance và unmarshal payload
// handler: đối tượng thực hiện xử lý job (riêng biệt)
func (w *Worker) RegisterJob(jobType string, jobTemplate Job, handler JobHandler) {
	// Tạo factory function từ job template
	w.factories[jobType] = func() Job {
		// Sử dụng reflection để tạo instance mới
		jobValue := reflect.ValueOf(jobTemplate)
		if jobValue.Kind() == reflect.Ptr {
			jobValue = jobValue.Elem()
		}
		newJobValue := reflect.New(jobValue.Type())
		return newJobValue.Interface().(Job)
	}

	// Lưu handler (có thể nil nếu muốn fallback dùng method trên job)
	if handler != nil {
		w.handlers[jobType] = handler
	}

	w.mux.HandleFunc(jobType, w.createHandler(jobType))
}

// createHandler tạo handler function cho asynq
func (w *Worker) createHandler(jobType string) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		// Lấy factory từ registry
		factory, exists := w.factories[jobType]
		if !exists {
			return fmt.Errorf("no job factory registered for job type: %s", jobType)
		}

		// Tạo job instance mới từ factory
		job := factory()
		if job == nil {
			return fmt.Errorf("failed to create job instance for type: %s", jobType)
		}

		// Unmarshal JobPayload vào job instance
		if err := w.populateJobData(job, t.Payload()); err != nil {
			return fmt.Errorf("failed to populate job data: %w", err)
		}

		var handlerErr error

		// Nếu có handler riêng, gọi handler đó
		if h, ok := w.handlers[jobType]; ok && h != nil {
			log.Printf("Processing job with external handler: %s", jobType)
			handlerErr = h.Handle(ctx, job)
		} else if legacyHandler, ok := any(job).(interface {
			Handle(context.Context) error
		}); ok {
			// Fallback: nếu job struct vẫn có method Handle (legacy), gọi nó
			log.Printf("Processing job with legacy job.Handle: %s", jobType)
			handlerErr = legacyHandler.Handle(ctx)
		} else {
			return fmt.Errorf("no handler found for job type: %s", jobType)
		}

		// Kiểm tra timeout: nếu context bị timeout hoặc error là timeout
		if handlerErr != nil {
			// Kiểm tra xem context có bị timeout không
			if ctx.Err() == context.DeadlineExceeded {
				return fmt.Errorf("job timeout: %w", handlerErr)
			}
			// Kiểm tra xem error có phải là timeout error không
			if errors.Is(handlerErr, context.DeadlineExceeded) {
				return fmt.Errorf("job timeout: %w", handlerErr)
			}
			return handlerErr
		}

		// Kiểm tra lại context sau khi handler hoàn thành (phòng trường hợp timeout xảy ra nhưng handler không trả về error)
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("job timeout: context deadline exceeded")
		}

		return nil
	}
}

// populateJobData điền data vào job instance từ JobPayload
func (w *Worker) populateJobData(job Job, payload []byte) error {
	var jobPayload JobPayload
	if err := json.Unmarshal(payload, &jobPayload); err != nil {
		return fmt.Errorf("failed to unmarshal job payload: %w", err)
	}

	// Marshal lại data để unmarshal vào job struct
	data, err := json.Marshal(jobPayload.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal job data: %w", err)
	}

	// Unmarshal vào job instance
	return json.Unmarshal(data, job)
}

// failableHandler mô tả handler có hook Failed
type failableHandler interface {
	JobHandler
	Failed(ctx context.Context, job Job, err error)
}

// Start chạy worker
func (w *Worker) Start() error {
	log.Println("Starting queue worker...")
	return w.server.Run(w.mux)
}

// Stop dừng worker
func (w *Worker) Stop() {
	w.server.Shutdown()
}

// Global worker instance
var defaultWorker *Worker

// InitWorker khởi tạo worker global với config
func InitWorker(redisAddr string, concurrency int, queues map[string]int) {
	defaultWorker = NewWorker(redisAddr, concurrency, queues)
}

// RegisterJobHandler đăng ký job handler global
// Thay đổi: nhận cả job template và handler
func RegisterJobHandler(jobType string, jobTemplate Job, handler JobHandler) {
	if defaultWorker == nil {
		log.Fatal("Worker not initialized. Call InitWorker() first.")
	}
	defaultWorker.RegisterJob(jobType, jobTemplate, handler)
}

// StartWorker chạy worker global
func StartWorker() error {
	if defaultWorker == nil {
		log.Fatal("Worker not initialized. Call InitWorker() first.")
	}
	return defaultWorker.Start()
}

// StopWorker dừng worker global
func StopWorker() {
	if defaultWorker != nil {
		defaultWorker.Stop()
	}
}
