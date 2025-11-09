package repo

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
		model "core-ledger/model/wealify"
)

type EmployeeRepo interface {
	creator[*model.Employee]
	GetByIds(ids []int64) ([]*model.Employee, error)
	GetOneByFields(ctx context.Context, fields map[string]interface{}) (*model.Employee, error)
	Count(fields map[string]interface{}) (int64, error)
	Update(*model.Employee) error
	GetMe(id int64) (*model.Employee, error)
}

type employeeRepo struct {
	db *gorm.DB
}

func NewEmployeeRepo(db *gorm.DB) EmployeeRepo {
	return &employeeRepo{db: db}
}

func (s *employeeRepo) Create(employees ...*model.Employee) error {
	return s.db.Create(employees).Error
}

func (s *employeeRepo) GetByIds(ids []int64) ([]*model.Employee, error) {
	var employees []*model.Employee
	return employees, s.db.Model(&model.Employee{}).Where("id in (?)", ids).Find(&employees).Error
}

func (s *employeeRepo) GetOneByFields(ctx context.Context, fields map[string]interface{}) (*model.Employee, error) {
	var employee *model.Employee
	return employee, s.db.WithContext(ctx).Model(&model.Employee{}).Where(fields).First(&employee).Error
}

func (s *employeeRepo) Count(fields map[string]interface{}) (int64, error) {
	return 0, nil
}

func (s *employeeRepo) Update(*model.Employee) error {
	return nil
}

func (s *employeeRepo) GetMe(id int64) (*model.Employee, error) {
	var employee *model.Employee
	query := s.db.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Info)}).
		Model(&model.Employee{}).
		Joins("Avatar").Joins("CallingCode").Joins("Country").
		Joins("Language").
		Preload("Permissions").
		//Joins("LEFT JOIN employee_permissions ep ON ep.employee_id = employees.id").
		//Joins("LEFT JOIN permissions p ON ep.permission_id = p.id", s.db.Model(&model.Permission{})).
		Where("employees.id = ?", id).Select(
		"employees.id",
		"employees.employee_id",
		"email",
		"full_name",
		"phone_number",
		"address",
		"date_of_birth",
		"two_factor_status",
		"two_factor_verification_status",
		"two_factor_enable_for",
		"authenticator_app_data_url",
		"two_factor_method",
		"registered_at",
		"changed_pw_at",
		"last_online_at",
		"employees.status",
		"employees.file_id",
		"employees.calling_code_id",
		"employees.country_id",
		"employees.language_id",
	).First(&employee)
	return employee, query.Error
}
