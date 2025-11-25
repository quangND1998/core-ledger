package swagger

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes registers Swagger documentation routes
func SetupRoutes(rg *gin.RouterGroup, h *SwaggerHandler) {
	swagger := rg.Group("swagger")
	{
		swagger.GET("/swagger.json", h.GetSwaggerJSON)
		swagger.GET("/", h.GetSwaggerUI)
		swagger.GET("/redoc", h.GetSwaggerUIAlternative)
	}
}


