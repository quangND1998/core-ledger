package middleware

import (
	"bytes"
	"context"
	config "core-ledger/configs"
	"core-ledger/model/dto"
	"core-ledger/pkg/ginhp"
	"core-ledger/pkg/logger"
	"core-ledger/pkg/repo"
	"core-ledger/pkg/utils/fingerprint"
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mssola/user_agent"
)

type MiddleWare interface {
	Authorize(c *gin.Context)
	Recovery(c *gin.Context)
	BasicAuth(c *gin.Context)
	ExtractFingerprint(c *gin.Context)
	VAGuard(c *gin.Context)
	WalletGuard(c *gin.Context)
}

type middleware struct {
	l            logger.CustomLogger
	customerRepo repo.CustomerRepo
	employeeRepo repo.EmployeeRepo
	basicAuth    map[string]string
}

func (m *middleware) ExtractFingerprint(c *gin.Context) {
	ip := c.GetHeader("X-Forward-For")
	if ip == "" {
		ip = c.ClientIP()
	}

	// Lấy các headers
	accept := c.GetHeader("Accept")
	acceptLang := c.GetHeader("Accept-Language")
	userAgentString := c.GetHeader("User-Agent")

	// Parse user-agent
	ua := user_agent.New(userAgentString)
	browserName, browserVersion := ua.Browser()
	os := ua.OS()
	device := ua.Platform()

	// Optional: tách OS version nếu muốn chia thành major/minor
	osMajor := ""
	osMinor := ""
	if parts := strings.Split(os, " "); len(parts) >= 2 {
		osMajor = parts[1]
		if len(parts) > 2 {
			osMinor = parts[2]
		}
	}

	c.Set(ginhp.ContextKeyFingerprint.String(), fingerprint.Fingerprint{
		Headers: fingerprint.AcceptHeader{
			Accept:          accept,
			AcceptLanguage:  acceptLang,
			UserAgentHeader: userAgentString,
		},
		IpAddress: fingerprint.IpAddress{
			Value: ip,
		},
		UserAgent: fingerprint.UserAgent{
			Device: fingerprint.DeviceInfo{
				Family:  device,
				Version: []string{}, // `ua-parser` không parse version cụ thể cho device
			},
			OS: fingerprint.OSInfo{
				Family: os,
				Major:  osMajor,
				Minor:  osMinor,
			},
			Browser: fingerprint.BrowserInfo{
				Family:  browserName,
				Version: browserVersion,
			},
		},
	})
}

// AuthMiddleware creates a new Gin middleware for JWT authentication.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := extractTokenFromHeader(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		claims := &dto.Claims{}

		jwtSecret := config.GetConfig().JWT.Secret
		if jwtSecret == "" {
			jwtSecret = "default-secret-key"
		}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil {
			if errors.Is(err, jwt.ErrSignatureInvalid) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token signature"})
				c.Abort()
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set user information in the context for downstream handlers
		c.Set("userID", claims.ID)
		c.Set("isEmployee", claims.IsEmployee)

		c.Next()
	}
}

func extractTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("Authorization header is required")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("Authorization header format must be Bearer {token}")
	}

	return parts[1], nil
}

// GetUserIDFromContext retrieves the user SystemPaymentID from the Gin context.
func GetUserIDFromContext(c *gin.Context) (uint64, error) {
	userID, ok := c.Get("userID")
	if !ok {
		return 0, errors.New("userID not found in context")
	}
	fmt.Printf("Type of myInt: %T\n", userID)
	id, ok := userID.(uint64)
	if !ok {
		return 0, errors.New("userID in context is not of type uint")
	}

	return id, nil
}

// EmployeeAuthMiddleware checks for a valid JWT and ensures the user is an employee.
func EmployeeAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// First, run the standard authentication check.
		AuthMiddleware()(c)

		// If the standard auth middleware aborted, c.IsAborted() will be true.
		if c.IsAborted() {
			return
		}

		// Now, check if the user is an employee.
		isEmployee, exists := c.Get("isEmployee")
		if !exists || !isEmployee.(bool) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Employee access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (m *middleware) BasicAuth(c *gin.Context) {
	username, password, ok := c.Request.BasicAuth()
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code": http.StatusUnauthorized,
		})
	}
	if pwd, ok := m.basicAuth[username]; !ok || pwd != password {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code": http.StatusUnauthorized,
		})
	}
	c.Next()
}

