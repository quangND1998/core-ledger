package excel

import (
	"core-ledger/pkg/logger"
	"core-ledger/pkg/queue"
)

type ExcelService interface {
}
type excelService struct {
	logger     logger.CustomLogger
	dispatcher queue.Dispatcher
}

func NewExcelService(dispatcher queue.Dispatcher) *excelService {
	return &excelService{
		logger:     logger.NewSystemLog("ExcelService"),
		dispatcher: dispatcher,
	}
}

func (s *excelService) ImportCoAccounts(tmpFile string) error {
	s.logger.Info("Importing co-accounts from file: %s", tmpFile)
	
	// TODO: Implement the logic to import co-accounts from the given fi
	return nil
}
