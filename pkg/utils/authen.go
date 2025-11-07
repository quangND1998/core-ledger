package utils

import (
	"core-ledger/pkg/ginhp"
	"encoding/base64"
	"errors"

	"github.com/gin-gonic/gin"
)

func WithAuthenticatedUser(c *gin.Context) (userID int64, err error) {
	rawCustomer, exist := c.Get(ginhp.ContextKeyAccountRequest.String())
	if !exist {
		return 0, errors.New("account not found in request")
	}
	customer, ok := rawCustomer.(*ginhp.AccountRequest)
	if !ok {
		return 0, errors.New("invalid customer type")
	}
	return customer.Customer.ID, nil
}

func GenerateBasicAuth(username, password string) string {
	authString := username + ":" + password
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(authString))
	return "Basic " + encodedAuth
}
