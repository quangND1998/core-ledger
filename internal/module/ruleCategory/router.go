package ruleCategory

import (
	"github.com/gin-gonic/gin"
)

func registerAPIRoutes(r *gin.RouterGroup, h *RuleCategoryHandler, middleware ...gin.HandlerFunc) {
	// Apply middleware to the group if provided
	tx := r.Group("rule-category", middleware...)
	{
		tx.GET("/list", h.List)

		// Add more routes here
		// tx.POST("", h.Create)
		// tx.GET("/:id", h.GetByID)
		// tx.PUT("/:id", h.Update)
		// tx.DELETE("/:id", h.Delete)
	}
}

// SetupRoutes registers transaction routes with optional middleware
// Usage:
//   - Without middleware: transactions.SetupRoutes(protected, handler)
//   - With middleware: transactions.SetupRoutes(protected, handler, authMiddleware, loggingMiddleware)
func SetupRoutes(rg *gin.RouterGroup, h *RuleCategoryHandler, middleware ...gin.HandlerFunc) {
	registerAPIRoutes(rg, h, middleware...)
}
