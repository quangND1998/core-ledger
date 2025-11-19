package user

import (
	"github.com/gin-gonic/gin"
)

func registerAPIRoutes(r *gin.RouterGroup, h *UserHandler, middleware ...gin.HandlerFunc) {
	// Apply middleware to the group if provided
	users := r.Group("users", middleware...)
	{
		// User management - CRUD
		users.GET("/list", h.ListUsers)
		users.GET("/:id", h.GetUser)
		users.POST("", h.CreateUser)
		users.PUT("/:id", h.UpdateUser)

		// User role & permission management
		users.POST("/:id/roles", h.SyncUserRoles)
		users.POST("/:id/permissions", h.SyncUserPermissions)

		// Legacy endpoints (for backward compatibility)
		users.POST("/give-permission", h.GivePermissionToUser)
		users.POST("/revoke-permission", h.RevokePermissionFromUser)
		users.POST("/assign-role", h.AssignRoleToUser)
		users.POST("/remove-role", h.RemoveRoleFromUser)
	}
}

// SetupRoutes registers user routes with optional middleware
// Usage:
//   - Without middleware: user.SetupRoutes(protected, handler)
//   - With middleware: user.SetupRoutes(protected, handler, authMiddleware, permissionMiddleware)
func SetupRoutes(rg *gin.RouterGroup, h *UserHandler, middleware ...gin.HandlerFunc) {
	registerAPIRoutes(rg, h, middleware...)
}

