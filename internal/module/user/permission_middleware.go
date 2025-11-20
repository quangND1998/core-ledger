package user

import (
	"core-ledger/pkg/database"
	"core-ledger/pkg/ginhp"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserPermissionMiddleware checks if the authenticated user has the required permission
// This middleware should be used AFTER UserAuthMiddleware
func UserPermissionMiddleware(permissionName string, guardName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if response already written (prevent double write)
		if c.Writer.Written() {
			return
		}

		// Get user from context (set by UserAuthMiddleware)
		user, err := GetUserFromContext(c)
		if err != nil {
			if !c.Writer.Written() {
				ginhp.RespondError(c, http.StatusUnauthorized, "User not authenticated")
			}
			c.Abort()
			return
		}

		// Use guard name from user if not provided
		checkGuardName := guardName
		if checkGuardName == "" {
			checkGuardName = user.GuardName
		}

		// Check permission
		hasPermission, err := user.HasPermission(database.Instance(), permissionName, checkGuardName)
		if err != nil {
			if !c.Writer.Written() {
				ginhp.RespondError(c, http.StatusInternalServerError, "Failed to check permission")
			}
			c.Abort()
			return
		}

		if !hasPermission {
			if !c.Writer.Written() {
				ginhp.RespondError(c, http.StatusForbidden, "Insufficient permissions")
			}
			c.Abort()
			return
		}

		c.Next()
	}
}

// UserRoleMiddleware checks if the authenticated user has the required role
// This middleware should be used AFTER UserAuthMiddleware
func UserRoleMiddleware(roleName string, guardName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context (set by UserAuthMiddleware)
		user, err := GetUserFromContext(c)
		if err != nil {
			ginhp.RespondError(c, http.StatusUnauthorized, "User not authenticated")
			c.Abort()
			return
		}

		// Use guard name from user if not provided
		checkGuardName := guardName
		if checkGuardName == "" {
			checkGuardName = user.GuardName
		}

		// Check role
		hasRole, err := user.HasRole(database.Instance(), roleName, checkGuardName)
		if err != nil {
			ginhp.RespondError(c, http.StatusInternalServerError, "Failed to check role")
			c.Abort()
			return
		}

		if !hasRole {
			ginhp.RespondError(c, http.StatusForbidden, "Insufficient role")
			c.Abort()
			return
		}

		c.Next()
	}
}

// UserAnyRoleMiddleware checks if the authenticated user has any of the required roles
// This middleware should be used AFTER UserAuthMiddleware
func UserAnyRoleMiddleware(roleNames []string, guardName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context (set by UserAuthMiddleware)
		user, err := GetUserFromContext(c)
		if err != nil {
			ginhp.RespondError(c, http.StatusUnauthorized, "User not authenticated")
			c.Abort()
			return
		}

		// Use guard name from user if not provided
		checkGuardName := guardName
		if checkGuardName == "" {
			checkGuardName = user.GuardName
		}

		// Check if user has any of the roles
		hasAnyRole, err := user.HasAnyRole(database.Instance(), roleNames, checkGuardName)
		if err != nil {
			ginhp.RespondError(c, http.StatusInternalServerError, "Failed to check roles")
			c.Abort()
			return
		}

		if !hasAnyRole {
			ginhp.RespondError(c, http.StatusForbidden, "Insufficient roles")
			c.Abort()
			return
		}

		c.Next()
	}
}

// UserAllRolesMiddleware checks if the authenticated user has all of the required roles
// This middleware should be used AFTER UserAuthMiddleware
func UserAllRolesMiddleware(roleNames []string, guardName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context (set by UserAuthMiddleware)
		user, err := GetUserFromContext(c)
		if err != nil {
			ginhp.RespondError(c, http.StatusUnauthorized, "User not authenticated")
			c.Abort()
			return
		}

		// Use guard name from user if not provided
		checkGuardName := guardName
		if checkGuardName == "" {
			checkGuardName = user.GuardName
		}

		// Check if user has all of the roles
		hasAllRoles, err := user.HasAllRoles(database.Instance(), roleNames, checkGuardName)
		if err != nil {
			ginhp.RespondError(c, http.StatusInternalServerError, "Failed to check roles")
			c.Abort()
			return
		}

		if !hasAllRoles {
			ginhp.RespondError(c, http.StatusForbidden, "Insufficient roles")
			c.Abort()
			return
		}

		c.Next()
	}
}

// UserAnyPermissionMiddleware checks if the authenticated user has any of the required permissions
// This middleware should be used AFTER UserAuthMiddleware
func UserAnyPermissionMiddleware(permissionNames []string, guardName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context (set by UserAuthMiddleware)
		user, err := GetUserFromContext(c)
		if err != nil {
			ginhp.RespondError(c, http.StatusUnauthorized, "User not authenticated")
			c.Abort()
			return
		}

		// Use guard name from user if not provided
		checkGuardName := guardName
		if checkGuardName == "" {
			checkGuardName = user.GuardName
		}

		// Check if user has any of the permissions
		hasAnyPermission, err := user.HasAnyPermission(database.Instance(), permissionNames, checkGuardName)
		if err != nil {
			ginhp.RespondError(c, http.StatusInternalServerError, "Failed to check permissions")
			c.Abort()
			return
		}

		if !hasAnyPermission {
			ginhp.RespondError(c, http.StatusForbidden, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// UserAllPermissionsMiddleware checks if the authenticated user has all of the required permissions
// This middleware should be used AFTER UserAuthMiddleware
func UserAllPermissionsMiddleware(permissionNames []string, guardName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context (set by UserAuthMiddleware)
		user, err := GetUserFromContext(c)
		if err != nil {
			ginhp.RespondError(c, http.StatusUnauthorized, "User not authenticated")
			c.Abort()
			return
		}

		// Use guard name from user if not provided
		checkGuardName := guardName
		if checkGuardName == "" {
			checkGuardName = user.GuardName
		}

		// Check if user has all of the permissions
		hasAllPermissions, err := user.HasAllPermissions(database.Instance(), permissionNames, checkGuardName)
		if err != nil {
			ginhp.RespondError(c, http.StatusInternalServerError, "Failed to check permissions")
			c.Abort()
			return
		}

		if !hasAllPermissions {
			ginhp.RespondError(c, http.StatusForbidden, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}
