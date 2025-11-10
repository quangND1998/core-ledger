package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"

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

// NewWorker tạo worker instance mới với config
func NewWorker(redisAddr string, concurrency int, queues map[string]int) *Worker {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: concurrency,
			Queues:      queues,
		},
	)

	return &Worker{
		server:    srv,
		mux:       asynq.NewServeMux(),
		factories: make(map[string]JobFactory),
		handlers:  make(map[string]JobHandler),
	}
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

		// Nếu có handler riêng, gọi handler đó
		if h, ok := w.handlers[jobType]; ok && h != nil {
			log.Printf("Processing job with external handler: %s", jobType)
			return h.Handle(ctx, job)
		}

		// Fallback: nếu job struct vẫn có method Handle (legacy), gọi nó
		if legacyHandler, ok := any(job).(interface {
			Handle(context.Context) error
		}); ok {
			log.Printf("Processing job with legacy job.Handle: %s", jobType)
			return legacyHandler.Handle(ctx)
		}

		return fmt.Errorf("no handler found for job type: %s", jobType)
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
