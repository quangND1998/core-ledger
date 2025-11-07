package constants

import "github.com/gin-gonic/gin"

type ConstantHandler interface {
	GetConstants(c *gin.Context)
	GetAccountLevels(c *gin.Context)
	GetAccountTypes(c *gin.Context)
	GetProviders(c *gin.Context)
	GetProviderTypes(c *gin.Context)
}
