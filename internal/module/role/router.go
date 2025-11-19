package role

import (
	"github.com/gin-gonic/gin"
)

func registerAPIRoutes(r *gin.RouterGroup, h *RoleHandler, middleware ...gin.HandlerFunc) {
	// Apply middleware to the group if provided
	roles := r.Group("roles", middleware...)
	{
		// Role management - CRUD
		roles.GET("/list", h.ListRoles)
		roles.GET("/:id", h.GetRole)
		roles.POST("", h.CreateRole)
		roles.PUT("/:id", h.UpdateRole)
		roles.DELETE("/:id", h.DeleteRole)

		// Role permission management
		roles.POST("/give-permission", h.GivePermissionToRole)
		roles.POST("/revoke-permission", h.RevokePermissionFromRole)
		roles.POST("/:id/permissions", h.SyncRolePermissions)
		roles.GET("/:id/permissions", h.GetRolePermissions)
	}
}

// SetupRoutes registers role routes with optional middleware
// Usage:
//   - Without middleware: role.SetupRoutes(protected, handler)
//   - With middleware: role.SetupRoutes(protected, handler, authMiddleware, permissionMiddleware)
func SetupRoutes(rg *gin.RouterGroup, h *RoleHandler, middleware ...gin.HandlerFunc) {
	registerAPIRoutes(rg, h, middleware...)
}

