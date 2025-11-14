package entries

import (
	"core-ledger/pkg/logger"
	"core-ledger/pkg/queue"
	"core-ledger/pkg/repo"

	"gorm.io/gorm"
)

type EntriesService struct {
	db            *gorm.DB
	coAccountRepo repo.CoAccountRepo
	entriesRepo   repo.EnTriesRepo
	snapShotRepo  repo.SnapshotRepo
	logger        logger.CustomLogger
	dispatcher    queue.Dispatcher
}

func NewEntriesService(dispatcher queue.Dispatcher, db *gorm.DB, entriesRepo repo.EnTriesRepo) *EntriesService {
	return &EntriesService{
		db: db,
		logger:      logger.NewSystemLog("EntriesService"),
		dispatcher:  dispatcher,
		entriesRepo: entriesRepo,
	}
}
