package repo

import (
	"context"
	model "core-ledger/model/core-ledger"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RuleValueRepo interface {
	creator[*model.RuleValue]
	// reader[*model.Journal, *dto.ListCustomerFilter]
	getByID[*model.RuleValue]
	updater[*model.RuleValue]
	Save(customer *model.RuleValue) error
	Upsert(accounts []*model.RuleValue, updateColumns []string) error
	List(ctx context.Context) ([]*model.RuleValue, error)
}

type ruleValueRepo struct {
	db *gorm.DB
}

func NewRuleValueRepo(db *gorm.DB) RuleValueRepo {
	return &ruleValueRepo{
		db: db,
	}
}

func (c *ruleValueRepo) List(ctx context.Context) ([]*model.RuleValue, error) {
	var accounts []*model.RuleValue
	if err := c.db.WithContext(ctx).Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}
func (c *ruleValueRepo) Save(customer *model.RuleValue) error {
	return c.db.Create(&customer).Error
}

func (c *ruleValueRepo) Create(customer ...*model.RuleValue) error {
	return c.db.Create(customer).Error
}

func (c *ruleValueRepo) GetByID(ctx context.Context, id int64) (*model.RuleValue, error) {
	customer := &model.RuleValue{}
	return customer, c.db.WithContext(ctx).First(&customer, "id = ?", id).Error
}

func (c *ruleValueRepo) Update(customer *model.RuleValue) error {
	return c.db.Save(&customer).Error
}

func (c *ruleValueRepo) UpdateSelectField(entity *model.RuleValue, fields map[string]interface{}) error {
	return c.db.Model(entity).Updates(fields).Error
}

func (c *ruleValueRepo) Upsert(accounts []*model.RuleValue, updateColumns []string) error {
	if len(accounts) == 0 {
		return nil
	}

	if len(updateColumns) == 0 {
		// Mặc định update tất cả trường có thể thay đổi
		updateColumns = []string{"name", "code", "updated_at"}
	}

	return c.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}}, // cột unique
		DoUpdates: clause.AssignmentColumns(updateColumns),
	}).Create(&accounts).Error
}
