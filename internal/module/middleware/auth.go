package middleware

import (
	"bytes"
	"core-ledger/pkg/ginhp"
	"core-ledger/pkg/logger"
	"core-ledger/pkg/utils/fingerprint"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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
	l logger.CustomLogger

	basicAuth map[string]string
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
func GetUserIDFromContext(c *gin.Context) (int64, error) {
	userID, ok := c.Get("userID")
	if !ok {
		return 0, errors.New("userID not found in context")
	}

	id, ok := userID.(int64)
	if !ok {
		return 0, errors.New("userID in context is not of type uint")
	}

	return id, nil
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
