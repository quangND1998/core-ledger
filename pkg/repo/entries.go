package repo

import (
	"context"
	model "core-ledger/model/core-ledger"
	"core-ledger/model/dto"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EnTriesRepo interface {
	creator[*model.Entry]
	// reader[*model.Entry, *dto.ListCustomerFilter]
	getByID[*model.Entry]
	updater[*model.Entry]
	Save(customer *model.Entry) error
	Upsert(accounts []*model.Entry, updateColumns []string) error
	GetByAccount(ctx context.Context, id int64) ([]model.Entry, error)
	PaginateWithScopes(ctx context.Context, filter *dto.ListEntrytFilter) (*dto.PaginationResponse[*model.Entry], error)
}

type enTriesRepo struct {
	db *gorm.DB
}

func NewEnTriesRepo(db *gorm.DB) EnTriesRepo {
	return &enTriesRepo{
		db: db,
	}
}
func (c *enTriesRepo) Save(customer *model.Entry) error {
	return c.db.Create(&customer).Error
}

func (c *enTriesRepo) Create(customer ...*model.Entry) error {
	return c.db.Create(customer).Error
}

func (c *enTriesRepo) GetByID(ctx context.Context, id int64) (*model.Entry, error) {
	customer := &model.Entry{}
	return customer, c.db.WithContext(ctx).First(&customer, "id = ?", id).Error
}

func (c *enTriesRepo) Update(customer *model.Entry) error {
	return c.db.Save(&customer).Error
}

func (c *enTriesRepo) UpdateSelectField(entity *model.Entry, fields map[string]interface{}) error {
	return c.db.Model(entity).Updates(fields).Error
}

func (c *enTriesRepo) Upsert(accounts []*model.Entry, updateColumns []string) error {
	if len(accounts) == 0 {
		return nil
	}

	if len(updateColumns) == 0 {
		// Mặc định update tất cả trường có thể thay đổi
		updateColumns = []string{"amount", "status", "Meta", "updated_at"}
	}

	return c.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}, {Name: "currency"}}, // cột unique
		DoUpdates: clause.AssignmentColumns(updateColumns),
	}).Create(&accounts).Error
}

func (c *enTriesRepo) GetByAccount(context context.Context, id int64) ([]model.Entry, error) {
	entries := []model.Entry{}
	return entries, c.db.WithContext(context).Where("account_id = ?", id).Find(&entries).Error
}

func (r *enTriesRepo) PaginateWithScopes(ctx context.Context, fields *dto.ListEntrytFilter) (*dto.PaginationResponse[*model.Entry], error) {
	params := BuildParamsFromFilter(fields)

	var items []*model.Entry

	// q = q.
	// 	Preload("Entries").Preload("Parent").
	// 	Preload("Children")
	limit := int64(25)
	page := int64(1)
	if fields.Limit != nil {
		limit = *fields.Limit
	}
	if fields.Page != nil {
		page = *fields.Page
	}

	pagination, err := CustomPaginate(r.db.Model(&model.CoaAccount{}), params, page, limit, &items)
	if err != nil {
		return nil, err
	}

	return pagination, nil
}
