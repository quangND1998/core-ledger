package repo

import (
	"context"

	"gorm.io/gorm/logger"

	"core-ledger/internal/module/constants"
	"core-ledger/model/dto"
	"core-ledger/model/enum"
	model "core-ledger/model/wealify"
	wv "core-ledger/pkg/utils/wrapvalue"
	"strings"
	"time"

	"gorm.io/gorm"
)

type TransactionRepo interface {
	creator[*model.Transaction]
	reader[*model.Transaction, *TransactionFilter]
	getByID[*model.Transaction]
	updater[*model.Transaction]
	GetList(ctx context.Context) ([]model.Transaction, error)
	Saves(trans []*model.Transaction) error
	GetByIds(ids []int64) ([]*model.Transaction, error)
	Count(fields map[string]interface{}) (int64, error)
	GetLastId(ctx context.Context) (int64, error)
	UpdateWaitingHPayTransactions(ctx context.Context) (int64, error)
	GetPendingAndProcessWithdrawAmount(ctx context.Context, userID int64, walletType model.WalletType) (float64, error)
	GetTotalTopUpSuccess(ctx context.Context, userID int64, walletType model.WalletType) (float64, error)
	InTx(ctx context.Context, fn func(tx *gorm.DB) error) error
	FindTxProcessFromPlatforms(ctx context.Context, platformIDs []string) ([]model.TransactionAutoProcessToApprove, error)
	UpdateTransactionNote(ctx context.Context, id string, note string) error
	FindApprovedTxByCustomerIDs(ctx context.Context, customerIDs []int64) ([]int64, error)
	UpdateByCondition(ctx context.Context, where map[string]interface{}, updates map[string]interface{}) error
	PaginateManual(fields *TransactionFilter) (*dto.PaginationResponse[*model.TransactionView], error)
}
type transactionRepo struct {
	db *gorm.DB
}

type TransactionFilter struct {
	dto.BasePaginationQuery
	GetFields             []string                   `json:"get_fields" form:"get_fields"`
	TransactionStatuses   []enum.TransactionStatus   `json:"transaction_statuses,omitempty" form:"transaction_statuses"`
	TransactionStatus     *string                    `json:"transaction_status,omitempty" form:"transaction_status"`
	VaTransactionStatuses []enum.VaTransactionStatus `json:"va_transaction_statuses,omitempty" form:"va_transaction_statuses"`
	VaTransactionStatus   *enum.VaTransactionStatus  `json:"va_transaction_status,omitempty" form:"va_transaction_status"`
	TransactionTypes      []enum.TransactionType     `json:"transaction_types,omitempty" form:"transaction_types"`
	TransactionType       *enum.TransactionType      `json:"transaction_type,omitempty" form:"transaction_type"`
	WalletTypes           []model.WalletType         `json:"wallet_types,omitempty" form:"wallet_types"`
	TopUpStatuses         []string                   `json:"top_up_statuses,omitempty" form:"top_up_statuses[]"`
	WithdrawStatuses      []string                   `json:"withdraw_statuses,omitempty" form:"withdraw_statuses[]"`
	MinAmount             *float64                   `json:"min_amount,omitempty" form:"min_amount"`
	MaxAmount             *float64                   `json:"max_amount,omitempty" form:"max_amount"`
	ThirdPartyStatus      *string                    `json:"third_party_status,omitempty" form:"third_party_status"`
	SystemPaymentIds      []string                   `json:"system_payment_ids,omitempty" form:"system_payment_ids"`
	Providers             []string                   `json:"providers,omitempty" form:"providers"`
	CustomerInfo          *string                    `json:"customer_info,omitempty" form:"customer_info"`
	OwnerId               *string                    `json:"owner_id,omitempty" form:"owner_id"`
	CustomerIds           []int64                    `json:"customer_ids,omitempty" form:"customer_ids"`
	IsConfirming          *string                    `json:"is_confirming,omitempty" form:"is_confirming"`
	TransactionIds        []string                   `json:"transaction_ids,omitempty" form:"transaction_ids"`
	InHouse               *bool                      `json:"in_house,omitempty" form:"in_house"`
	TopUpType             *string                    `json:"top_up_type,omitempty" form:"top_up_type"`

	// builder
	NeedHistory bool
}

