package transactions

import (
	"context"
	model "core-ledger/model/wealify"
	"core-ledger/pkg/logger"
	"core-ledger/pkg/repo"
)

type TransactionService interface {
	listTransaction(ctx context.Context) ([]model.Transaction, error)
}
type transactionService struct {
	repo.TransactionRepo
	logger logger.CustomLogger
}

func (s *transactionService) listTransaction(ctx context.Context) ([]model.Transaction, error) {
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
func NewTransactionService(transactionRepo repo.TransactionRepo) TransactionService {
	return &transactionService{
		TransactionRepo: transactionRepo,
		logger:          logger.NewSystemLog("TransactionService"),
	}
}
