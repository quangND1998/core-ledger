package repo

import (
	"context"
	model "core-ledger/model/core-ledger"
	"core-ledger/model/dto"

	"gorm.io/gorm"
)

type PermissionRepo interface {
	creator[*model.Permission]
	reader[*model.Permission, *dto.BasePaginationQuery]
	getByID[*model.Permission]
	updater[*model.Permission]
	deleter[*model.Permission]
	FindByName(ctx context.Context, name string, guardName string) (*model.Permission, error)
	FindOrCreate(ctx context.Context, name string, guardName string) (*model.Permission, error)
	GetByNames(ctx context.Context, names []string, guardName string) ([]*model.Permission, error)
}

type permissionRepo struct {
	db *gorm.DB
}

func NewPermissionRepo(db *gorm.DB) PermissionRepo {
	return &permissionRepo{db: db}
}

func (r *permissionRepo) Create(permissions ...*model.Permission) error {
	return r.db.Create(permissions).Error
}

func (r *permissionRepo) GetByID(ctx context.Context, id int64) (*model.Permission, error) {
	var permission *model.Permission
	return permission, r.db.WithContext(ctx).First(&permission, "id = ?", id).Error
}

func (r *permissionRepo) GetOneByFields(ctx context.Context, fields map[string]interface{}, preloads ...string) (*model.Permission, error) {
	var permission *model.Permission
	query := r.db.WithContext(ctx).Model(&model.Permission{})
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	return permission, query.Where(fields).First(&permission).Error
}

func (r *permissionRepo) GetManyByFields(ctx context.Context, fields map[string]interface{}, preloads ...string) ([]*model.Permission, error) {
	var permissions []*model.Permission
	query := r.db.WithContext(ctx).Model(&model.Permission{}).Where(fields)
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	return permissions, query.Find(&permissions).Error
}

func (r *permissionRepo) UpdateSelectField(entity *model.Permission, fields map[string]interface{}) error {
	return r.db.Model(entity).Updates(fields).Error
}

func (r *permissionRepo) Delete(entity *model.Permission) error {
	return r.db.Delete(entity).Error
}

func (r *permissionRepo) FindByName(ctx context.Context, name string, guardName string) (*model.Permission, error) {
	var permission *model.Permission
	err := r.db.WithContext(ctx).
		Where("name = ? AND guard_name = ?", name, guardName).
		First(&permission).Error
	return permission, err
}

func (r *permissionRepo) FindOrCreate(ctx context.Context, name string, guardName string) (*model.Permission, error) {
	var permission *model.Permission
	err := r.db.WithContext(ctx).
		Where("name = ? AND guard_name = ?", name, guardName).
		FirstOrCreate(&permission, model.Permission{
			Name:      name,
			GuardName: guardName,
		}).Error
	return permission, err
}

func (r *permissionRepo) GetByNames(ctx context.Context, names []string, guardName string) ([]*model.Permission, error) {
	var permissions []*model.Permission
	err := r.db.WithContext(ctx).
		Where("name IN ? AND guard_name = ?", names, guardName).
		Find(&permissions).Error
	return permissions, err
}

func (r *permissionRepo) Paginate(fields *dto.BasePaginationQuery) (*dto.PaginationResponse[*model.Permission], error) {
	var items []*model.Permission
	var total int64
	query := r.db.Model(&model.Permission{})

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

	return &dto.PaginationResponse[*model.Permission]{
		Items:     items,
		Total:     total,
		Limit:     int64(limit),
		Page:      int64(offset/limit + 1),
		TotalPage: totalPage,
		NextPage:  nextPage,
		PrevPage:  prevPage,
	}, nil
}

