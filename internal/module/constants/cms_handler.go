package constants

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type cmsConstantHandler struct{}

func NewCmsConstantHandler() ConstantHandler {
	return &cmsConstantHandler{}
}

// GetConstants godoc
// @Summary Get all constants
// @Description Get all system constants
// @Tags constants
// @Accept jsonfield
// @Produce jsonfield
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /constants [get]
func (h *cmsConstantHandler) GetConstants(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   Constants,
	})
}

// GetAccountLevels godoc
// @Summary Get account levels
// @Description Get account level constants
// @Tags constants
// @Accept jsonfield
// @Produce jsonfield
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /constants/account-levels [get]
func (h *cmsConstantHandler) GetAccountLevels(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   Constants["account_level"],
	})
}

// GetAccountTypes godoc
// @Summary Get account types
// @Description Get account type constants
// @Tags constants
// @Accept jsonfield
// @Produce jsonfield
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /constants/account-types [get]
func (h *cmsConstantHandler) GetAccountTypes(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   Constants["account_type"],
	})
}

// GetProviders godoc
// @Summary Get providers
// @Description Get provider constants
// @Tags constants
// @Accept jsonfield
// @Produce jsonfield
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /constants/providers [get]
func (h *cmsConstantHandler) GetProviders(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   Constants["provider"],
	})
}

// GetProviderTypes godoc
// @Summary Get provider types
// @Description Get provider type constants
// @Tags constants
// @Accept jsonfield
// @Produce jsonfield
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /constants/provider-types [get]
func (h *cmsConstantHandler) GetProviderTypes(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   Constants["provider_type"],
	})
}
