package model

const OldTableNameEmployeePermission = "employee-permission"

type OldEmployeePermission struct {
	EmployeeID   int32 `gorm:"column:employee_id;primaryKey" json:"employee_id"`
	PermissionID int32 `gorm:"column:permission_id;primaryKey" json:"permission_id"`
}

// TableName EmployeePermission's table name
func (*OldEmployeePermission) TableName() string {
	return OldTableNameEmployeePermission
}
