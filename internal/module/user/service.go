package user

import (
	"context"
	model "core-ledger/model/core-ledger"
	"core-ledger/pkg/logger"
	"core-ledger/pkg/queue"
	"core-ledger/pkg/repo"

	"gorm.io/gorm"
)

type UserService struct {
	db            *gorm.DB
	userRepo       repo.UserRepo
	roleRepo       repo.RoleRepo
	permissionRepo repo.PermissionRepo
	logger         logger.CustomLogger
	dispatcher     queue.Dispatcher
}

func NewUserService(dispatcher queue.Dispatcher, db *gorm.DB, userRepo repo.UserRepo, roleRepo repo.RoleRepo, permissionRepo repo.PermissionRepo) *UserService {
	return &UserService{
		db:            db,
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		logger:         logger.NewSystemLog("UserService"),
		dispatcher:     dispatcher,
	}
}

// User permission methods
func (s *UserService) GivePermissionToUser(ctx context.Context, userID uint64, permissionName string, guardName string) error {
	user, err := s.userRepo.GetByID(ctx, int64(userID))
	if err != nil {
		return err
	}
	return user.GivePermissionTo(s.db, permissionName, guardName)
}

func (s *UserService) RevokePermissionFromUser(ctx context.Context, userID uint64, permissionName string, guardName string) error {
	user, err := s.userRepo.GetByID(ctx, int64(userID))
	if err != nil {
		return err
	}
	return user.RevokePermissionTo(s.db, permissionName, guardName)
}

func (s *UserService) AssignRoleToUser(ctx context.Context, userID uint64, roleName string, guardName string) error {
	user, err := s.userRepo.GetByID(ctx, int64(userID))
	if err != nil {
		return err
	}
	return user.AssignRole(s.db, roleName, guardName)
}

func (s *UserService) RemoveRoleFromUser(ctx context.Context, userID uint64, roleName string, guardName string) error {
	user, err := s.userRepo.GetByID(ctx, int64(userID))
	if err != nil {
		return err
	}
	return user.RemoveRole(s.db, roleName, guardName)
}

func (s *UserService) SyncUserPermissions(ctx context.Context, userID uint64, permissionNames []string, guardName string) error {
	user, err := s.userRepo.GetByID(ctx, int64(userID))
	if err != nil {
		return err
	}
	return user.SyncPermissions(s.db, permissionNames, guardName)
}

func (s *UserService) SyncUserRoles(ctx context.Context, userID uint64, roleNames []string, guardName string) error {
	user, err := s.userRepo.GetByID(ctx, int64(userID))
	if err != nil {
		return err
	}
	return user.SyncRoles(s.db, roleNames, guardName)
}

func (s *UserService) SyncUserRolesByIDs(ctx context.Context, userID uint64, roleIDs []uint64) error {
	user, err := s.userRepo.GetByID(ctx, int64(userID))
	if err != nil {
		return err
	}
	return user.SyncRolesByIDs(s.db, roleIDs)
}

func (s *UserService) UserHasPermission(ctx context.Context, userID uint64, permissionName string, guardName string) (bool, error) {
	user, err := s.userRepo.GetByID(ctx, int64(userID))
	if err != nil {
		return false, err
	}
	return user.HasPermission(s.db, permissionName, guardName)
}

func (s *UserService) UserHasRole(ctx context.Context, userID uint64, roleName string, guardName string) (bool, error) {
	user, err := s.userRepo.GetByID(ctx, int64(userID))
	if err != nil {
		return false, err
	}
	return user.HasRole(s.db, roleName, guardName)
}

func (s *UserService) UserHasAnyPermission(ctx context.Context, userID uint64, permissionNames []string, guardName string) (bool, error) {
	user, err := s.userRepo.GetByID(ctx, int64(userID))
	if err != nil {
		return false, err
	}
	return user.HasAnyPermission(s.db, permissionNames, guardName)
}

func (s *UserService) UserHasAllPermissions(ctx context.Context, userID uint64, permissionNames []string, guardName string) (bool, error) {
	user, err := s.userRepo.GetByID(ctx, int64(userID))
	if err != nil {
		return false, err
	}
	return user.HasAllPermissions(s.db, permissionNames, guardName)
}

func (s *UserService) GetAllUserPermissions(ctx context.Context, userID uint64, guardName string) ([]model.Permission, error) {
	user, err := s.userRepo.GetByID(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	return user.GetAllPermissions(s.db, guardName)
}

func (s *UserService) GetUserRoleNames(ctx context.Context, userID uint64, guardName string) ([]string, error) {
	user, err := s.userRepo.GetByID(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	return user.GetRoleNames(s.db, guardName)
}

func (s *UserService) GetUserPermissionNames(ctx context.Context, userID uint64, guardName string) ([]string, error) {
	user, err := s.userRepo.GetByID(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	return user.GetPermissionNames(s.db, guardName)
}

