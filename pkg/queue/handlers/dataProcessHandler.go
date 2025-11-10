package handlers

import (
	"context"
	"core-ledger/pkg/queue"
	"core-ledger/pkg/queue/jobs"
	"core-ledger/pkg/repo"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
)

// DataProcessHandler xử lý DataProcessJob
type DataProcessHandler struct {
	repo.TransactionRepo
	// thêm dependency nếu cần (ví dụ: services, repos)
}

func NewDataProcessHandler(transactionRepo repo.TransactionRepo) *DataProcessHandler {
	return &DataProcessHandler{
		TransactionRepo: transactionRepo,
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
	log.Printf("Handling DataProcessJob: Type=%s, Action=%s", job.ProcessType, job.Action)
	// Log retry info
	if n, ok := asynq.GetRetryCount(ctx); ok {
		log.Printf("Retry count (previous): %d, current attempt: %d", n, n+1)
	}
	if max, ok := asynq.GetMaxRetry(ctx); ok {
		log.Printf("Max retry configured: %d", max)
	}
	return fmt.Errorf("forced failure for testing")
	// Test lỗi: đặt Action="fail" để cố tình trả về lỗi (kích hoạt retry/Failed)
	if job.Action == "fail" {
		return fmt.Errorf("forced failure for testing")
	}
	// var transactions []model.Transaction
	// transactions, err := h.TransactionRepo.GetList(ctx)
	// log.Printf("Fetched %d transactions", len(transactions))
	// if err != nil {
	// 	return err
	// }
	// TODO: business logic xử lý theo job.ProcessType / job.Action / job.Data
	_ = job
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