func (s *transactionRepo) FindApprovedTxByCustomerIDs(ctx context.Context, customerIDs []int64) ([]int64, error) {
	var rows []struct {
		CustomerID int64 `json:"customer_id"`
	}

	err := s.db.WithContext(ctx).
		Model(&model.Transaction{}).
		Select("DISTINCT(customer_id)").
		Where("customer_id IN (?)", customerIDs).
		Where("transaction_status = ?", enum.TransactionStatusApproved).
		Where("transaction_type = ?", enum.TransactionTypeTopUp).
		Where("virtual_account_id is not null").
		Find(&rows).
		Error
	if err != nil {
		return nil, err
	}

	ids := make([]int64, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.CustomerID)
	}

	return ids, nil
}

func (s *transactionRepo) UpdateTransactionNote(ctx context.Context, transactionID string, note string) error {
	return s.db.WithContext(ctx).Exec("UPDATE transactions t SET auto_approve_note = ? WHERE id = ?", note, transactionID).Error
}

func (s *transactionRepo) FindTxProcessFromPlatforms(ctx context.Context, platformIDs []string) ([]model.TransactionAutoProcessToApprove, error) {
	var platforms []model.TransactionAutoProcessToApprove

	err := s.db.WithContext(ctx).
		Table("transactions t").
		Joins("INNER JOIN `virtual-accounts` va ON va.id = t.virtual_account_id and va.platform_id is not null").
		Joins("INNER JOIN `customers` c ON c.id = t.customer_id").
		Where("t.transaction_type = ?", model.TransactionVcTypeTopUp).
		Where("t.transaction_status = ?", enum.TransactionStatusProcess).
		Where("va.platform_id IN ?", platformIDs).
		Where("t.auto_approve_note is null").
		Where("t.updated_at <= ?", time.Now().AddDate(0, 0, -1)). // For auto approve, the Transaction must be in Process Status at least 1 day.
		Select("c.kyb_status, t.amount, t.id, t.transaction_status, va.platform_id, t.remark").
		Find(&platforms).
		Error
	if err != nil {
		return nil, err
	}

	return platforms, nil
}

func NewTransactionRepo(db *gorm.DB) TransactionRepo {
	return &transactionRepo{db: db}
}

//	func (r *transactionRepo) WithSchema(schema string) *transactionRepo {
//		schemaDB := r.db.Session(&gorm.Session{})
//		schemaDB.Exec("SET search_path TO " + schema)
//		return &transactionRepo{db: schemaDB}
//	}
func (s *transactionRepo) GetList(ctx context.Context) ([]model.Transaction, error) {
	var products []model.Transaction
	err := s.db.WithContext(ctx).Limit(1000).Offset(0).Find(&products).Error

	return products, err
}

func (s *transactionRepo) Create(transactions ...*model.Transaction) error {
	return s.db.Create(transactions).Error
}

func (s *transactionRepo) GetByID(ctx context.Context, id int64) (*model.Transaction, error) {
	customer := &model.Transaction{}
	return customer, s.db.WithContext(ctx).First(&customer, "id = ?", id).Error
}

func (s *transactionRepo) Saves(trans []*model.Transaction) error {
	return s.db.Model(&model.Transaction{}).Save(trans).Error
}

func (s *transactionRepo) GetByIds(ids []int64) ([]*model.Transaction, error) {
	var sessions []*model.Transaction
	return sessions, s.db.Model(&model.Transaction{}).Where("id in (?)", ids).Find(&sessions).Error
}

func (s *transactionRepo) GetOneByFields(ctx context.Context, fields map[string]interface{}, preloads ...string) (*model.Transaction, error) {
	var customer *model.Transaction
	query := s.db.WithContext(ctx).Model(&model.Transaction{})
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	return customer, query.Where(fields).First(&customer).Error
}

