package user

import (
	"core-ledger/pkg/repo"
	"github.com/gin-gonic/gin"
)

// Example: Cách sử dụng permission middleware với routes

// Example 1: Apply permission middleware cho toàn bộ group
func ExampleSetupRoutesWithPermission(rg *gin.RouterGroup, h *UserHandler, userRepo repo.UserRepo) {
	// Tạo middleware
	authMiddleware := UserAuthMiddleware(userRepo)
	permissionMiddleware := UserPermissionMiddleware("manage-users", "web")
	
	// Apply cả auth và permission middleware cho toàn bộ routes
	SetupRoutes(rg, h, authMiddleware, permissionMiddleware)
}

// Example 2: Apply permission middleware cho từng route cụ thể
func ExampleSetupRoutesWithSpecificPermissions(rg *gin.RouterGroup, h *UserHandler, userRepo repo.UserRepo) {
	authMiddleware := UserAuthMiddleware(userRepo)
	
	users := rg.Group("users", authMiddleware)
	{
		// Public routes (chỉ cần auth)
		users.GET("/list", h.ListUsers)
		users.GET("/:id", h.GetUser)
		
		// Routes cần permission cụ thể
		users.POST("", 
			h.CreateUser, 
			UserPermissionMiddleware("create-users", "web"))
		
		users.PUT("/:id", 
			h.UpdateUser, 
			UserPermissionMiddleware("update-users", "web"))
		
		// Routes cần role cụ thể
		users.POST("/:id/roles", 
			h.SyncUserRoles, 
			UserRoleMiddleware("admin", "web"))
		
		// Routes cần bất kỳ role nào
		users.POST("/:id/permissions", 
			h.SyncUserPermissions, 
			UserAnyRoleMiddleware([]string{"admin", "super-admin"}, "web"))
	}
}

// Example 3: Combine multiple permissions
func ExampleSetupRoutesWithMultiplePermissions(rg *gin.RouterGroup, h *UserHandler, userRepo repo.UserRepo) {
	authMiddleware := UserAuthMiddleware(userRepo)
	
	users := rg.Group("users", authMiddleware)
	{
		// Cần cả "view-users" VÀ "manage-users"
		users.GET("/admin/list", 
			h.ListUsers, 
			UserAllPermissionsMiddleware([]string{"view-users", "manage-users"}, "web"))
		
		// Cần "edit-users" HOẶC "admin" role
		users.PUT("/:id", 
			h.UpdateUser, 
			UserAnyPermissionMiddleware([]string{"edit-users", "update-users"}, "web"))
	}
}

