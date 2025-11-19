package user

import (
	model "core-ledger/model/core-ledger"
	"core-ledger/pkg/ginhp"
	"core-ledger/pkg/jwt"
	"core-ledger/pkg/repo"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserAuthMiddleware creates a middleware for JWT authentication
func UserAuthMiddleware(userRepo repo.UserRepo) gin.HandlerFunc {
	jwtService := jwt.NewUserJWTService()

	return func(c *gin.Context) {
		tokenString, err := extractTokenFromHeader(c)
		if err != nil {
			ginhp.RespondError(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		// Validate token
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			ginhp.RespondError(c, http.StatusUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}

		// Get user from database
		user, err := userRepo.GetByID(c, int64(claims.UserID))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ginhp.RespondError(c, http.StatusUnauthorized, "User not found")
			} else {
				ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
			}
			c.Abort()
			return
		}

		// Set user in context
		c.Set("user", user)
		c.Set("userID", user.ID)
		c.Set("userEmail", user.Email)
		c.Set("guardName", user.GuardName)

		c.Next()
	}
}

// extractTokenFromHeader extracts the JWT token from the Authorization header
func extractTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("Authorization header is required")
	}

	// Check if it starts with "Bearer "
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return "", errors.New("Authorization header must start with 'Bearer '")
	}

	return authHeader[7:], nil
}

// GetUserFromContext retrieves the user from the Gin context
func GetUserFromContext(c *gin.Context) (*model.User, error) {
	user, exists := c.Get("user")
	if !exists {
		return nil, errors.New("user not found in context")
	}

	u, ok := user.(*model.User)
	if !ok {
		return nil, errors.New("invalid user type in context")
	}

	return u, nil
}

// GetUserIDFromContext retrieves the user ID from the Gin context
func GetUserIDFromContext(c *gin.Context) (uint64, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, errors.New("userID not found in context")
	}

	id, ok := userID.(uint64)
	if !ok {
		return 0, errors.New("invalid userID type in context")
	}

	return id, nil
}

