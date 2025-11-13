package coaaccount

import (
	"context"
	model "core-ledger/model/core-ledger"
	"core-ledger/model/dto"
	"core-ledger/pkg/logger"
	"core-ledger/pkg/queue"
	"core-ledger/pkg/repo"

	"gorm.io/gorm"
)

type CoaAccountService struct {
	db            *gorm.DB
	coAccountRepo repo.CoAccountRepo
	entriesRepo   repo.EnTriesRepo
	snapShotRepo  repo.SnapshotRepo
	logger        logger.CustomLogger
	dispatcher    queue.Dispatcher
}

func NewCoaAccountService(dispatcher queue.Dispatcher, db *gorm.DB, coAccountRepo repo.CoAccountRepo, entriesRepo repo.EnTriesRepo, snapShotRepo repo.SnapshotRepo) *CoaAccountService {
	return &CoaAccountService{
		db:            db,
		coAccountRepo: coAccountRepo,
		logger:        logger.NewSystemLog("CoaAccountService"),
		dispatcher:    dispatcher,
		entriesRepo:   entriesRepo,
		snapShotRepo:  snapShotRepo,
	}
}

func (c *CoaAccountService) GetCoaAccountDetail(ctx context.Context, id int64) (*dto.CoaAccountDetailResponse, error) {
	data := &dto.CoaAccountDetailResponse{
		CoaAccount: nil,
		Entries:    []model.Entry{},
		Snapshots:  []model.Snapshot{},
	}

	var err error
	if data.CoaAccount, err = c.coAccountRepo.GetByID(ctx, id); err != nil {
		return data, err
	}
	if data.Entries, err = c.entriesRepo.GetByAccount(ctx, id); err != nil {
		return data, err
	}
	if data.Snapshots, err = c.snapShotRepo.GetByAccount(ctx, id); err != nil {
		return data, err
	}

	return data, nil
}
