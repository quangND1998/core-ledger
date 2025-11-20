package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Entity
	ID        uint64    `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Email     string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"`
	FullName  string    `gorm:"type:varchar(255)" json:"full_name"`
	GuardName string    `gorm:"type:varchar(50);default:'web'" json:"guard_name"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	// Relations
	// Note: Polymorphic many2many không được GORM hỗ trợ trực tiếp
	// Sử dụng các method AssignRole, GivePermissionTo để gán role/permission
	// và LoadRoles, LoadPermissions để load dữ liệu
	Roles       []Role       `gorm:"-" json:"roles,omitempty"`       // Load bằng method LoadRoles()
	Permissions []Permission `gorm:"-" json:"permissions,omitempty"` // Load bằng method LoadPermissions()
}

func (u *User) TableName() string {
	return "users"
}

// HasPermission checks if user has a specific permission (directly or via role)
func (u *User) HasPermission(db *gorm.DB, permissionName string, guardName string) (bool, error) {
	if guardName == "" {
		guardName = u.GuardName
	}

	// Check direct permission
	var directCount int64
	err := db.Model(&ModelHasPermission{}).
		Joins("JOIN permissions ON model_has_permissions.permission_id = permissions.id").
		Where("model_has_permissions.model_id = ? AND model_has_permissions.model_type = ?", u.ID, "User").
		Where("permissions.name = ? AND permissions.guard_name = ?", permissionName, guardName).
		Count(&directCount).Error
	if err != nil {
		return false, err
	}
	if directCount > 0 {
		return true, nil
	}

	// Check via role: Load roles của user, rồi check permissions của từng role
	var userRoles []Role
	err = db.Model(&ModelHasRole{}).
		Joins("JOIN roles ON model_has_roles.role_id = roles.id").
		Where("model_has_roles.model_id = ? AND model_has_roles.model_type = ?", u.ID, "User").
		Where("roles.guard_name = ?", guardName).
		Select("roles.*").
		Find(&userRoles).Error
	if err != nil {
		return false, err
	}

	// Check permissions của từng role (sử dụng relationship)
	for _, role := range userRoles {
		var rolePermissions []Permission
		err = db.Model(&role).
			Association("Permissions").
			Find(&rolePermissions, "name = ? AND guard_name = ?", permissionName, guardName)
		if err == nil && len(rolePermissions) > 0 {
			return true, nil
		}
	}

	return false, nil
}

// HasAnyPermission checks if user has any of the specified permissions
func (u *User) HasAnyPermission(db *gorm.DB, permissionNames []string, guardName string) (bool, error) {
	if guardName == "" {
		guardName = u.GuardName
	}

	// Check direct permissions
	var directCount int64
	err := db.Model(&ModelHasPermission{}).
		Joins("JOIN permissions ON model_has_permissions.permission_id = permissions.id").
		Where("model_has_permissions.model_id = ? AND model_has_permissions.model_type = ?", u.ID, "User").
		Where("permissions.name IN ? AND permissions.guard_name = ?", permissionNames, guardName).
		Count(&directCount).Error
	if err != nil {
		return false, err
	}
	if directCount > 0 {
		return true, nil
	}

	// Check via roles: Load roles và check permissions qua relationship
	var userRoles []Role
	err = db.Model(&ModelHasRole{}).
		Joins("JOIN roles ON model_has_roles.role_id = roles.id").
		Where("model_has_roles.model_id = ? AND model_has_roles.model_type = ?", u.ID, "User").
		Where("roles.guard_name = ?", guardName).
		Select("roles.*").
		Find(&userRoles).Error
	if err != nil {
		return false, err
	}

	// Check permissions của từng role
	for _, role := range userRoles {
		var rolePermissions []Permission
		err = db.Model(&role).
			Association("Permissions").
			Find(&rolePermissions, "name IN ? AND guard_name = ?", permissionNames, guardName)
		if err == nil && len(rolePermissions) > 0 {
			return true, nil
		}
	}

	return false, nil
}

// HasAllPermissions checks if user has all of the specified permissions
func (u *User) HasAllPermissions(db *gorm.DB, permissionNames []string, guardName string) (bool, error) {
	if guardName == "" {
		guardName = u.GuardName
	}

	// Get all permissions user has (direct + via roles)
	allPermissions, err := u.GetAllPermissions(db, guardName)
	if err != nil {
		return false, err
	}

	// Create a map for quick lookup
	permMap := make(map[string]bool)
	for _, perm := range allPermissions {
		permMap[perm.Name] = true
	}

	// Check if user has all required permissions
	for _, permName := range permissionNames {
		if !permMap[permName] {
			return false, nil
		}
	}

	return true, nil
}

// GetAllPermissions returns all permissions user has (direct + via roles)
func (u *User) GetAllPermissions(db *gorm.DB, guardName string) ([]Permission, error) {
	if guardName == "" {
		guardName = u.GuardName
	}

	var permissions []Permission

	// Get direct permissions
	var directPerms []Permission
	err := db.Model(&ModelHasPermission{}).
		Joins("JOIN permissions ON model_has_permissions.permission_id = permissions.id").
		Where("model_has_permissions.model_id = ? AND model_has_permissions.model_type = ?", u.ID, "User").
		Where("permissions.guard_name = ?", guardName).
		Select("permissions.*").
		Find(&directPerms).Error
	if err != nil {
		return nil, err
	}

	// Get permissions via roles: Load roles trước, rồi load permissions của từng role
	var userRoles []Role
	err = db.Model(&ModelHasRole{}).
		Joins("JOIN roles ON model_has_roles.role_id = roles.id").
		Where("model_has_roles.model_id = ? AND model_has_roles.model_type = ?", u.ID, "User").
		Where("roles.guard_name = ?", guardName).
		Select("roles.*").
		Find(&userRoles).Error
	if err != nil {
		return nil, err
	}

	// Load permissions của từng role qua relationship
	permMap := make(map[uint64]Permission)
	for _, perm := range directPerms {
		permMap[perm.ID] = perm
	}

	for _, role := range userRoles {
		var rolePermissions []Permission
		err = db.Model(&role).
			Association("Permissions").
			Find(&rolePermissions, "guard_name = ?", guardName)
		if err == nil {
			for _, perm := range rolePermissions {
				permMap[perm.ID] = perm
			}
		}
	}

	permissions = make([]Permission, 0, len(permMap))
	for _, perm := range permMap {
		permissions = append(permissions, perm)
	}

	return permissions, nil
}

// GetRoleNames returns all role names user has
func (u *User) GetRoleNames(db *gorm.DB, guardName string) ([]string, error) {
	if guardName == "" {
		guardName = u.GuardName
	}

	var roleNames []string
	err := db.Model(&ModelHasRole{}).
		Joins("JOIN roles ON model_has_roles.role_id = roles.id").
		Where("model_has_roles.model_id = ? AND model_has_roles.model_type = ?", u.ID, "User").
		Where("roles.guard_name = ?", guardName).
		Pluck("roles.name", &roleNames).Error
	if err != nil {
		return nil, err
	}

	return roleNames, nil
}

// GetPermissionNames returns direct permission names user has (not via roles)
func (u *User) GetPermissionNames(db *gorm.DB, guardName string) ([]string, error) {
	if guardName == "" {
		guardName = u.GuardName
	}

	var permissionNames []string
	err := db.Model(&ModelHasPermission{}).
		Joins("JOIN permissions ON model_has_permissions.permission_id = permissions.id").
		Where("model_has_permissions.model_id = ? AND model_has_permissions.model_type = ?", u.ID, "User").
		Where("permissions.guard_name = ?", guardName).
		Pluck("permissions.name", &permissionNames).Error
	if err != nil {
		return nil, err
	}

	return permissionNames, nil
}

// HasRole checks if user has a specific role
func (u *User) HasRole(db *gorm.DB, roleName string, guardName string) (bool, error) {
	if guardName == "" {
		guardName = u.GuardName
	}

	var count int64
	err := db.Model(&ModelHasRole{}).
		Joins("JOIN roles ON model_has_roles.role_id = roles.id").
		Where("model_has_roles.model_id = ? AND model_has_roles.model_type = ?", u.ID, "User").
		Where("roles.name = ? AND roles.guard_name = ?", roleName, guardName).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// HasAnyRole checks if user has any of the specified roles
func (u *User) HasAnyRole(db *gorm.DB, roleNames []string, guardName string) (bool, error) {
	if guardName == "" {
		guardName = u.GuardName
	}

	var count int64
	err := db.Model(&ModelHasRole{}).
		Joins("JOIN roles ON model_has_roles.role_id = roles.id").
		Where("model_has_roles.model_id = ? AND model_has_roles.model_type = ?", u.ID, "User").
		Where("roles.name IN ? AND roles.guard_name = ?", roleNames, guardName).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// HasAllRoles checks if user has all of the specified roles
func (u *User) HasAllRoles(db *gorm.DB, roleNames []string, guardName string) (bool, error) {
	if guardName == "" {
		guardName = u.GuardName
	}

	var count int64
	err := db.Model(&ModelHasRole{}).
		Joins("JOIN roles ON model_has_roles.role_id = roles.id").
		Where("model_has_roles.model_id = ? AND model_has_roles.model_type = ?", u.ID, "User").
		Where("roles.name IN ? AND roles.guard_name = ?", roleNames, guardName).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return int64(len(roleNames)) == count, nil
}

// GivePermissionTo assigns a permission directly to user
func (u *User) GivePermissionTo(db *gorm.DB, permissionName string, guardName string) error {
	if guardName == "" {
		guardName = u.GuardName
	}

	// Find or create permission
	var permission Permission
	err := db.Where("name = ? AND guard_name = ?", permissionName, guardName).FirstOrCreate(&permission, Permission{
		Name:      permissionName,
		GuardName: guardName,
	}).Error
	if err != nil {
		return err
	}

	// Check if already assigned
	var existing ModelHasPermission
	err = db.Where("permission_id = ? AND model_id = ? AND model_type = ?", permission.ID, u.ID, "User").
		First(&existing).Error
	if err == nil {
		// Already assigned
		return nil
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	// Assign permission
	modelHasPermission := ModelHasPermission{
		PermissionID: permission.ID,
		ModelID:      u.ID,
		ModelType:    "User",
	}
	return db.Create(&modelHasPermission).Error
}

// RevokePermissionTo removes a permission from user
func (u *User) RevokePermissionTo(db *gorm.DB, permissionName string, guardName string) error {
	if guardName == "" {
		guardName = u.GuardName
	}

	var permission Permission
	err := db.Where("name = ? AND guard_name = ?", permissionName, guardName).First(&permission).Error
	if err != nil {
		return err
	}

	return db.Where("permission_id = ? AND model_id = ? AND model_type = ?", permission.ID, u.ID, "User").
		Delete(&ModelHasPermission{}).Error
}

// AssignRole assigns a role to user
func (u *User) AssignRole(db *gorm.DB, roleName string, guardName string) error {
	if guardName == "" {
		guardName = u.GuardName
	}

	// Find or create role
	var role Role
	err := db.Where("name = ? AND guard_name = ?", roleName, guardName).FirstOrCreate(&role, Role{
		Name:      roleName,
		GuardName: guardName,
	}).Error
	if err != nil {
		return err
	}

	// Check if already assigned
	var existing ModelHasRole
	err = db.Where("role_id = ? AND model_id = ? AND model_type = ?", role.ID, u.ID, "User").
		First(&existing).Error
	if err == nil {
		// Already assigned
		return nil
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	// Assign role
	modelHasRole := ModelHasRole{
		RoleID:    role.ID,
		ModelID:   u.ID,
		ModelType: "User",
	}
	return db.Create(&modelHasRole).Error
}

// RemoveRole removes a role from user
func (u *User) RemoveRole(db *gorm.DB, roleName string, guardName string) error {
	if guardName == "" {
		guardName = u.GuardName
	}

	var role Role
	err := db.Where("name = ? AND guard_name = ?", roleName, guardName).First(&role).Error
	if err != nil {
		return err
	}

	return db.Where("role_id = ? AND model_id = ? AND model_type = ?", role.ID, u.ID, "User").
		Delete(&ModelHasRole{}).Error
}

// SyncPermissions syncs user permissions (removes all and assigns new ones)
func (u *User) SyncPermissions(db *gorm.DB, permissionNames []string, guardName string) error {
	if guardName == "" {
		guardName = u.GuardName
	}

	// Remove all existing permissions
	err := db.Where("model_id = ? AND model_type = ?", u.ID, "User").
		Delete(&ModelHasPermission{}).Error
	if err != nil {
		return err
	}

	// Assign new permissions
	for _, permName := range permissionNames {
		if err := u.GivePermissionTo(db, permName, guardName); err != nil {
			return err
		}
	}

	return nil
}

// SyncRoles syncs user roles (removes all and assigns new ones)
func (u *User) SyncRoles(db *gorm.DB, roleNames []string, guardName string) error {
	if guardName == "" {
		guardName = u.GuardName
	}

	// Remove all existing roles
	err := db.Where("model_id = ? AND model_type = ?", u.ID, "User").
		Delete(&ModelHasRole{}).Error
	if err != nil {
		return err
	}

	// Assign new roles
	for _, roleName := range roleNames {
		if err := u.AssignRole(db, roleName, guardName); err != nil {
			return err
		}
	}

	return nil
}

// SyncRolesByIDs syncs user roles by role IDs (removes all and assigns new ones)
func (u *User) SyncRolesByIDs(db *gorm.DB, roleIDs []uint64) error {
	// Remove all existing roles
	err := db.Where("model_id = ? AND model_type = ?", u.ID, "User").
		Delete(&ModelHasRole{}).Error
	if err != nil {
		return err
	}

	// Assign new roles by IDs
	if len(roleIDs) > 0 {
		modelHasRoles := make([]ModelHasRole, 0, len(roleIDs))
		for _, roleID := range roleIDs {
			modelHasRoles = append(modelHasRoles, ModelHasRole{
				RoleID:    roleID,
				ModelID:   u.ID,
				ModelType: "User",
			})
		}
		return db.Create(&modelHasRoles).Error
	}

	return nil
}

// LoadRoles loads roles for user into u.Roles field
func (u *User) LoadRoles(db *gorm.DB, guardName string) error {
	if guardName == "" {
		guardName = u.GuardName
	}

	var roles []Role
	err := db.Model(&ModelHasRole{}).
		Joins("JOIN roles ON model_has_roles.role_id = roles.id").
		Where("model_has_roles.model_id = ? AND model_has_roles.model_type = ?", u.ID, "User").
		Where("roles.guard_name = ?", guardName).
		Select("roles.*").
		Find(&roles).Error
	if err != nil {
		return err
	}

	u.Roles = roles
	return nil
}

// LoadPermissions loads direct permissions for user into u.Permissions field
func (u *User) LoadPermissions(db *gorm.DB, guardName string) error {
	if guardName == "" {
		guardName = u.GuardName
	}

	var permissions []Permission
	err := db.Model(&ModelHasPermission{}).
		Joins("JOIN permissions ON model_has_permissions.permission_id = permissions.id").
		Where("model_has_permissions.model_id = ? AND model_has_permissions.model_type = ?", u.ID, "User").
		Where("permissions.guard_name = ?", guardName).
		Select("permissions.*").
		Find(&permissions).Error
	if err != nil {
		return err
	}

	u.Permissions = permissions
	return nil
}
