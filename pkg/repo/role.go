package repo

import (
	"context"
	model "core-ledger/model/core-ledger"
	"core-ledger/model/dto"

	"gorm.io/gorm"
)

type RoleRepo interface {
	creator[*model.Role]
	reader[*model.Role, *dto.BasePaginationQuery]
	getByID[*model.Role]
	updater[*model.Role]
	deleter[*model.Role]
	FindByName(ctx context.Context, name string, guardName string) (*model.Role, error)
	FindOrCreate(ctx context.Context, name string, guardName string) (*model.Role, error)
	GetByNames(ctx context.Context, names []string, guardName string) ([]*model.Role, error)
	GivePermissionTo(ctx context.Context, roleID uint64, permissionID uint64) error
	RevokePermissionTo(ctx context.Context, roleID uint64, permissionID uint64) error
	SyncPermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error
	GetPermissions(ctx context.Context, roleID uint64) ([]*model.Permission, error)
}

type roleRepo struct {
	db *gorm.DB
}

func NewRoleRepo(db *gorm.DB) RoleRepo {
	return &roleRepo{db: db}
}

func (r *roleRepo) Create(roles ...*model.Role) error {
	return r.db.Create(roles).Error
}

func (r *roleRepo) GetByID(ctx context.Context, id int64) (*model.Role, error) {
	var role *model.Role
	return role, r.db.WithContext(ctx).Preload("Permissions").First(&role, "id = ?", id).Error
}

func (r *roleRepo) GetOneByFields(ctx context.Context, fields map[string]interface{}, preloads ...string) (*model.Role, error) {
	var role *model.Role
	query := r.db.WithContext(ctx).Model(&model.Role{})
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	return role, query.Where(fields).First(&role).Error
}

func (r *roleRepo) GetManyByFields(ctx context.Context, fields map[string]interface{}, preloads ...string) ([]*model.Role, error) {
	var roles []*model.Role
	query := r.db.WithContext(ctx).Model(&model.Role{}).Where(fields)
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	return roles, query.Find(&roles).Error
}

func (r *roleRepo) UpdateSelectField(entity *model.Role, fields map[string]interface{}) error {
	return r.db.Model(entity).Updates(fields).Error
}

func (r *roleRepo) Delete(entity *model.Role) error {
	return r.db.Delete(entity).Error
}

func (r *roleRepo) FindByName(ctx context.Context, name string, guardName string) (*model.Role, error) {
	var role *model.Role
	err := r.db.WithContext(ctx).
		Where("name = ? AND guard_name = ?", name, guardName).
		First(&role).Error
	return role, err
}

func (r *roleRepo) FindOrCreate(ctx context.Context, name string, guardName string) (*model.Role, error) {
	var role *model.Role
	err := r.db.WithContext(ctx).
		Where("name = ? AND guard_name = ?", name, guardName).
		FirstOrCreate(&role, model.Role{
			Name:      name,
			GuardName: guardName,
		}).Error
	return role, err
}

func (r *roleRepo) GetByNames(ctx context.Context, names []string, guardName string) ([]*model.Role, error) {
	var roles []*model.Role
	err := r.db.WithContext(ctx).
		Where("name IN ? AND guard_name = ?", names, guardName).
		Find(&roles).Error
	return roles, err
}

func (r *roleRepo) GivePermissionTo(ctx context.Context, roleID uint64, permissionID uint64) error {
	// Check if already exists
	var existing model.RoleHasPermission
	err := r.db.WithContext(ctx).
		Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		First(&existing).Error
	if err == nil {
		// Already assigned
		return nil
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	// Assign permission
	roleHasPermission := model.RoleHasPermission{
		RoleID:       roleID,
		PermissionID: permissionID,
	}
	return r.db.WithContext(ctx).Create(&roleHasPermission).Error
}

func (r *roleRepo) RevokePermissionTo(ctx context.Context, roleID uint64, permissionID uint64) error {
	return r.db.WithContext(ctx).
		Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Delete(&model.RoleHasPermission{}).Error
}

func (r *roleRepo) SyncPermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error {
	// Remove all existing permissions
	err := r.db.WithContext(ctx).
		Where("role_id = ?", roleID).
		Delete(&model.RoleHasPermission{}).Error
	if err != nil {
		return err
	}

	// Assign new permissions
	if len(permissionIDs) > 0 {
		roleHasPermissions := make([]*model.RoleHasPermission, 0, len(permissionIDs))
		for _, permID := range permissionIDs {
			roleHasPermissions = append(roleHasPermissions, &model.RoleHasPermission{
				RoleID:       roleID,
				PermissionID: permID,
			})
		}
		return r.db.WithContext(ctx).Create(&roleHasPermissions).Error
	}

	return nil
}

func (r *roleRepo) GetPermissions(ctx context.Context, roleID uint64) ([]*model.Permission, error) {
	var permissions []*model.Permission
	err := r.db.WithContext(ctx).
		Joins("JOIN role_has_permissions ON permissions.id = role_has_permissions.permission_id").
		Where("role_has_permissions.role_id = ?", roleID).
		Find(&permissions).Error
	return permissions, err
}

func (r *roleRepo) Paginate(fields *dto.BasePaginationQuery) (*dto.PaginationResponse[*model.Role], error) {
	var items []*model.Role
	var total int64
	query := r.db.Model(&model.Role{}).Preload("Permissions")

	err := query.Count(&total).Error
	if err != nil {
		return nil, err
	}

	limit := 25
	offset := 0
	var page int64 = 1
	if fields != nil {
		if fields.Limit != nil {
			limit = int(*fields.Limit)
		}
		if fields.Page != nil {
			offset = int(*fields.Page-1) * limit
			page = *fields.Page
		}
	}

	err = query.Limit(limit).Offset(offset).Find(&items).Error
	if err != nil {
		return nil, err
	}

	totalPage := total/int64(limit) + 1
	var nextPage *int64
	var prevPage *int64
	if page < totalPage {
		nextPage = &[]int64{page + 1}[0]
	}
	if page > 1 {
		prevPage = &[]int64{page - 1}[0]
	}

	return &dto.PaginationResponse[*model.Role]{
		Items:     items,
		Total:     total,
		Limit:     int64(limit),
		Page:      int64(offset/limit + 1),
		TotalPage: totalPage,
		NextPage:  nextPage,
		PrevPage:  prevPage,
	}, nil
}

