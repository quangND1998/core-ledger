package role

import (
	"context"
	"core-ledger/pkg/logger"
	"core-ledger/pkg/queue"
	"core-ledger/pkg/repo"

	"gorm.io/gorm"
)

type RoleService struct {
	db            *gorm.DB
	roleRepo       repo.RoleRepo
	permissionRepo repo.PermissionRepo
	logger         logger.CustomLogger
	dispatcher     queue.Dispatcher
}

func NewRoleService(dispatcher queue.Dispatcher, db *gorm.DB, roleRepo repo.RoleRepo, permissionRepo repo.PermissionRepo) *RoleService {
	return &RoleService{
		db:            db,
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		logger:         logger.NewSystemLog("RoleService"),
		dispatcher:     dispatcher,
	}
}

// Role permission methods
func (s *RoleService) GivePermissionToRole(ctx context.Context, roleID uint64, permissionName string, guardName string) error {
	role, err := s.roleRepo.GetByID(ctx, int64(roleID))
	if err != nil {
		return err
	}
	return role.GivePermissionTo(s.db, permissionName, guardName)
}

func (s *RoleService) RevokePermissionFromRole(ctx context.Context, roleID uint64, permissionName string, guardName string) error {
	role, err := s.roleRepo.GetByID(ctx, int64(roleID))
	if err != nil {
		return err
	}
	return role.RevokePermissionTo(s.db, permissionName, guardName)
}

func (s *RoleService) SyncRolePermissions(ctx context.Context, roleID uint64, permissionNames []string, guardName string) error {
	role, err := s.roleRepo.GetByID(ctx, int64(roleID))
	if err != nil {
		return err
	}
	return role.SyncPermissions(s.db, permissionNames, guardName)
}

func (s *RoleService) SyncRolePermissionsByIDs(ctx context.Context, roleID uint64, permissionIDs []uint64) error {
	return s.roleRepo.SyncPermissions(ctx, roleID, permissionIDs)
}

