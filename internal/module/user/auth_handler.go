package user

import (
	"core-ledger/internal/module/validate"
	model "core-ledger/model/core-ledger"
	"core-ledger/model/dto"
	"core-ledger/pkg/database"
	"core-ledger/pkg/ginhp"
	"core-ledger/pkg/jwt"
	"core-ledger/pkg/logger"
	"core-ledger/pkg/queue"
	"core-ledger/pkg/repo"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	logger     logger.CustomLogger
	userRepo   repo.UserRepo
	jwtService *jwt.UserJWTService
	dispatcher queue.Dispatcher
}

func NewAuthHandler(userRepo repo.UserRepo, dispatcher queue.Dispatcher) *AuthHandler {
	return &AuthHandler{
		logger:     logger.NewSystemLog("AuthHandler"),
		userRepo:   userRepo,
		jwtService: jwt.NewUserJWTService(),
		dispatcher: dispatcher,
	}
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest represents the register request body
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FullName  string `json:"full_name"`
	GuardName string `json:"guard_name"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         struct {
		ID          uint64             `json:"id"`
		Email       string             `json:"email"`
		FullName    string             `json:"full_name"`
		GuardName   string             `json:"guard_name"`
		Roles       []model.Role       `json:"roles"`
		Permissions []model.Permission `json:"permissions"`
	} `json:"user"`
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		out := validate.FormatErrorMessage(req, err)
		ginhp.RespondErrorValidate(c, http.StatusUnprocessableEntity, "Invalid input", out)

		return
	}

	// Find user by email
	user, err := h.userRepo.GetOneByFields(c, map[string]interface{}{"email": req.Email})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ginhp.RespondError(c, http.StatusUnauthorized, "Invalid email or password")
		} else {
			ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		ginhp.RespondError(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Generate tokens
	accessToken, err := h.jwtService.GenerateToken(user)
	if err != nil {
		ginhp.RespondError(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	refreshToken, err := h.jwtService.GenerateRefreshToken(user)
	if err != nil {
		ginhp.RespondError(c, http.StatusInternalServerError, "Failed to generate refresh token")
		return
	}

	// Calculate expiration time
	expiresAt := time.Now().Add(24 * time.Hour) // Default 24 hours

	// Load roles and permissions
	roles, err := h.loadUserRoles(user)
	if err != nil {
		h.logger.Error("Failed to load user roles", err)
		// Continue without roles (don't fail login)
		roles = []model.Role{}
	}

	permissions, err := user.GetAllPermissions(database.Instance(), user.GuardName)
	if err != nil {
		h.logger.Error("Failed to load user permissions", err)
		// Continue without permissions (don't fail login)
		permissions = []model.Permission{}
	}

	response := AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}
	response.User.ID = user.ID
	response.User.Email = user.Email
	response.User.FullName = user.FullName
	response.User.GuardName = user.GuardName
	response.User.Roles = roles
	response.User.Permissions = permissions

	c.JSON(http.StatusOK, dto.PreResponse{
		Data: response,
	})
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Check if user already exists
	existing, err := h.userRepo.GetOneByFields(c, map[string]interface{}{"email": req.Email})
	if err == nil && existing != nil {
		ginhp.RespondError(c, http.StatusConflict, "Email already registered")
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		ginhp.RespondError(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Set default guard name
	guardName := req.GuardName
	if guardName == "" {
		guardName = "web"
	}

	// Create user
	user := &model.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		FullName:  req.FullName,
		GuardName: guardName,
	}

	if err := h.userRepo.Create(user); err != nil {
		ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Generate tokens
	accessToken, err := h.jwtService.GenerateToken(user)
	if err != nil {
		ginhp.RespondError(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	refreshToken, err := h.jwtService.GenerateRefreshToken(user)
	if err != nil {
		ginhp.RespondError(c, http.StatusInternalServerError, "Failed to generate refresh token")
		return
	}

	expiresAt := time.Now().Add(24 * time.Hour)

	// Load roles and permissions
	roles, err := h.loadUserRoles(user)
	if err != nil {
		h.logger.Error("Failed to load user roles", err)
		roles = []model.Role{}
	}

	permissions, err := user.GetAllPermissions(database.Instance(), user.GuardName)
	if err != nil {
		h.logger.Error("Failed to load user permissions", err)
		permissions = []model.Permission{}
	}

	response := AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}
	response.User.ID = user.ID
	response.User.Email = user.Email
	response.User.FullName = user.FullName
	response.User.GuardName = user.GuardName
	response.User.Roles = roles
	response.User.Permissions = permissions

	c.JSON(http.StatusCreated, dto.PreResponse{
		Data: response,
	})
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Validate refresh token
	claims, err := h.jwtService.ValidateToken(req.RefreshToken)
	if err != nil {
		ginhp.RespondError(c, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	// Get user
	user, err := h.userRepo.GetByID(c, int64(claims.UserID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ginhp.RespondError(c, http.StatusUnauthorized, "User not found")
		} else {
			ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// Generate new tokens
	accessToken, err := h.jwtService.GenerateToken(user)
	if err != nil {
		ginhp.RespondError(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	refreshToken, err := h.jwtService.GenerateRefreshToken(user)
	if err != nil {
		ginhp.RespondError(c, http.StatusInternalServerError, "Failed to generate refresh token")
		return
	}

	expiresAt := time.Now().Add(24 * time.Hour)

	// Load roles and permissions
	roles, err := h.loadUserRoles(user)
	if err != nil {
		h.logger.Error("Failed to load user roles", err)
		roles = []model.Role{}
	}

	permissions, err := user.GetAllPermissions(database.Instance(), user.GuardName)
	if err != nil {
		h.logger.Error("Failed to load user permissions", err)
		permissions = []model.Permission{}
	}

	response := AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}
	response.User.ID = user.ID
	response.User.Email = user.Email
	response.User.FullName = user.FullName
	response.User.GuardName = user.GuardName
	response.User.Roles = roles
	response.User.Permissions = permissions

	c.JSON(http.StatusOK, dto.PreResponse{
		Data: response,
	})
}

// LogoutRequest represents the logout request body
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// LogoutResponse represents the logout response
type LogoutResponse struct {
	Message string `json:"message"`
}

// Logout handles user logout and blacklists tokens
func (h *AuthHandler) Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Logout can work without body, but if body is provided, validate it
		// For now, we'll make refresh_token optional
	}

	var userID uint64

	// Get access token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		// Extract token from "Bearer {token}"
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			accessToken := parts[1]
			// Parse access token without blacklist check (to get claims)
			claims, err := h.jwtService.ValidateTokenWithoutBlacklistCheck(accessToken)
			if err == nil {
				userID = claims.UserID
				// Blacklist access token
				if err := h.jwtService.BlacklistTokenByClaims(accessToken, claims); err != nil {
					h.logger.Error("Failed to blacklist access token", err)
				} else {
					h.logger.Info("Access token blacklisted", map[string]interface{}{
						"user_id": userID,
					})
				}
			}
		}
	}

	// Blacklist refresh token if provided
	if req.RefreshToken != "" {
		claims, err := h.jwtService.ValidateTokenWithoutBlacklistCheck(req.RefreshToken)
		if err == nil {
			if userID == 0 {
				userID = claims.UserID
			}
			// Blacklist refresh token
			if err := h.jwtService.BlacklistTokenByClaims(req.RefreshToken, claims); err != nil {
				h.logger.Error("Failed to blacklist refresh token", err)
			} else {
				h.logger.Info("Refresh token blacklisted", map[string]interface{}{
					"user_id": userID,
				})
			}
		} else {
			h.logger.Warn("Invalid refresh token provided during logout", err)
		}
	}

	// Log logout event for audit purposes
	if userID > 0 {
		h.logger.Info("User logged out", map[string]interface{}{
			"user_id": userID,
		})
	}

	// Return success response
	response := LogoutResponse{
		Message: "Logout successful. Tokens have been revoked.",
	}

	c.JSON(http.StatusOK, dto.PreResponse{
		Data: response,
	})
}

// loadUserRoles loads all roles for a user
func (h *AuthHandler) loadUserRoles(user *model.User) ([]model.Role, error) {
	var roles []model.Role
	err := database.Instance().
		Model(&model.ModelHasRole{}).
		Joins("JOIN roles ON model_has_roles.role_id = roles.id").
		Where("model_has_roles.model_id = ? AND model_has_roles.model_type = ?", user.ID, "User").
		Where("roles.guard_name = ?", user.GuardName).
		Select("roles.*").
		Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}