func (w *CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)    // ghi vào buffer
	return len(b), nil // vẫn ghi ra luôn
}

func (m *middleware) Authorize(c *gin.Context) {
	apiKey := c.GetHeader("x-api-key")
	if apiKey != "" {
		m.authorizeApiKey(c, apiKey)
	} else {
		bearerToken := c.GetHeader("Authorization")
		parts := strings.Split(bearerToken, " ")
		if len(parts) != 2 {
			ginhp.RespondError(c, http.StatusUnauthorized, "invalid token paths")
			return
		}
		claims := &dto.Claims{}
		token, err := jwt.ParseWithClaims(parts[1], claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		logrus.Info(*claims)
		logrus.Info(token)
		if err != nil {
			logrus.Info(err)
			ginhp.RespondError(c, http.StatusUnauthorized, "Unauthorized")
			return
		}
		accountReq := &ginhp.AccountRequest{
			CallingCode:     ginhp.CallingCodeRelation{},
			TwoFactorStatus: claims.TwoFactorStatus,
			RegisteredAt:    time.Now(),
			Status:          true,
			IsDeleted:       false,
			IsEmployee:      claims.IsEmployee,
		}
		var customerReq *ginhp.CustomerRequest
		var employeeReq *ginhp.EmployeeRequest
		if claims.IsEmployee {
			accountReq.Employee, err = m.employeeRepo.GetOneByFields(c, map[string]interface{}{
				"id": claims.ID,
			})
			if err != nil {
				ginhp.ReturnBadRequestError(c, err)
				return
			}
			employeeReq = &ginhp.EmployeeRequest{
				ID: claims.ID,
			}
		} else {
			accountReq.Customer, err = m.customerRepo.GetOneByFields(c, map[string]interface{}{
				"id": claims.ID,
			})
			if err != nil {
				ginhp.ReturnBadRequestError(c, err)
				return
			}
			var delegationAccountID *string
			if claims.DelegationAccountID != "" {
				delegationAccountID = &claims.DelegationAccountID
			}
			customerReq = &ginhp.CustomerRequest{
				Customer:            accountReq.Customer,
				DelegationAccountID: delegationAccountID,
			}
		}

		c.Set(ginhp.ContextKeyAccountRequest.String(), accountReq)
		c.Set(ginhp.ContextKeyCustomerRequest.String(), customerReq)
		c.Set(ginhp.ContextKeyEmployeeRequest.String(), employeeReq)
	}

}

func (m *middleware) Recovery(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			m.l.Error(fmt.Sprintf("🔥 Panic recovered: %v\n%s", err, debug.Stack()))
			ginhp.RespondError(c, http.StatusInternalServerError, err.(error).Error())
		}
	}()
	c.Next()
}

func (m *middleware) authorizeApiKey(c *gin.Context, apiKey string) {
	customer, err := m.customerRepo.GetByApiKey(context.Background(), apiKey)
	if err != nil {
		ginhp.RespondError(c, http.StatusUnauthorized, "Invalid api key")
	}
	accountReq := &ginhp.AccountRequest{
		CallingCode:  ginhp.CallingCodeRelation{},
		RegisteredAt: time.Now(),
		Status:       true,
		IsDeleted:    false,
	}
	accountReq.Customer = customer
	customerReq := &ginhp.CustomerRequest{
		Customer: accountReq.Customer,
	}
	c.Set(ginhp.ContextKeyAccountRequest.String(), accountReq)
	c.Set(ginhp.ContextKeyCustomerRequest.String(), customerReq)
}

func NewMiddleware(customerRepo repo.CustomerRepo, employeeRepo repo.EmployeeRepo) MiddleWare {
	return &middleware{
		l:            logger.NewSystemLog("Middleware"),
		customerRepo: customerRepo,
		employeeRepo: employeeRepo,
	}
}
