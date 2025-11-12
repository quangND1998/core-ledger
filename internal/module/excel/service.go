package excel

import (
	"context"
	"core-ledger/pkg/logger"
	"core-ledger/pkg/queue"
	"core-ledger/pkg/queue/jobs"
	"core-ledger/pkg/repo"
	"log"

	"gorm.io/gorm"
)

type ExcelService struct {
	db            *gorm.DB
	coAccountRepo repo.CoAccountRepo
	logger        logger.CustomLogger
	dispatcher    queue.Dispatcher
}

func NewExcelService(dispatcher queue.Dispatcher, db *gorm.DB, coAccountRepo repo.CoAccountRepo) *ExcelService {
	return &ExcelService{
		db:            db,
		coAccountRepo: coAccountRepo,
		logger:        logger.NewSystemLog("ExcelService"),
		dispatcher:    dispatcher,
	}
}

func (s *ExcelService) ImportCoAccounts(ctx context.Context, tmpFile string) error {
	s.logger.Info("Importing co-accounts from file: %s", tmpFile)
	dataJob := jobs.NewImportCoaAccount("import_coa_account", "import", jobs.DataImportCoaAccount{
		TmpFile: tmpFile,
	})
	dataJob.SetQueue("critical")
	if err := s.dispatcher.Dispatch(dataJob); err != nil {
		log.Printf("❌ Failed to dispatch data job: %v", err)

	}
	// TODO: Implement the logic to import co-accounts from the given fi
	return nil
}
