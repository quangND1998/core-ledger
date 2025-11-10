package repo

import (
	"context"
	model "core-ledger/model/core-ledger"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type JournalRepo interface {
	creator[*model.Journal]
	// reader[*model.Journal, *dto.ListCustomerFilter]
	getByID[*model.Journal]
	updater[*model.Journal]
	Save(customer *model.Journal) error
	Upsert(accounts []*model.Journal, updateColumns []string) error
}

type journalRepo struct {
	db *gorm.DB
}

func NewJournalRepo(db *gorm.DB) JournalRepo {
	return &journalRepo{
		db: db,
	}
}
func (c *journalRepo) Save(customer *model.Journal) error {
	return c.db.Create(&customer).Error
}

func (c *journalRepo) Create(customer ...*model.Journal) error {
	return c.db.Create(customer).Error
}

func (c *journalRepo) GetByID(ctx context.Context, id int64) (*model.Journal, error) {
	customer := &model.Journal{}
	return customer, c.db.WithContext(ctx).First(&customer, "id = ?", id).Error
}

func (c *journalRepo) Update(customer *model.Journal) error {
	return c.db.Save(&customer).Error
}

func (c *journalRepo) UpdateSelectField(entity *model.Journal, fields map[string]interface{}) error {
	return c.db.Model(entity).Updates(fields).Error
}

func (c *journalRepo) Upsert(accounts []*model.Journal, updateColumns []string) error {
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
