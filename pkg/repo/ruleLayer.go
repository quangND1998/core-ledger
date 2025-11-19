package repo

import (
	"context"
	model "core-ledger/model/core-ledger"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RuleLayerRepo interface {
	creator[*model.AccountRuleLayer]
	// reader[*model.Journal, *dto.ListCustomerFilter]
	getByID[*model.AccountRuleLayer]
	updater[*model.AccountRuleLayer]
	Save(customer *model.AccountRuleLayer) error
	Upsert(accounts []*model.AccountRuleLayer, updateColumns []string) error
	List(ctx context.Context) ([]*model.AccountRuleLayer, error)
}

type ruleLayerRepo struct {
	db *gorm.DB
}

func NewRuleLayerRepo(db *gorm.DB) RuleLayerRepo {
	return &ruleLayerRepo{
		db: db,
	}
}

func (c *ruleLayerRepo) List(ctx context.Context) ([]*model.AccountRuleLayer, error) {
	var accounts []*model.AccountRuleLayer
	if err := c.db.WithContext(ctx).Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}
func (c *ruleLayerRepo) Save(customer *model.AccountRuleLayer) error {
	return c.db.Create(&customer).Error
}

func (c *ruleLayerRepo) Create(customer ...*model.AccountRuleLayer) error {
	return c.db.Create(customer).Error
}

func (c *ruleLayerRepo) GetByID(ctx context.Context, id int64) (*model.AccountRuleLayer, error) {
	customer := &model.AccountRuleLayer{}
	return customer, c.db.WithContext(ctx).First(&customer, "id = ?", id).Error
}

func (c *ruleLayerRepo) Update(customer *model.AccountRuleLayer) error {
	return c.db.Save(&customer).Error
}

func (c *ruleLayerRepo) UpdateSelectField(entity *model.AccountRuleLayer, fields map[string]interface{}) error {
	return c.db.Model(entity).Updates(fields).Error
}

func (c *ruleLayerRepo) Upsert(accounts []*model.AccountRuleLayer, updateColumns []string) error {
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
