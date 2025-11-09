package repo

import (
	"context"
	"core-ledger/model/dto"
	model "core-ledger/model/wealify"
	"fmt"

	"gorm.io/gorm"
)

type CustomerRepo interface {
	creator[*model.Customer]
	reader[*model.Customer, *dto.ListCustomerFilter]
	getByID[*model.Customer]
	updater[*model.Customer]
	Save(customer *model.Customer) error
	GetByIds(ids []int64) ([]*model.Customer, error)
	Count(fields map[string]interface{}) (int64, error)
	GetLastId(ctx context.Context) (int64, error)
	GetAvailableBanks(ctx context.Context, customerID int64) ([]*model.CustomerAvailableBank, error)
	GetAvailablePlatforms(ctx context.Context, customerID int64) ([]*model.CustomerPlatform, error)
	GetByApiKey(ctx context.Context, apkKey string) (*model.Customer, error)
	FindCustomerByEmail(ctx context.Context, email string) (*model.Customer, error)

	FindCustomerWithoutKYB(ctx context.Context, limit, offset int) ([]*model.Customer, error)
	UpdateKYBStatus(ctx context.Context, ids []int64, status int) error
	GetTotalWalletsBalance(ctx context.Context, customerIDs []int64, walletType model.WalletType) (float64, error)
}

type customerRepo struct {
	db *gorm.DB
}

func (s *customerRepo) FindCustomerWithoutKYB(ctx context.Context, limit, offset int) ([]*model.Customer, error) {
	var customers []*model.Customer
	err := s.db.WithContext(ctx).Model(&model.Customer{}).Where("kyb_status = ?", 0).Limit(limit).Offset(offset).Find(&customers).Error
	if err != nil {
		return nil, err
	}
	return customers, nil
}
func (s *customerRepo) UpdateKYBStatus(ctx context.Context, ids []int64, status int) error {
	return s.db.WithContext(ctx).Exec("UPDATE customers c SET kyb_status = ? WHERE c.id IN ?", status, ids).Error
}

func NewCustomerRepo(db *gorm.DB) CustomerRepo {
	return &customerRepo{db: db}
}

func (s *customerRepo) Create(customer ...*model.Customer) error {
	return s.db.Create(customer).Error
}

func (s *customerRepo) GetByID(ctx context.Context, id int64) (*model.Customer, error) {
	customer := &model.Customer{}
	return customer, s.db.WithContext(ctx).First(&customer, "id = ?", id).Error
}

func (s *customerRepo) Save(customer *model.Customer) error {
	return s.db.Save(customer).Error
}

func (s *customerRepo) GetByIds(ids []int64) ([]*model.Customer, error) {
	var sessions []*model.Customer
	return sessions, s.db.Model(&model.Session{}).Where("id in (?)", ids).Find(&sessions).Error
}

func (s *customerRepo) GetOneByFields(ctx context.Context, fields map[string]interface{}, preloads ...string) (*model.Customer, error) {
	var customer *model.Customer
	query := s.db.WithContext(ctx).Model(model.Customer{})
	for _, preload := range preloads {
		query.Preload(preload)
	}
	return customer, query.Where(fields).First(&customer).Error
}

func (s *customerRepo) GetManyByFields(ctx context.Context, fields map[string]interface{}, preloads ...string) ([]*model.Customer, error) {
	fmt.Println("implement me")
	return []*model.Customer{}, nil
}

func (s *customerRepo) Count(fields map[string]interface{}) (int64, error) {
	var total int64
	query := s.db.Model((*model.Customer)(nil))
	for key, value := range fields {
		query.Where(fmt.Sprintf("%s = ?", key), value)
	}
	err := query.Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (s *customerRepo) UpdateSelectField(customer *model.Customer, fields map[string]interface{}) error {
	return s.db.Model(&customer).Updates(fields).Error
}

func (s *customerRepo) GetLastId(ctx context.Context) (int64, error) {
	var lastId int64
	return lastId, s.db.WithContext(ctx).Model(&model.Customer{}).Select("MAX(id) as id").Scan(&lastId).Error
}

func (s *customerRepo) GetAvailableBanks(ctx context.Context, customerID int64) ([]*model.CustomerAvailableBank, error) {
	var res []*model.CustomerAvailableBank
	return res, s.db.Model(model.CustomerAvailableBank{}).Preload("Bank").
		Where("customer_id = ?", customerID).Find(&res).Error
}

func (s *customerRepo) GetAvailablePlatforms(ctx context.Context, customerID int64) ([]*model.CustomerPlatform, error) {
	var res []*model.CustomerPlatform
	return res, s.db.Model(model.CustomerPlatform{}).
		Preload("Customer").Preload("Platform").
		Where("customer_id = ?", customerID).Scan(&res).Error
}

func (s *customerRepo) GetByApiKey(ctx context.Context, apiKey string) (*model.Customer, error) {
	var res *model.Customer
	return res, s.db.WithContext(ctx).Joins("JOIN integration_partners ip ON customers.id = ip.user_id").
		Where("ip.api_key = ?", apiKey).First(&res).Error
}

func (s *customerRepo) Paginate(fields *dto.ListCustomerFilter) (*dto.PaginationResponse[*model.Customer], error) {
	var items []*model.Customer
	var total int64
	tx := s.db.Model(&model.Customer{})
	if fields.Keyword != nil {
		tx.Where("CONCAT(email, full_name, referral_code) = ?", fields.Keyword)
	}
	if len(fields.IDs) > 0 {
		tx.Where("id IN (?)", fields.IDs)
	}
	if fields.Tier != nil {
		tx.Where("tier = ?", fields.Tier)
	}

	if fields.Type != nil {
		tx.Where("account_type = ?", fields.Type)
	}

	err := tx.Count(&total).Error
	if err != nil {
		return nil, err
	}

	limit := 1000
	offset := 0
	if fields.Limit != nil {
		limit = int(*fields.Limit)
	}
	if fields.Page != nil {
		offset = int(*fields.Page-1) * limit
	}
	err = tx.Limit(limit).Offset(offset).Find(&items).Error
	if err != nil {
		return nil, err
	}
	return &dto.PaginationResponse[*model.Customer]{
		Items: items,
		Total: total,
		Limit: int64(limit),
		Page:  int64(offset),
	}, nil
}
func (s *customerRepo) FindCustomerByEmail(ctx context.Context, email string) (*model.Customer, error) {
	var c model.Customer
	err := s.db.
		WithContext(ctx).
		Model(&c).
		Where("email = ?", email).
		First(&c).
		Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *customerRepo) GetTotalWalletsBalance(ctx context.Context, customerIDs []int64, walletType model.WalletType) (float64, error) {
	qb := s.db.WithContext(ctx).Table(model.TableNameWallet).Where("type = ?", walletType).Select("COALESCE(SUM(balance), 0)")
	if customerIDs != nil && len(customerIDs) > 0 {
		qb.Where("customer_id in (?)", customerIDs)
	}
	var totalBalance float64
	return totalBalance, qb.Scan(&totalBalance).Error
}
