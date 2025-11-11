package handlers

import (
	"context"
	model "core-ledger/model/wealify"
	"core-ledger/pkg/queue"
	"core-ledger/pkg/queue/jobs"
	"core-ledger/pkg/repo"
	"fmt"
	"log"
)

// DataProcessHandler xử lý DataProcessJob
type MyJobHandler struct {
	repo.TransactionRepo
	// thêm dependency nếu cần (ví dụ: services, repos)
}

func NewMyJobHandler(transactionRepo repo.TransactionRepo) *MyJobHandler {
	return &MyJobHandler{
		TransactionRepo: transactionRepo,
	}
}

// NewDataProcessRegistration: provider đăng ký job/handler vào group "queue-registrations"
func NewMyJobHandlerRegistration(h *MyJobHandler) queue.Registration {
	return queue.Registration{
		Type:     "my:job",
		Template: &jobs.MyJob{},
		Handler:  h,
	}
}

func (h *MyJobHandler) Handle(ctx context.Context, j queue.Job) error {
	// kiểu assert về concrete job
	job, ok := j.(*jobs.MyJob)
	if !ok {
		return fmt.Errorf("invalid job type, expect *MyJob")
	}

	var transactions []model.Transaction
	transactions, err := h.TransactionRepo.GetList(ctx)
	log.Printf("Fetched %d transactions", len(transactions))
	if err != nil {
		return err
	}
	// TODO: business logic xử lý theo job.ProcessType / job.Action / job.Data
	_ = job
	return nil // trả về error để asynq retry nếu cần
}

// Failed: hook được gọi khi job đã hết retry hoặc timeout
func (h *MyJobHandler) Failed(ctx context.Context, j queue.Job, err error) {
	// cố gắng assert đúng loại job để log chi tiết
	if job, ok := j.(*jobs.DataProcessJob); ok {
		log.Printf("[FAILED] DataProcessJob Type=%s Action=%s Error=%v", job.ProcessType, job.Action, err)
	} else {
		log.Printf("[FAILED] DataProcessJob Error=%v", err)
	}
	// TODO: Có thể ghi log vào DB, tạo transaction_log, hoặc đẩy sang channel cảnh báo...
	_ = h // giữ chỗ nếu sau này cần dùng repo để lưu DB
}
