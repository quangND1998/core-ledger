package handlers

import (
	"context"
	"core-ledger/pkg/queue"
	"core-ledger/pkg/queue/jobs"
	"core-ledger/pkg/repo"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

// DataProcessHandler xử lý DataProcessJob
type DataProcessHandler struct {
	repo.TransactionRepo
	dispatcher queue.Dispatcher

	// thêm dependency nếu cần (ví dụ: services, repos)
}

func NewDataProcessHandler(transactionRepo repo.TransactionRepo, dispatcher queue.Dispatcher) *DataProcessHandler {
	return &DataProcessHandler{
		TransactionRepo: transactionRepo,
		dispatcher:      dispatcher,
	}
}

// NewDataProcessRegistration: provider đăng ký job/handler vào group "queue-registrations"
func NewDataProcessRegistration(h *DataProcessHandler) queue.Registration {
	return queue.Registration{
		Type:     "data:process",
		Template: &jobs.DataProcessJob{},
		Handler:  h,
	}
}

func (h *DataProcessHandler) Handle(ctx context.Context, j queue.Job) error {
	// kiểu assert về concrete job
	job, ok := j.(*jobs.DataProcessJob)
	if !ok {
		return fmt.Errorf("invalid job type, expect *DataProcessJob")
	}

	// Log retry info
	currentAttempt := 1
	if n, ok := asynq.GetRetryCount(ctx); ok {
		currentAttempt = n + 1
	}
	maxRetry := 0
	if max, ok := asynq.GetMaxRetry(ctx); ok {
		maxRetry = max
	}

	// Log backoff info nếu có
	backoff := job.GetBackoff()
	backoffInfo := "none (using default)"
	if len(backoff) > 0 {
		backoffInfo = fmt.Sprintf("%v seconds", backoff)
	}

	log.Printf("📦 [Job] Handling DataProcessJob: Type=%s, Action=%s | Attempt=%d/%d | Backoff=%s",
		job.ProcessType, job.Action, currentAttempt, maxRetry+1, backoffInfo)

	// Test lỗi: đặt Action="fail" để cố tình trả về lỗi (kích hoạt retry/Failed)
	if job.Action == "fail" {
		log.Printf("❌ [Job] Forcing failure for testing (will retry with backoff)")
		return fmt.Errorf("forced failure for testing")
	}

	dataJob := jobs.NewMyJob("user_analytics", "export", map[string]interface{}{
		"user_id": "test_user_123",
		"format":  "json",
		"filters": map[string]interface{}{
			"date_from": "2024-01-01",
			"date_to":   "2024-12-31",
			"status":    "active",
		},
	})

	dataJob.SetOptions(map[string]interface{}{
		"include_headers": true,
		"date_format":     "ISO",
		"compression":     "none",
		"max_records":     1000,
	})
	dataJob.SetQueue("critical")
	if err := h.dispatcher.Dispatch(dataJob, queue.Timeout(1*time.Second)); err != nil {
		log.Printf("❌ Failed to dispatch data job: %v", err)
		return err
	}
	// var transactions []model.Transaction
	// transactions, err := h.TransactionRepo.GetList(ctx)
	// log.Printf("Fetched %d transactions", len(transactions))
	// if err != nil {
	// 	return err
	// }
	// TODO: business logic xử lý theo job.ProcessType / job.Action / job.Data

	return nil // trả về error để asynq retry nếu cần
}

// Failed: hook được gọi khi job đã hết retry hoặc timeout
func (h *DataProcessHandler) Failed(ctx context.Context, j queue.Job, err error) {
	// cố gắng assert đúng loại job để log chi tiết
	if job, ok := j.(*jobs.DataProcessJob); ok {
		log.Printf("[FAILED] DataProcessJob Type=%s Action=%s Error=%v", job.ProcessType, job.Action, err)
	} else {
		log.Printf("[FAILED] DataProcessJob Error=%v", err)
	}
	// TODO: Có thể ghi log vào DB, tạo transaction_log, hoặc đẩy sang channel cảnh báo...
	_ = h // giữ chỗ nếu sau này cần dùng repo để lưu DB
}
