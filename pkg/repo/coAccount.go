package repo

import (
	"context"
	model "core-ledger/model/core-ledger"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CoAccountRepo interface {
	creator[*model.CoaAccount]
	// reader[*model.CoaAccount, *dto.ListCustomerFilter]
	getByID[*model.CoaAccount]
	updater[*model.CoaAccount]
	Save(customer *model.CoaAccount) error
	Upsert(accounts []*model.CoaAccount, updateColumns []string) error
}

type coAccountRepo struct {
	db *gorm.DB
}

func NewCoAccountRepo(db *gorm.DB) CoAccountRepo {
	return &coAccountRepo{
		db: db,
	}
}
func (c *coAccountRepo) Save(customer *model.CoaAccount) error {
	return c.db.Create(&customer).Error
}

func (c *coAccountRepo) Create(customer ...*model.CoaAccount) error {
	return c.db.Create(customer).Error
}

func (c *coAccountRepo) GetByID(ctx context.Context, id int64) (*model.CoaAccount, error) {
	customer := &model.CoaAccount{}
	return customer, c.db.WithContext(ctx).First(&customer, "id = ?", id).Error
}

func (c *coAccountRepo) Update(customer *model.CoaAccount) error {
	return c.db.Save(&customer).Error
}

func (c *coAccountRepo) UpdateSelectField(entity *model.CoaAccount, fields map[string]interface{}) error {
	return c.db.Model(entity).Updates(fields).Error
}

func (c *coAccountRepo) Upsert(accounts []*model.CoaAccount, updateColumns []string) error {
	if len(accounts) == 0 {
		return nil
	}

	if len(updateColumns) == 0 {
		// Mặc định update tất cả trường có thể thay đổi
		updateColumns = []string{"name", "type", "parent_id", "status", "provider", "network", "tags", "metadata", "updated_at"}
	}

	return c.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}, {Name: "currency"}}, // cột unique
		DoUpdates: clause.AssignmentColumns(updateColumns),
	}).Create(&accounts).Error
}