func (t *transactionRepo) GetManyByFields(ctx context.Context, fields map[string]interface{}, preloads ...string) ([]*model.Transaction, error) {
	var trans []*model.Transaction
	query := t.db.WithContext(ctx).Model(&model.Transaction{}).Where(fields)
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	return trans, query.Find(&trans).Error
}

func (t *transactionRepo) GetPendingAndProcessWithdrawAmount(ctx context.Context, userID int64, walletType model.WalletType) (float64, error) {
	var res float64
	return res, t.db.WithContext(ctx).Model(&model.Transaction{}).
		Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Info)}).
		Select("COALESCE(SUM(transactions.amount), 0)").
		Joins("LEFT JOIN wallets sw ON sw.id = transactions.sent_wallet_id").
		Joins("LEFT JOIN wallets rw ON rw.id = transactions.received_wallet_id").
		Where("transactions.transaction_type = ?", enum.TransactionTypeWithdrawal).
		Where("transactions.transaction_status IN (?)", []enum.TransactionStatus{enum.TransactionStatusPending, enum.TransactionStatusProcess}).
		Where("transactions.customer_id = ?", userID).
		Where("sw.type = ? OR rw.type = ?", walletType, walletType).
		Where(`(transactions.transaction_linked_id IS NULL AND
			transactions.id NOT IN (
				SELECT t2.transaction_linked_id from transactions t2
                where t2.transaction_linked_id is not null
			) OR
       		transactions.transaction_linked_id IS NOT NULL
       	)`).
		Scan(&res).Error
}

// 3856

func (t *transactionRepo) GetTotalTopUpSuccess(ctx context.Context, userID int64, walletType model.WalletType) (float64, error) {
	var res float64
	return res, t.db.WithContext(ctx).Model(&model.Transaction{}).
		Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Info)}).
		Select("COALESCE(SUM(transactions.payment_amount), 0)").
		Joins("LEFT JOIN wallets sw ON sw.id = transactions.sent_wallet_id").
		Joins("LEFT JOIN wallets rw ON rw.id = transactions.received_wallet_id").
		Where("transactions.status = true").
		Where("transactions.is_deleted = false").
		Where("transactions.transaction_type = ? or transactions.transaction_type = ? and transactions.received_wallet_id is not null", enum.TransactionTypeTopUp, enum.TransactionTypeAdjustment).
		Where("transactions.transaction_status = ?", enum.TransactionStatusApproved).
		Where("transactions.customer_id = ?", userID).
		Where("sw.type = ? OR rw.type = ?", walletType, walletType).Scan(&res).Error
}

func (s *transactionRepo) Count(fields map[string]interface{}) (int64, error) {
	return 0, nil
}

func (s *transactionRepo) UpdateSelectField(customer *model.Transaction, fields map[string]interface{}) error {
	return s.db.Model(&customer).Updates(fields).Error
}

func (s *transactionRepo) GetLastId(ctx context.Context) (int64, error) {
	var lastId int64
	return lastId, s.db.WithContext(ctx).Model(&model.Transaction{}).Select("MAX(id) as id").Scan(&lastId).Error
}

