package permission

import (
	"core-ledger/pkg/logger"
	"core-ledger/pkg/queue"
	"core-ledger/pkg/repo"

	"gorm.io/gorm"
)

type PermissionService struct {
	db            *gorm.DB
	permissionRepo repo.PermissionRepo
	logger         logger.CustomLogger
	dispatcher     queue.Dispatcher
}

func NewPermissionService(dispatcher queue.Dispatcher, db *gorm.DB, permissionRepo repo.PermissionRepo) *PermissionService {
	return &PermissionService{
		db:            db,
		permissionRepo: permissionRepo,
		logger:         logger.NewSystemLog("PermissionService"),
		dispatcher:     dispatcher,
	}
}
