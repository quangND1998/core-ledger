package user

import (
	"github.com/gin-gonic/gin"
)

func registerAuthRoutes(r *gin.RouterGroup, h *AuthHandler) {
	auth := r.Group("auth")
	{
		// Public auth endpoints
		auth.POST("/login", h.Login)
		auth.POST("/register", h.Register)
		auth.POST("/refresh", h.RefreshToken)
		auth.POST("/logout", h.Logout)
	}
}

// SetupAuthRoutes registers authentication routes
func SetupAuthRoutes(rg *gin.RouterGroup, h *AuthHandler) {
	registerAuthRoutes(rg, h)
}