func (s *transactionRepo) Paginate(fields *TransactionFilter) (*dto.PaginationResponse[*model.Transaction], error) {
	var items []*model.Transaction
	var total int64
	query := s.db.Model(&model.Transaction{}).Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Info)}).
		Joins("Customer").
		Joins("Currency").
		Joins("VirtualAccount").
		Joins("VirtualAccount.Platform"). // TODO add filter and join, just use for export currently
		Joins("SystemPayment").           // TODO refactor just need when export file
		Joins("SystemPayment.BankAccount").
		Joins("SystemPayment.EWallet").
		Joins("ConfirmTransaction"). // note: for detail only
		Preload("SubTransactions").
		Where("transactions.is_deleted = FALSE").
		Where("transactions.status = TRUE").
		Where("transactions.transaction_linked_id IS NULL").
		Where("transactions.is_vc_transaction = FALSE").
		Order("transactions.created_at DESC")

	if fields.NeedHistory {
		query = query.Preload("Histories")
	}

	if fields.Keyword != nil && *fields.Keyword != "" {
		likeQuery := "%" + strings.ToLower(*fields.Keyword) + "%"
		query = query.Where(
			"VirtualAccount.card_number LIKE ? OR transactions.transaction_id LIKE ?", // TODO query sub transaction
			likeQuery,
			likeQuery,
		)
	}

	if fields.StringIDs != nil && len(fields.StringIDs) > 0 {
		query = query.Where("transactions.id IN (?)", fields.StringIDs)
	}

	if fields.CustomerIds != nil && len(fields.CustomerIds) > 0 {
		query = query.Where("transactions.customer_id IN (?)", fields.CustomerIds)
	}

	if fields.Providers != nil && len(fields.Providers) > 0 {
		query = query.Where("transactions.provider IN (?)", fields.Providers)
	}

	if fields.OwnerId != nil {
		query = query.Where("transactions.customer_id = ? OR transactions.sender_id = ? OR transactions.receiver_id = ?",
			fields.OwnerId,
			fields.OwnerId,
			fields.OwnerId,
		).Where("transactions.transaction_status != 'ON_HOLD'")
	}

	if fields.InHouse != nil {
		query = query.Where("Customer.in_house = ?", fields.InHouse)
	}

	if len(fields.WalletTypes) > 0 {
		query = query.
			Joins("LEFT JOIN wallets sw ON sw.id = transactions.sent_wallet_id").
			Joins("LEFT JOIN wallets rw ON rw.id = transactions.received_wallet_id").
			Where("sw.type IN (?) OR rw.type IN (?)", fields.WalletTypes, fields.WalletTypes)
	}

	if len(fields.TransactionStatuses) > 0 {
		query = query.Where("transactions.transaction_status IN (?)", fields.TransactionStatuses)
	}

	if len(fields.TransactionTypes) > 0 {
		query = query.Where("transactions.transaction_type IN (?)", fields.TransactionTypes)
	}

	// old filter
	if fields.TransactionStatus != nil {
		query = query.Where("transactions.transaction_status = ?", fields.TransactionStatus)
	}

	if fields.TransactionType != nil {
		if *fields.TransactionType == enum.TransactionTypeTopUp {
			query = query.Where("transactions.transaction_type = ? or transactions.transaction_type = ? and transactions.received_wallet_id is not null", fields.TransactionType, enum.TransactionTypeAdjustment)
		} else if *fields.TransactionType == enum.TransactionTypeWithdrawal {
			query = query.Where("transactions.transaction_type = ? or transactions.transaction_type = ? and transactions.sent_wallet_id is not null", fields.TransactionType, enum.TransactionTypeAdjustment)
		} else {
			query = query.Where("transactions.transaction_type = ?", fields.TransactionType)
		}
	}

	if fields.VaTransactionStatus != nil {
		query = query.Where("va_transaction_status = ?", fields.VaTransactionStatus)
	}

	if fields.TopUpType != nil {
		if *fields.TopUpType == "QR" {
			query = query.Where("qr_image_url IS NOT NULL")
		} else {
			query = query.Where("qr_image_url IS NULL")
		}
	}

	if len(fields.SystemPaymentIds) > 0 {
		query = query.Where("system_payment_id IN (?)", fields.SystemPaymentIds)
	}

	layout := "2006-01-02"
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	if fields.StartDate != nil {
		startTime, err := time.ParseInLocation(layout, *fields.StartDate, loc)
		if err == nil {
			query = query.Where("transactions.created_at >= ?", startTime)
		}
	}

	if fields.EndDate != nil {
		endTime, err := time.ParseInLocation(layout, *fields.EndDate, loc)
		if err == nil {
			endTime = endTime.Add(24 * time.Hour)
			query = query.Where("transactions.created_at <= ?", endTime)
		}
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, err
	}

	if len(fields.GetFields) != 0 {
		query = query.Select(fields.GetFields)
	}

	limit := 1000000
	offset := 0
	var page int64 = 1
	if fields.Limit != nil {
		limit = int(*fields.Limit)
	}
	if fields.Page != nil {
		offset = int(*fields.Page-1) * limit
		page = *fields.Page
	}

	err = query.Limit(limit).Offset(offset).Find(&items).Error
	if err != nil {
		return nil, err
	}

	totalPage := total/int64(limit) + 1
	var nextPage *int64
	var prevPage *int64
	if page < totalPage {
		nextPage = wv.ToPointer(page + 1)
	}
	if page > 1 {
		prevPage = wv.ToPointer(page - 1)
	}

	return &dto.PaginationResponse[*model.Transaction]{
		Items:     items,
		Total:     total,
		Limit:     int64(limit),
		Page:      int64(offset/limit + 1),
		TotalPage: totalPage,
		NextPage:  nextPage,
		PrevPage:  prevPage,
	}, nil
}

