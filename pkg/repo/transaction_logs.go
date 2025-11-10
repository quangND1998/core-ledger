package repo

import (
	"context"
	model "core-ledger/model/core-ledger"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TransactionLogRepo interface {
	creator[*model.TransactionLog]
	// reader[*model.TransactionLog, *dto.ListCustomerFilter]
	getByID[*model.TransactionLog]
	updater[*model.TransactionLog]
	Save(customer *model.TransactionLog) error
	Upsert(accounts []*model.TransactionLog, updateColumns []string) error
}

type transactionLogRepo struct {
	db *gorm.DB
}

func NewTransactionLogRepo(db *gorm.DB) TransactionLogRepo {
	return &transactionLogRepo{
		db: db,
	}
}
func (c *transactionLogRepo) Save(customer *model.TransactionLog) error {
	return c.db.Create(&customer).Error
}

func (c *transactionLogRepo) Create(customer ...*model.TransactionLog) error {
	return c.db.Create(customer).Error
}

func (c *transactionLogRepo) GetByID(ctx context.Context, id int64) (*model.TransactionLog, error) {
	customer := &model.TransactionLog{}
	return customer, c.db.WithContext(ctx).First(&customer, "id = ?", id).Error
}

func (c *transactionLogRepo) Update(customer *model.TransactionLog) error {
	return c.db.Save(&customer).Error
}

func (c *transactionLogRepo) UpdateSelectField(entity *model.TransactionLog, fields map[string]interface{}) error {
	return c.db.Model(entity).Updates(fields).Error
}

func (c *transactionLogRepo) Upsert(accounts []*model.TransactionLog, updateColumns []string) error {
	if len(accounts) == 0 {
		return nil
	}

	if len(updateColumns) == 0 {
		// Mặc định update tất cả trường có thể thay đổi
		updateColumns = []string{"status", "posted_at", "Meta", "updated_at"}
	}

	return c.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}, {Name: "currency"}}, // cột unique
		DoUpdates: clause.AssignmentColumns(updateColumns),
	}).Create(&accounts).Error
}
