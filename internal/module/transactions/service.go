package transactions

import (
	"context"
	model "core-ledger/model/wealify"
	"core-ledger/pkg/logger"
	"core-ledger/pkg/queue"
	"core-ledger/pkg/queue/jobs"
	"core-ledger/pkg/repo"
	"log"
	"time"
)

type TransactionService interface {
	listTransaction(ctx context.Context) ([]model.Transaction, error)
}
type transactionService struct {
	repo.TransactionRepo
	userRepo   repo.UserRepo
	logger     logger.CustomLogger
	dispatcher queue.Dispatcher
}

func (s *transactionService) ListTransaction(ctx context.Context) ([]model.Transaction, error) {
	var transactions []model.Transaction
	transactions, err := s.TransactionRepo.GetList(ctx)
	s.logger.Info("Fetched %d transactions ", len(transactions))
	if err != nil {
		return nil, err
	}
	dataJob := jobs.NewDataProcessJob("user_analytics", "export", map[string]interface{}{
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
	dataJob.SetBackoff([]int{10, 20, 30})
	if err := s.dispatcher.Dispatch(dataJob, queue.Timeout(1*time.Second)); err != nil {
		log.Printf("❌ Failed to dispatch data job: %v", err)
		return nil, err
	}
	// TODO: Implement the logic to fetch transactions from the database based on the request parameters
	// and populate the transactions slice accordingly.
	return transactions, nil
}

func (s *transactionService) List(ctx context.Context) ([]model.Transaction, error) {
	var transactions []model.Transaction
	transactions, err := s.TransactionRepo.GetList(ctx)
	s.logger.Info("Fetched %d transactions", len(transactions))
	if err != nil {
		return nil, err
	}

	// TODO: Implement the logic to fetch transactions from the database based on the request parameters
	// and populate the transactions slice accordingly.
	return transactions, nil
}
func NewTransactionService(transactionRepo repo.TransactionRepo, userRepo repo.UserRepo, dispatcher queue.Dispatcher) *transactionService {
	return &transactionService{
		TransactionRepo: transactionRepo,
		userRepo:        userRepo,
		logger:          logger.NewSystemLog("TransactionService"),
		dispatcher:      dispatcher,
	}
}
