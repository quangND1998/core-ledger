package coaaccount

import (
	"github.com/gin-gonic/gin"
)

func registerAPIRoutes(r *gin.RouterGroup, h *CoaAccountHandler, middleware ...gin.HandlerFunc) {
	// Apply middleware to the group if provided
	tx := r.Group("coa-accounts", middleware...)
	{
		tx.GET("/list", h.List)
		// tx.GET("/list", user.UserAnyPermissionMiddleware([]string{"coa.view"}, "web"), h.List)
		tx.GET("/:id", h.GetCoaAccountDetail)
		tx.POST("export", h.ExportCoaAccounts)
		// Add more routes here
		tx.POST("", h.Create)
		tx.POST("/:account_no/exist", h.ExistAccoutnNo)
		tx.PUT("/:id/update", h.Update)
		tx.POST("update-status", h.UpdateStatus)
		// tx.GET("/:id", h.GetByID)
		// tx.PUT("/:id", h.Update)
		// tx.DELETE("/:id", h.Delete)
	}
}

// RegisterRequestCoaAccountRoutes registers routes for request_coa_account
func RegisterRequestCoaAccountRoutes(r *gin.RouterGroup, h *RequestCoaAccountHandler, middleware ...gin.HandlerFunc) {
	req := r.Group("request-coa-accounts", middleware...)
	{
		req.GET("", h.GetList)              // Get list of requests
		req.GET("/:id", h.GetDetail)        // Get request detail
		req.POST("", h.Create)              // Create new request
		req.PUT("/:id", h.Update)           // Update rejected request
		req.POST("/:id/approve", h.Approve) // Approve request
		req.POST("/:id/reject", h.Reject)   // Reject request
	}
}

// SetupRoutes registers transaction routes with optional middleware
// Usage:
//   - Without middleware: transactions.SetupRoutes(protected, handler)
//   - With middleware: transactions.SetupRoutes(protected, handler, authMiddleware, loggingMiddleware)
func SetupRoutes(rg *gin.RouterGroup, h *CoaAccountHandler, hr *RequestCoaAccountHandler, middleware ...gin.HandlerFunc) {
	registerAPIRoutes(rg, h, middleware...)
	RegisterRequestCoaAccountRoutes(rg, hr, middleware...)
}
