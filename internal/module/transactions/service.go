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
	job := jobs.NewDataProcessJob("test", "fail", map[string]interface{}{
		"test": "backoff",
	})
	job.SetBackoff([]int{2, 5, 10}) // Custom backoff: 2s, 5s, 10s
	job.SetRetry(3)                 // Cho phép retry 3 lần

	if err := s.dispatcher.Dispatch(job, queue.Timeout(1*time.Second)); err != nil {
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
