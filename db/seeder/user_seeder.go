package seeder

import (
	"fmt"
	model "core-ledger/model/core-ledger"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeederUser(db *gorm.DB) error {
	guardName := "web"

	// 1. Tạo Permissions
	permissions := []struct {
		Name        string
		Description string
	}{
		{Name: "users.view", Description: "Xem danh sách users"},
		{Name: "users.create", Description: "Tạo user mới"},
		{Name: "users.edit", Description: "Chỉnh sửa user"},
		{Name: "users.delete", Description: "Xóa user"},
		{Name: "roles.view", Description: "Xem danh sách roles"},
		{Name: "roles.create", Description: "Tạo role mới"},
		{Name: "roles.edit", Description: "Chỉnh sửa role"},
		{Name: "roles.delete", Description: "Xóa role"},
		{Name: "permissions.view", Description: "Xem danh sách permissions"},
		{Name: "permissions.manage", Description: "Quản lý permissions"},
		{Name: "coa.view", Description: "Xem chart of accounts"},
		{Name: "coa.create", Description: "Tạo chart of account"},
		{Name: "coa.edit", Description: "Chỉnh sửa chart of account"},
		{Name: "coa.delete", Description: "Xóa chart of account"},
		{Name: "journals.view", Description: "Xem journals"},
		{Name: "journals.create", Description: "Tạo journal"},
		{Name: "journals.edit", Description: "Chỉnh sửa journal"},
		{Name: "journals.delete", Description: "Xóa journal"},
		{Name: "reports.view", Description: "Xem báo cáo"},
		{Name: "settings.manage", Description: "Quản lý cài đặt"},
	}

	permissionMap := make(map[string]*model.Permission)
	for _, perm := range permissions {
		var existing model.Permission
		err := db.Where("name = ? AND guard_name = ?", perm.Name, guardName).First(&existing).Error
		if err == nil {
			// Đã tồn tại
			fmt.Printf("Permission đã tồn tại: %s\n", perm.Name)
			permissionMap[perm.Name] = &existing
			continue
		}
		if err != gorm.ErrRecordNotFound {
			return fmt.Errorf("lỗi khi kiểm tra permission %s: %w", perm.Name, err)
		}

		// Tạo mới
		newPermission := model.Permission{
			Name:      perm.Name,
			GuardName: guardName,
		}
		if err := db.Create(&newPermission).Error; err != nil {
			return fmt.Errorf("lỗi khi tạo permission %s: %w", perm.Name, err)
		}
		fmt.Printf("Đã tạo permission: %s\n", perm.Name)
		permissionMap[perm.Name] = &newPermission
	}

	// 2. Tạo Roles
	roles := []struct {
		Name        string
		Description string
		Permissions []string // Danh sách permission names
	}{
		{
			Name:        "super_admin",
			Description: "Super Admin - Toàn quyền",
			Permissions: []string{
				"users.view", "users.create", "users.edit", "users.delete",
				"roles.view", "roles.create", "roles.edit", "roles.delete",
				"permissions.view", "permissions.manage",
				"coa.view", "coa.create", "coa.edit", "coa.delete",
				"journals.view", "journals.create", "journals.edit", "journals.delete",
				"reports.view", "settings.manage",
			},
		},
		{
			Name:        "admin",
			Description: "Admin - Quản trị viên",
			Permissions: []string{
				"users.view", "users.create", "users.edit",
				"coa.view", "coa.create", "coa.edit",
				"journals.view", "journals.create", "journals.edit",
				"reports.view",
			},
		},
		{
			Name:        "accountant",
			Description: "Kế toán viên",
			Permissions: []string{
				"coa.view",
				"journals.view", "journals.create", "journals.edit",
				"reports.view",
			},
		},
		{
			Name:        "viewer",
			Description: "Người xem - Chỉ xem",
			Permissions: []string{
				"coa.view",
				"journals.view",
				"reports.view",
			},
		},
	}

	roleMap := make(map[string]*model.Role)
	for _, roleData := range roles {
		var existing model.Role
		err := db.Where("name = ? AND guard_name = ?", roleData.Name, guardName).First(&existing).Error
		if err == nil {
			// Đã tồn tại
			fmt.Printf("Role đã tồn tại: %s\n", roleData.Name)
			roleMap[roleData.Name] = &existing
		} else if err == gorm.ErrRecordNotFound {
			// Tạo mới
			newRole := model.Role{
				Name:      roleData.Name,
				GuardName: guardName,
			}
			if err := db.Create(&newRole).Error; err != nil {
				return fmt.Errorf("lỗi khi tạo role %s: %w", roleData.Name, err)
			}
			fmt.Printf("Đã tạo role: %s\n", roleData.Name)
			roleMap[roleData.Name] = &newRole
		} else {
			return fmt.Errorf("lỗi khi kiểm tra role %s: %w", roleData.Name, err)
		}

		// Gán permissions cho role
		role := roleMap[roleData.Name]
		for _, permName := range roleData.Permissions {
			perm, ok := permissionMap[permName]
			if !ok {
				fmt.Printf("⚠️  Warning: Permission %s không tồn tại, bỏ qua\n", permName)
				continue
			}

			// Kiểm tra xem đã gán chưa
			var existingRolePerm model.RoleHasPermission
			err := db.Where("role_id = ? AND permission_id = ?", role.ID, perm.ID).
				First(&existingRolePerm).Error
			if err == nil {
				// Đã gán rồi
				continue
			}
			if err != gorm.ErrRecordNotFound {
				return fmt.Errorf("lỗi khi kiểm tra role_has_permission: %w", err)
			}

			// Gán permission
			roleHasPermission := model.RoleHasPermission{
				RoleID:       role.ID,
				PermissionID: perm.ID,
			}
			if err := db.Create(&roleHasPermission).Error; err != nil {
				return fmt.Errorf("lỗi khi gán permission %s cho role %s: %w", permName, roleData.Name, err)
			}
			fmt.Printf("  ✓ Đã gán permission '%s' cho role '%s'\n", permName, roleData.Name)
		}
	}

	// 3. Tạo Users
	users := []struct {
		Email    string
		Password string
		FullName string
		RoleName string
	}{
		{
			Email:    "superadmin@example.com",
			Password: "SuperAdmin123!",
			FullName: "Super Administrator",
			RoleName: "super_admin",
		},
		{
			Email:    "admin@example.com",
			Password: "Admin123!",
			FullName: "Administrator",
			RoleName: "admin",
		},
		{
			Email:    "accountant@example.com",
			Password: "Accountant123!",
			FullName: "Kế toán viên",
			RoleName: "accountant",
		},
		{
			Email:    "viewer@example.com",
			Password: "Viewer123!",
			FullName: "Người xem",
			RoleName: "viewer",
		},
	}

	for _, userData := range users {
		// Kiểm tra user đã tồn tại chưa (chỉ query ID để tránh load relationships)
		var existing model.User
		err := db.Select("id", "email").Where("email = ?", userData.Email).First(&existing).Error
		
		var userID uint64
		
		if err == nil {
			// User đã tồn tại
			fmt.Printf("User đã tồn tại: %s (ID: %d)\n", userData.Email, existing.ID)
			userID = existing.ID
		} else if err == gorm.ErrRecordNotFound {
			// User chưa tồn tại, tạo mới
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
			if err != nil {
				return fmt.Errorf("lỗi khi hash password cho user %s: %w", userData.Email, err)
			}

			newUser := model.User{
				Email:     userData.Email,
				Password:  string(hashedPassword),
				FullName:  userData.FullName,
				GuardName: guardName,
			}
			if err := db.Create(&newUser).Error; err != nil {
				return fmt.Errorf("lỗi khi tạo user %s: %w", userData.Email, err)
			}
			fmt.Printf("Đã tạo user: %s (%s) - ID: %d\n", userData.Email, userData.FullName, newUser.ID)
			userID = newUser.ID
		} else {
			return fmt.Errorf("lỗi khi kiểm tra user %s: %w", userData.Email, err)
		}

		// Load user object để sử dụng method AssignRole
		var user model.User
		if err := db.Select("id", "email", "guard_name").Where("id = ?", userID).First(&user).Error; err != nil {
			return fmt.Errorf("lỗi khi load user %s (ID: %d): %w", userData.Email, userID, err)
		}

		// Gán role cho user sử dụng method AssignRole
		if err := user.AssignRole(db, userData.RoleName, guardName); err != nil {
			return fmt.Errorf("lỗi khi gán role %s cho user %s: %w", userData.RoleName, userData.Email, err)
		}
		fmt.Printf("  ✓ Đã gán role '%s' cho user '%s' (user_id: %d)\n", userData.RoleName, userData.Email, userID)
		
		// Verify sau khi gán role
		var verifyRole model.ModelHasRole
		role, ok := roleMap[userData.RoleName]
		if ok {
			if err := db.Where("role_id = ? AND model_id = ? AND model_type = ?", role.ID, userID, "User").
				First(&verifyRole).Error; err != nil {
				fmt.Printf("  ⚠️  Warning: Không thể verify role sau khi gán: %v\n", err)
			} else {
				fmt.Printf("  ✓ Verified: model_has_role đã được tạo (role_id: %d, model_id: %d, model_type: %s)\n", 
					verifyRole.RoleID, verifyRole.ModelID, verifyRole.ModelType)
			}
		}
	}

	// 4. Kiểm tra và hiển thị kết quả
	fmt.Println("\n==================================================")
	fmt.Println("📊 KIỂM TRA KẾT QUẢ SEEDER:")
	fmt.Println("==================================================")
	
	// Đếm số lượng records trong các bảng
	var roleCount int64
	db.Model(&model.Role{}).Count(&roleCount)
	fmt.Printf("✓ Roles: %d records\n", roleCount)
	
	var permissionCount int64
	db.Model(&model.Permission{}).Count(&permissionCount)
	fmt.Printf("✓ Permissions: %d records\n", permissionCount)
	
	var roleHasPermCount int64
	db.Model(&model.RoleHasPermission{}).Count(&roleHasPermCount)
	fmt.Printf("✓ Role Has Permissions: %d records\n", roleHasPermCount)
	
	var userCount int64
	db.Model(&model.User{}).Count(&userCount)
	fmt.Printf("✓ Users: %d records\n", userCount)
	
	var modelHasRoleCount int64
	db.Model(&model.ModelHasRole{}).Count(&modelHasRoleCount)
	fmt.Printf("✓ Model Has Roles: %d records\n", modelHasRoleCount)
	
	var modelHasPermCount int64
	db.Model(&model.ModelHasPermission{}).Count(&modelHasPermCount)
	fmt.Printf("✓ Model Has Permissions: %d records\n", modelHasPermCount)
	
	// Hiển thị chi tiết user và role của họ
	fmt.Println("\n📋 CHI TIẾT USER VÀ ROLE:")
	var userRoles []struct {
		UserEmail string
		UserID    uint64
		RoleName  string
		RoleID    uint64
	}
	db.Table("users").
		Select("users.email as user_email, users.id as user_id, roles.name as role_name, roles.id as role_id").
		Joins("JOIN model_has_roles ON users.id = model_has_roles.model_id AND model_has_roles.model_type = 'User'").
		Joins("JOIN roles ON model_has_roles.role_id = roles.id").
		Scan(&userRoles)
	
	if len(userRoles) == 0 {
		fmt.Println("  ⚠️  Không tìm thấy user nào có role!")
	} else {
		for _, ur := range userRoles {
			fmt.Printf("  - User: %s (ID: %d) -> Role: %s (ID: %d)\n", ur.UserEmail, ur.UserID, ur.RoleName, ur.RoleID)
		}
	}

	// Hiển thị chi tiết records trong model_has_roles
	fmt.Println("\n📋 CHI TIẾT MODEL_HAS_ROLES:")
	var allModelHasRoles []model.ModelHasRole
	if err := db.Find(&allModelHasRoles).Error; err != nil {
		fmt.Printf("  ⚠️  Lỗi khi query model_has_roles: %v\n", err)
	} else {
		if len(allModelHasRoles) == 0 {
			fmt.Println("  ⚠️  Bảng model_has_roles trống!")
		} else {
			fmt.Printf("  Tổng số records: %d\n", len(allModelHasRoles))
			for i, mhr := range allModelHasRoles {
				if i < 10 { // Chỉ hiển thị 10 records đầu
					fmt.Printf("    [%d] role_id: %d, model_id: %d, model_type: %s\n", 
						i+1, mhr.RoleID, mhr.ModelID, mhr.ModelType)
				}
			}
			if len(allModelHasRoles) > 10 {
				fmt.Printf("    ... và %d records khác\n", len(allModelHasRoles)-10)
			}
		}
	}

	// Hiển thị chi tiết records trong model_has_permissions
	fmt.Println("\n📋 CHI TIẾT MODEL_HAS_PERMISSIONS:")
	var allModelHasPerms []model.ModelHasPermission
	if err := db.Find(&allModelHasPerms).Error; err != nil {
		fmt.Printf("  ⚠️  Lỗi khi query model_has_permissions: %v\n", err)
	} else {
		if len(allModelHasPerms) == 0 {
			fmt.Println("  ℹ️  Bảng model_has_permissions trống (user có permission thông qua role)")
		} else {
			fmt.Printf("  Tổng số records: %d\n", len(allModelHasPerms))
			for i, mhp := range allModelHasPerms {
				if i < 10 { // Chỉ hiển thị 10 records đầu
					fmt.Printf("    [%d] permission_id: %d, model_id: %d, model_type: %s\n", 
						i+1, mhp.PermissionID, mhp.ModelID, mhp.ModelType)
				}
			}
			if len(allModelHasPerms) > 10 {
				fmt.Printf("    ... và %d records khác\n", len(allModelHasPerms)-10)
			}
		}
	}
	
	fmt.Println("\n✅ Seeder User, Role, Permission hoàn thành!")
	return nil
}
