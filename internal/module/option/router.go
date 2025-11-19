package option

import (
	"github.com/gin-gonic/gin"
)

func registerAPIRoutes(r *gin.RouterGroup, h *OptionHandler, middleware ...gin.HandlerFunc) {
	// Apply middleware to the group if provided
	tx := r.Group("rule", middleware...)
	{
		// tx.GET("types", h.GetRuleTypes)
		// tx.GET("options/tree", h.GetRuleOptionTree)
		tx.GET("options/full", h.GetFullRuleData)
		// tx.GET("types/:typeId/groups", h.GetRuleGroups)
		// tx.GET("groups/:groupId/steps", h.GetRuleOptionSteps)
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
func SetupRoutes(rg *gin.RouterGroup, h *OptionHandler, middleware ...gin.HandlerFunc) {
	registerAPIRoutes(rg, h, middleware...)
}
