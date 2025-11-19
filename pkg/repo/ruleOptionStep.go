package repo

import (
	"context"
	model "core-ledger/model/core-ledger"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RuleOptionStepRepo interface {
	creator[*model.AccountRuleOptionStep]
	// reader[*model.Journal, *dto.ListCustomerFilter]
	getByID[*model.AccountRuleOptionStep]
	updater[*model.AccountRuleOptionStep]
	Save(customer *model.AccountRuleOptionStep) error
	Upsert(accounts []*model.AccountRuleOptionStep, updateColumns []string) error
	List(ctx context.Context) ([]*model.AccountRuleOptionStep, error)
}

type ruleOptionStepRepo struct {
	db *gorm.DB
}

func NewRuleOptionStepRepo(db *gorm.DB) RuleOptionStepRepo {
	return &ruleOptionStepRepo{
		db: db,
	}
}

func (c *ruleOptionStepRepo) List(ctx context.Context) ([]*model.AccountRuleOptionStep, error) {
	var accounts []*model.AccountRuleOptionStep
	if err := c.db.WithContext(ctx).Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}
func (c *ruleOptionStepRepo) Save(customer *model.AccountRuleOptionStep) error {
	return c.db.Create(&customer).Error
}

func (c *ruleOptionStepRepo) Create(customer ...*model.AccountRuleOptionStep) error {
	return c.db.Create(customer).Error
}

func (c *ruleOptionStepRepo) GetByID(ctx context.Context, id int64) (*model.AccountRuleOptionStep, error) {
	customer := &model.AccountRuleOptionStep{}
	return customer, c.db.WithContext(ctx).First(&customer, "id = ?", id).Error
}

func (c *ruleOptionStepRepo) Update(customer *model.AccountRuleOptionStep) error {
	return c.db.Save(&customer).Error
}

func (c *ruleOptionStepRepo) UpdateSelectField(entity *model.AccountRuleOptionStep, fields map[string]interface{}) error {
	return c.db.Model(entity).Updates(fields).Error
}

func (c *ruleOptionStepRepo) Upsert(accounts []*model.AccountRuleOptionStep, updateColumns []string) error {
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
