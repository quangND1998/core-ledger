package repo

import (
	"context"
	model "core-ledger/model/core-ledger"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RuleCategoryRepo interface {
	creator[*model.RuleCategory]
	// reader[*model.Journal, *dto.ListCustomerFilter]
	getByID[*model.RuleCategory]
	updater[*model.RuleCategory]
	Save(customer *model.RuleCategory) error
	Upsert(accounts []*model.RuleCategory, updateColumns []string) error
	List(ctx context.Context) ([]*model.RuleCategory, error)
}

type ruleCategoryRepo struct {
	db *gorm.DB
}

func NewRuleCategoryRepo(db *gorm.DB) RuleCategoryRepo {
	return &ruleCategoryRepo{
		db: db,
	}
}

func (c *ruleCategoryRepo) List(ctx context.Context) ([]*model.RuleCategory, error) {
	var accounts []*model.RuleCategory
	if err := c.db.WithContext(ctx).Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}
func (c *ruleCategoryRepo) Save(customer *model.RuleCategory) error {
	return c.db.Create(&customer).Error
}

func (c *ruleCategoryRepo) Create(customer ...*model.RuleCategory) error {
	return c.db.Create(customer).Error
}

func (c *ruleCategoryRepo) GetByID(ctx context.Context, id int64) (*model.RuleCategory, error) {
	customer := &model.RuleCategory{}
	return customer, c.db.WithContext(ctx).First(&customer, "id = ?", id).Error
}

func (c *ruleCategoryRepo) Update(customer *model.RuleCategory) error {
	return c.db.Save(&customer).Error
}

func (c *ruleCategoryRepo) UpdateSelectField(entity *model.RuleCategory, fields map[string]interface{}) error {
	return c.db.Model(entity).Updates(fields).Error
}

func (c *ruleCategoryRepo) Upsert(accounts []*model.RuleCategory, updateColumns []string) error {
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
