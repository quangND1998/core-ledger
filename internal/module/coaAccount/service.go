package coaaccount

import (
	"core-ledger/pkg/logger"
	"core-ledger/pkg/queue"
	"core-ledger/pkg/repo"

	"gorm.io/gorm"
)

type CoaAccountService struct {
	db            *gorm.DB
	coAccountRepo repo.CoAccountRepo
	logger        logger.CustomLogger
	dispatcher    queue.Dispatcher
}

func NewCoaAccountService(dispatcher queue.Dispatcher, db *gorm.DB, coAccountRepo repo.CoAccountRepo) *CoaAccountService {
	return &CoaAccountService{
		db:            db,
		coAccountRepo: coAccountRepo,
		logger:        logger.NewSystemLog("CoaAccountService"),
		dispatcher:    dispatcher,
	}
}
