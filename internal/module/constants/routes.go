package constants

import (
	"github.com/gin-gonic/gin"
)

func registerAPIRoutes(r *gin.RouterGroup, h ConstantHandler, middleware ...gin.HandlerFunc) {
	cc := r.Group("constants", middleware...)
	{
		cc.GET("", h.GetConstants)
		cc.GET("account-levels", h.GetAccountLevels)
		cc.GET("account-types", h.GetAccountTypes)
		cc.GET("providers", h.GetProviders)
		cc.GET("providers-types", h.GetProviderTypes)
	}
}

func registerCMSRoutes(r *gin.RouterGroup, h ConstantHandler, middleware ...gin.HandlerFunc) {
	cc := r.Group("cms/constants", middleware...)
	{
		cc.GET("", h.GetConstants)
		cc.GET("account-levels", h.GetAccountLevels)
		cc.GET("account-types", h.GetAccountTypes)
		cc.GET("providers", h.GetProviders)
		cc.GET("providers-types", h.GetProviderTypes)
	}
}

func SetupConstantRoutes(router *gin.RouterGroup, handler, cmsHandler ConstantHandler, middleware ...gin.HandlerFunc) {
	registerAPIRoutes(router, handler, middleware...)
	registerCMSRoutes(router, cmsHandler, middleware...)
}
