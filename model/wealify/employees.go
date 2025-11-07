package model

import (
	"time"
)

const TableNameEmployee = "employees"

// Employee mapped from table <employees>
type Employee struct {
	CreatedAt                   time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP(6)" json:"created_at,omitempty"`
	UpdatedAt                   time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP(6)" json:"updated_at,omitempty"`
	Status                      bool      `gorm:"column:status;not null;default:1" json:"status,omitempty"`
	IsDeleted                   bool      `gorm:"column:is_deleted;not null" json:"is_deleted,omitempty"`
	ID                          int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id,omitempty"`
	Code                        string    `gorm:"column:employee_id;not null" json:"employee_id,omitempty"`
	FullName                    string    `gorm:"column:full_name;not null" json:"full_name,omitempty"`
	Email                       string    `gorm:"column:email;not null" json:"email,omitempty"`
	PhoneNumber                 string    `gorm:"column:phone_number" json:"phone_number,omitempty"`
	RefBy                       string    `gorm:"column:ref_by" json:"ref_by,omitempty"`
	RefCount                    int32     `gorm:"column:ref_count;not null" json:"ref_count,omitempty"`
	Address                     string    `gorm:"column:address" json:"address,omitempty"`
	DateOfBirth                 time.Time `gorm:"column:date_of_birth" json:"date_of_birth,omitempty"`
	Password                    string    `gorm:"column:password;not null" json:"password,omitempty"`
	TwoFactorStatus             string    `gorm:"column:two_factor_status;not null;default:DISABLE" json:"two_factor_status,omitempty"`
	TwoFactorVerificationStatus string    `gorm:"column:two_factor_verification_status;not null;default:UNVERIFIED" json:"two_factor_verification_status,omitempty"`
	TwoFactorMethod             string    `gorm:"column:two_factor_method;not null;default:EMAIL" json:"two_factor_method,omitempty"`
	AuthenticatorAppSecretKey   string    `gorm:"column:authenticator_app_secret_key" json:"authenticator_app_secret_key,omitempty"`
	AuthenticatorAppDataURL     string    `gorm:"column:authenticator_app_data_url" json:"authenticator_app_data_url,omitempty"`
	RegisteredAt                time.Time `gorm:"column:registered_at;not null" json:"registered_at,omitempty"`
	ChangedPwAt                 time.Time `gorm:"column:changed_pw_at;not null" json:"changed_pw_at,omitempty"`
	LastOnlineAt                time.Time `gorm:"column:last_online_at;not null" json:"last_online_at,omitempty"`

	FileID        string `gorm:"column:file_id" json:"file_id,omitempty"`
	CountryID     string `gorm:"column:country_id" json:"country_id,omitempty"`
	CallingCodeID string `gorm:"column:calling_code_id" json:"calling_code_id,omitempty"`
	LanguageID    string `gorm:"column:language_id" json:"language_id,omitempty"`

	Avatar              *File                 `gorm:"foreignKey:FileID;references:ID" json:"avatar,omitempty"`
	CallingCode         *CallingCode          `gorm:"foreignKey:CallingCodeID;references:ID" json:"calling_code,omitempty"`
	Country             *Countrie             `gorm:"foreignKey:CountryID;references:ID" json:"country,omitempty"`
	Language            *Language             `gorm:"foreignKey:LanguageID;references:ID" json:"language,omitempty"`
	EmployeePermissions []*EmployeePermission `gorm:"foreignKey:EmployeeID;references:ID" json:"employee_permissions,omitempty"`
	Permissions         []*Permission         `gorm:"many2many:employee_permissions;joinForeignKey:EmployeeID;joinReferences:PermissionID" json:"permissions,omitempty"`
}

// TableName Employee's table name
func (*Employee) TableName() string {
	return TableNameEmployee
}