func (t *transactionRepo) UpdateWaitingHPayTransactions(ctx context.Context) (int64, error) {
	result := t.db.Table(wv.ToPointer(model.Transaction{}).TableName()).
		Joins("ConfirmTransaction").
		Joins("Currency").
		Joins("Customer").
		Joins("ReceivedWallet", t.db.Where(&model.Wallet{Type: model.WalletTypeVA})).
		Joins("SentWallet").
		Joins("VirtualAccount").
		Where("transactions.va_transaction_status = ?", enum.VaTransactionStatusWaiting).
		Where("transactions.provider = 'H_PAY'").
		Where("transactions.transaction_type = ?", constants.TRANSACTION_TYPE_TOP_UP).
		Where("transactions.transaction_status IN (?)", []enum.TransactionStatus{enum.TransactionStatusPending, enum.TransactionStatusProcess}).
		Where("transactions.created_at < NOW() - INTERVAL 1 DAY").
		Where("transactions.updated_at < NOW() - INTERVAL 1 DAY").
		Updates(map[string]interface{}{
			"transactions.transaction_status":    enum.TransactionStatusProcess,
			"transactions.va_transaction_status": enum.VaTransactionStatusSuccess,
		})
	return result.RowsAffected, result.Error
}
func (t *transactionRepo) InTx(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return t.db.WithContext(ctx).Transaction(fn)
}

func (t *transactionRepo) UpdateByCondition(ctx context.Context, where map[string]interface{}, updates map[string]interface{}) error {
	return t.db.Model(&model.Transaction{}).Where(where).Updates(updates).Error
}

