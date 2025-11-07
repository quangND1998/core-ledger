package model

const TableNameEmployeePermission = "employee_permissions"

type EmployeePermission struct {
	EmployeeID   int64       `gorm:"column:employee_id;primaryKey" json:"employee_id"`
	PermissionID int64       `gorm:"column:permission_id;primaryKey" json:"permission_id"`
	Employee     *Employee   `gorm:"foreignKey:EmployeeID;references:ID" json:"employee"`
	Permission   *Permission `gorm:"foreignKey:PermissionID;references:ID" json:"permission"`
}

// TableName EmployeePermission's table name
func (*EmployeePermission) TableName() string {
	return TableNameEmployeePermission
}
