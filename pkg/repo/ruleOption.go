package repo

import (
	"context"
	model "core-ledger/model/core-ledger"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RuleOptionRepo interface {
	creator[*model.AccountRuleOption]
	// reader[*model.Journal, *dto.ListCustomerFilter]
	getByID[*model.AccountRuleOption]
	updater[*model.AccountRuleOption]
	Save(customer *model.AccountRuleOption) error
	Upsert(accounts []*model.AccountRuleOption, updateColumns []string) error
	List(ctx context.Context) ([]*model.AccountRuleOption, error)
	GetOneByFields(ctx context.Context, fields map[string]interface{}, preloads ...string) (*model.AccountRuleOption, error)
	GetManyByFields(ctx context.Context, fields map[string]interface{}, preloads ...string) ([]*model.AccountRuleOption, error)
}

type ruleOptionRepo struct {
	db *gorm.DB
}

func NewRuleOptionRepo(db *gorm.DB) RuleOptionRepo {
	return &ruleOptionRepo{
		db: db,
	}
}

func (c *ruleOptionRepo) List(ctx context.Context) ([]*model.AccountRuleOption, error) {
	var accounts []*model.AccountRuleOption
	if err := c.db.WithContext(ctx).Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}
func (c *ruleOptionRepo) Save(customer *model.AccountRuleOption) error {
	return c.db.Create(&customer).Error
}

func (c *ruleOptionRepo) Create(customer ...*model.AccountRuleOption) error {
	return c.db.Create(customer).Error
}

func (c *ruleOptionRepo) GetByID(ctx context.Context, id int64) (*model.AccountRuleOption, error) {
	customer := &model.AccountRuleOption{}
	return customer, c.db.WithContext(ctx).First(&customer, "id = ?", id).Error
}

func (c *ruleOptionRepo) Update(customer *model.AccountRuleOption) error {
	return c.db.Save(&customer).Error
}

func (c *ruleOptionRepo) UpdateSelectField(entity *model.AccountRuleOption, fields map[string]interface{}) error {
	return c.db.Model(entity).Updates(fields).Error
}

func (c *ruleOptionRepo) Upsert(accounts []*model.AccountRuleOption, updateColumns []string) error {
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
func (s *ruleOptionRepo) GetOneByFields(ctx context.Context, fields map[string]interface{}, preloads ...string) (*model.AccountRuleOption, error) {
	var ruleOption *model.AccountRuleOption
	query := s.db.WithContext(ctx).Model(&model.AccountRuleOption{})
	for _, preload := range preloads {
		query.Joins(preload).Preload(preload)
	}
	return ruleOption, query.Where(fields).First(&ruleOption).Error
}

func (s *ruleOptionRepo) GetManyByFields(ctx context.Context, fields map[string]interface{}, preloads ...string) ([]*model.AccountRuleOption, error) {
	var ruleOptions []*model.AccountRuleOption
	return ruleOptions, s.db.WithContext(ctx).Model(&model.AccountRuleOption{}).Where(fields).Find(&ruleOptions).Error
}
