package middleware

import (
	"github.com/gin-gonic/gin"
	"core-ledger/internal/core"
	"core-ledger/pkg/ginhp"
)

func (m *middleware) WalletGuard(c *gin.Context) {
	customer := ginhp.GetCustomerReq(c)
	if customer == nil {
		return
	}
	if !customer.WealifyWalletEnable {
		ginhp.RespondOKWithError(c, core.NewError(core.ErrCodeWealifyWalletFeatureInactive))
		return
	}
}

func (m *middleware) VAGuard(c *gin.Context) {
	customer := ginhp.GetCustomerReq(c)
	if customer == nil {
		return
	}
	if !customer.VaEnable {
		ginhp.RespondOKWithError(c, core.NewError(core.ErrCodeVAFeatureInactive))
		return
	}
}
