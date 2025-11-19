package middleware

import (
	model "core-ledger/model/core-ledger"
	"core-ledger/pkg/database"
	"core-ledger/pkg/repo"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PermissionMiddleware checks if the authenticated user has the required permission
func PermissionMiddleware(permissionName string, guardName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (set by AuthMiddleware)
		userID, ok := c.Get("userID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		userIDInt64, ok := userID.(int64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
			c.Abort()
			return
		}

		// Get user from database
		userRepo := repo.NewUserRepoIml(database.Instance())
		user, err := userRepo.GetByID(c, userIDInt64)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check permission"})
			}
			c.Abort()
			return
		}

		// Check permission
		hasPermission, err := user.HasPermission(database.Instance(), permissionName, guardName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check permission"})
			c.Abort()
			return
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		// Store user in context for handlers
		c.Set("user", user)
		c.Next()
	}
}

// RoleMiddleware checks if the authenticated user has the required role
func RoleMiddleware(roleName string, guardName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (set by AuthMiddleware)
		userID, ok := c.Get("userID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		userIDInt64, ok := userID.(int64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
			c.Abort()
			return
		}

		// Get user from database
		userRepo := repo.NewUserRepoIml(database.Instance())
		user, err := userRepo.GetByID(c, userIDInt64)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check role"})
			}
			c.Abort()
			return
		}

		// Check role
		hasRole, err := user.HasRole(database.Instance(), roleName, guardName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check role"})
			c.Abort()
			return
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient role"})
			c.Abort()
			return
		}

		// Store user in context for handlers
		c.Set("user", user)
		c.Next()
	}
}

// AnyRoleMiddleware checks if the authenticated user has any of the required roles
func AnyRoleMiddleware(roleNames []string, guardName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (set by AuthMiddleware)
		userID, ok := c.Get("userID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		userIDInt64, ok := userID.(int64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
			c.Abort()
			return
		}

		// Get user from database
		userRepo := repo.NewUserRepoIml(database.Instance())
		user, err := userRepo.GetByID(c, userIDInt64)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check roles"})
			}
			c.Abort()
			return
		}

		// Check if user has any of the roles
		hasAnyRole, err := user.HasAnyRole(database.Instance(), roleNames, guardName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check roles"})
			c.Abort()
			return
		}

		if !hasAnyRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient roles"})
			c.Abort()
			return
		}

		// Store user in context for handlers
		c.Set("user", user)
		c.Next()
	}
}

// GetUserFromContext retrieves the user from Gin context
func GetUserFromContext(c *gin.Context) (*model.User, error) {
	user, ok := c.Get("user")
	if !ok {
		return nil, errors.New("user not found in context")
	}

	userModel, ok := user.(*model.User)
	if !ok {
		return nil, errors.New("invalid user type in context")
	}

	return userModel, nil
}

