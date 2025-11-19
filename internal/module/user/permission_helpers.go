package user

import (
	"core-ledger/pkg/database"
	"core-ledger/pkg/ginhp"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CheckPermission checks if the current user has a specific permission
// Use this in handlers to check permissions programmatically
func CheckPermission(c *gin.Context, permissionName string, guardName string) (bool, error) {
	user, err := GetUserFromContext(c)
	if err != nil {
		return false, err
	}

	checkGuardName := guardName
	if checkGuardName == "" {
		checkGuardName = user.GuardName
	}

	return user.HasPermission(database.Instance(), permissionName, checkGuardName)
}

// CheckRole checks if the current user has a specific role
// Use this in handlers to check roles programmatically
func CheckRole(c *gin.Context, roleName string, guardName string) (bool, error) {
	user, err := GetUserFromContext(c)
	if err != nil {
		return false, err
	}

	checkGuardName := guardName
	if checkGuardName == "" {
		checkGuardName = user.GuardName
	}

	return user.HasRole(database.Instance(), roleName, checkGuardName)
}

// RequirePermission is a helper that checks permission and returns error if not authorized
// Use this in handlers when you need to check permission and return error
func RequirePermission(c *gin.Context, permissionName string, guardName string) error {
	hasPermission, err := CheckPermission(c, permissionName, guardName)
	if err != nil {
		ginhp.RespondError(c, http.StatusInternalServerError, "Failed to check permission")
		return err
	}

	if !hasPermission {
		ginhp.RespondError(c, http.StatusForbidden, "Insufficient permissions")
		return errors.New("insufficient permissions")
	}

	return nil
}

// RequireRole is a helper that checks role and returns error if not authorized
// Use this in handlers when you need to check role and return error
func RequireRole(c *gin.Context, roleName string, guardName string) error {
	hasRole, err := CheckRole(c, roleName, guardName)
	if err != nil {
		ginhp.RespondError(c, http.StatusInternalServerError, "Failed to check role")
		return err
	}

	if !hasRole {
		ginhp.RespondError(c, http.StatusForbidden, "Insufficient role")
		return errors.New("insufficient role")
	}

	return nil
}

