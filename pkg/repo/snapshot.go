package repo

import (
	"context"
	model "core-ledger/model/core-ledger"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SnapshotRepo interface {
	creator[*model.Snapshot]
	// reader[*model.Snapshot, *dto.ListCustomerFilter]
	getByID[*model.Snapshot]
	updater[*model.Snapshot]
	Save(customer *model.Snapshot) error
	Upsert(accounts []*model.Snapshot, updateColumns []string) error
}

type snapShotRepo struct {
	db *gorm.DB
}

func NewSnapshotRepo(db *gorm.DB) SnapshotRepo {
	return &snapShotRepo{
		db: db,
	}
}
func (c *snapShotRepo) Save(customer *model.Snapshot) error {
	return c.db.Create(&customer).Error
}

func (c *snapShotRepo) Create(customer ...*model.Snapshot) error {
	return c.db.Create(customer).Error
}

func (c *snapShotRepo) GetByID(ctx context.Context, id int64) (*model.Snapshot, error) {
	customer := &model.Snapshot{}
	return customer, c.db.WithContext(ctx).First(&customer, "id = ?", id).Error
}

func (c *snapShotRepo) Update(customer *model.Snapshot) error {
	return c.db.Save(&customer).Error
}

func (c *snapShotRepo) UpdateSelectField(entity *model.Snapshot, fields map[string]interface{}) error {
	return c.db.Model(entity).Updates(fields).Error
}

func (c *snapShotRepo) Upsert(accounts []*model.Snapshot, updateColumns []string) error {
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