func (s *transactionRepo) PaginateManual(fields *TransactionFilter) (*dto.PaginationResponse[*model.TransactionView], error) {
	var items []*model.TransactionView
	var total int64

	query := s.db.Table(model.TableNameTransaction).
		Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Info)}).
		Joins("LEFT JOIN customers AS customer ON customer.id = transactions.customer_id").
		Joins("LEFT JOIN currencies AS currencies ON currencies.symbol = transactions.currency_symbol").
		Joins("LEFT JOIN `virtual-accounts` AS va ON va.id = transactions.virtual_account_id").
		Joins("LEFT JOIN `platforms` AS platforms ON va.platform_id = platforms.id").
		Joins("LEFT JOIN `providers` AS providers ON transactions.provider = providers.value").
		Where("transactions.is_deleted = FALSE").
		Where("transactions.status = TRUE").
		Where("transactions.transaction_linked_id IS NULL").
		Where("transactions.is_vc_transaction = FALSE").
		Where("transactions.transaction_status != 'ON_HOLD'").
		Order("transactions.created_at DESC")

	if fields.Keyword != nil && *fields.Keyword != "" {
		likeQuery := "%" + strings.ToLower(*fields.Keyword) + "%"
		query = query.Where(
			"va.card_number LIKE ? OR transactions.transaction_id LIKE ?", // TODO query sub transaction
			likeQuery,
			likeQuery,
		)
	}

	if len(fields.CustomerIds) > 0 {
		query = query.Where("transactions.customer_id IN (?)", fields.CustomerIds)
	}

	if len(fields.Providers) > 0 {
		query = query.Where("transactions.provider IN (?)", fields.Providers)
	}

	if fields.OwnerId != nil {
		query = query.Where("transactions.customer_id = ? OR transactions.sender_id = ? OR transactions.receiver_id = ?",
			fields.OwnerId,
			fields.OwnerId,
			fields.OwnerId,
		).Where("transactions.transaction_status != 'ON_HOLD'")
	}

	if fields.InHouse != nil {
		query = query.Where("Customer.in_house = ?", fields.InHouse)
	}

	if len(fields.WalletTypes) > 0 {
		query = query.
			Joins("LEFT JOIN wallets AS sw ON sw.id = transactions.sent_wallet_id").
			Joins("LEFT JOIN wallets AS rw ON rw.id = transactions.received_wallet_id").
			Where("sw.type IN (?) OR rw.type IN (?)", fields.WalletTypes, fields.WalletTypes)
	}

	if len(fields.TransactionStatuses) > 0 {
		query = query.Where("transaction_status IN (?)", fields.TransactionStatuses)
	}

	if len(fields.TransactionTypes) > 0 {
		query = query.Where("transaction_type IN (?)", fields.TransactionTypes)
	}

	// old filter
	if fields.TransactionStatus != nil {
		query = query.Where("transaction_status = ?", fields.TransactionStatus)
	}

	if fields.TransactionType != nil {
		query = query.Where("transaction_type = ?", fields.TransactionType)
	}

	if fields.VaTransactionStatus != nil {
		query = query.Where("va_transaction_status = ?", fields.VaTransactionStatus)
	}

	if fields.TransactionType == nil && (len(fields.TopUpStatuses) > 0 || len(fields.WithdrawStatuses) > 0) {
		var conditions []string
		var args []interface{}

		if len(fields.TopUpStatuses) > 0 {
			conditions = append(conditions, "(transactions.transaction_type = 'TOP_UP' AND transactions.transaction_status IN (?))")
			args = append(args, fields.TopUpStatuses)
		}

		if len(fields.WithdrawStatuses) > 0 {
			conditions = append(conditions, "(transactions.transaction_type = 'WITHDRAWAL' AND transactions.transaction_status IN (?))")
			args = append(args, fields.WithdrawStatuses)
		}

		if len(conditions) > 0 {
			whereClause := strings.Join(conditions, " OR ")
			query = query.Where(whereClause, args...)
		}
	}

	if fields.TopUpType != nil {
		if *fields.TopUpType == "QR" {
			query = query.Where("qr_image_url IS NOT NULL")
		} else {
			query = query.Where("qr_image_url IS NULL")
		}
	}

	if len(fields.SystemPaymentIds) > 0 {
		query = query.Where("system_payment_id IN (?)", fields.SystemPaymentIds)
	}

	layout := "2006-01-02"
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	if fields.StartDate != nil {
		startTime, err := time.ParseInLocation(layout, *fields.StartDate, loc)
		if err == nil {
			query = query.Where("transactions.created_at >= ?", startTime)
		}
	}

	if fields.EndDate != nil {
		endTime, err := time.ParseInLocation(layout, *fields.EndDate, loc)
		if err == nil {
			endTime = endTime.Add(24 * time.Hour)
			query = query.Where("transactions.created_at <= ?", endTime)
		}
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, err
	}

	if len(fields.GetFields) != 0 {
		query = query.Select(fields.GetFields)
	}

	limit := 1000
	offset := 0
	var page int64 = 1
	if fields.Limit != nil {
		limit = int(*fields.Limit)
	}
	if fields.Page != nil {
		offset = int(*fields.Page-1) * limit
		page = *fields.Page
	}

	err = query.Limit(limit).Offset(offset).Find(&items).Error
	if err != nil {
		return nil, err
	}

	totalPage := total/int64(limit) + 1
	var nextPage *int64
	var prevPage *int64
	if page < totalPage {
		nextPage = wv.ToPointer(page + 1)
	}
	if page > 1 {
		prevPage = wv.ToPointer(page - 1)
	}

	return &dto.PaginationResponse[*model.TransactionView]{
		Items:     items,
		Total:     total,
		Limit:     int64(limit),
		Page:      int64(offset/limit + 1),
		TotalPage: totalPage,
		NextPage:  nextPage,
		PrevPage:  prevPage,
	}, nil
}
