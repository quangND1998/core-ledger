package permission

import (
	"github.com/gin-gonic/gin"
)

func registerAPIRoutes(r *gin.RouterGroup, h *PermissionHandler, middleware ...gin.HandlerFunc) {
	// Apply middleware to the group if provided
	perm := r.Group("permissions", middleware...)
	{
		// Permission management - CRUD
		perm.GET("/list", h.ListPermissions)
		perm.GET("/:id", h.GetPermission)
		perm.POST("", h.CreatePermission)
		perm.PUT("/:id", h.UpdatePermission)
		perm.DELETE("/:id", h.DeletePermission)
	}
}

// SetupRoutes registers permission routes with optional middleware
// Usage:
//   - Without middleware: permission.SetupRoutes(protected, handler)
//   - With middleware: permission.SetupRoutes(protected, handler, authMiddleware, permissionMiddleware)
func SetupRoutes(rg *gin.RouterGroup, h *PermissionHandler, middleware ...gin.HandlerFunc) {
	registerAPIRoutes(rg, h, middleware...)
}
