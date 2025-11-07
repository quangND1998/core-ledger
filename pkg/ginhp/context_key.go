package ginhp

import (
	"github.com/gin-gonic/gin"
)

type ContextKey string

const (
	ContextKeyAccountRequest  ContextKey = "accountRequest"
	ContextKeyCustomerRequest ContextKey = "customer_request"
	ContextKeyEmployeeRequest ContextKey = "employee_request"
	ContextKeyFingerprint     ContextKey = "Fingerprint"
)

func (t ContextKey) String() string {
	return string(t)
}

func GetByKey[T any](c *gin.Context, key ContextKey) T {
	return c.MustGet(string(key)).(T)
}
