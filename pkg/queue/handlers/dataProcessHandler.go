package handlers

import (
	"context"
	"core-ledger/pkg/queue"
	"core-ledger/pkg/repo"
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

func (h *DataProcessHandler) Handle(ctx context.Context, j queue.Job) error {
	// kiểu assert về concrete job

	return nil
}
