package jwt

import (
	config "core-ledger/configs"
	"core-ledger/model/core-ledger"
	"core-ledger/pkg/database"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// UserClaims represents the JWT claims for User authentication
type UserClaims struct {
	UserID    uint64 `json:"user_id"`
	Email     string `json:"email"`
	GuardName string `json:"guard_name"`
	jwt.RegisteredClaims
}

// UserJWTService handles JWT operations for User authentication
type UserJWTService struct {
	secret     string
	expiresIn  time.Duration
	refreshIn  time.Duration
}

// NewUserJWTService creates a new UserJWTService
func NewUserJWTService() *UserJWTService {
	cfg := config.GetConfig()
	secret := cfg.JWT.Secret
	if secret == "" {
		secret = "default-user-jwt-secret-key-change-in-production"
	}

	expiresIn := time.Hour * 24 // Default 24 hours
	if cfg.JWT.ExpiresIn > 0 {
		expiresIn = time.Duration(cfg.JWT.ExpiresIn) * time.Second
	}

	refreshIn := time.Hour * 24 * 7 // Default 7 days

	return &UserJWTService{
		secret:    secret,
		expiresIn: expiresIn,
		refreshIn: refreshIn,
	}
}

// GenerateToken generates a JWT token for a user
func (s *UserJWTService) GenerateToken(user *model.User) (string, error) {
	now := time.Now()
	claims := &UserClaims{
		UserID:    user.ID,
		Email:     user.Email,
		GuardName: user.GuardName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.expiresIn)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "core-ledger",
			Subject:   user.Email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

// GenerateRefreshToken generates a refresh token for a user
func (s *UserJWTService) GenerateRefreshToken(user *model.User) (string, error) {
	now := time.Now()
	claims := &UserClaims{
		UserID:    user.ID,
		Email:     user.Email,
		GuardName: user.GuardName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshIn)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "core-ledger",
			Subject:   user.Email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

// ValidateToken validates a JWT token and returns the claims
func (s *UserJWTService) ValidateToken(tokenString string) (*UserClaims, error) {
	// Check if token is blacklisted first
	isBlacklisted, err := s.IsTokenBlacklisted(tokenString)
	if err != nil {
		// Log error but continue validation (fail open if database is down)
		// In production, you might want to fail closed
	}
	if isBlacklisted {
		return nil, errors.New("token has been revoked")
	}

	claims := &UserClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// GetSecret returns the JWT secret (for testing purposes)
func (s *UserJWTService) GetSecret() string {
	return s.secret
}

// BlacklistToken adds a token to the blacklist in database
func (s *UserJWTService) BlacklistToken(tokenString string, expirationTime time.Duration) error {
	// Parse token để lấy claims (user_id, expires_at)
	claims, err := s.ValidateTokenWithoutBlacklistCheck(tokenString)
	if err != nil {
		// Nếu không parse được, vẫn có thể blacklist bằng hash
		// Nhưng không có user_id và expires_at
		tokenHash := s.hashToken(tokenString)
		blacklistedToken := &model.BlacklistedToken{
			TokenHash: tokenHash,
			UserID:    0, // Unknown user
			ExpiresAt: time.Now().Add(expirationTime),
		}
		return database.Instance().Create(blacklistedToken).Error
	}

	// Hash token để làm key (không lưu plain token)
	tokenHash := s.hashToken(tokenString)
	
	// Tính expires_at từ claims
	expiresAt := claims.ExpiresAt.Time
	if expirationTime > 0 {
		// Nếu có expirationTime, dùng nó
		expiresAt = time.Now().Add(expirationTime)
	}

	// Kiểm tra xem đã blacklist chưa
	var existing model.BlacklistedToken
	err = database.Instance().
		Where("token_hash = ?", tokenHash).
		First(&existing).Error
	
	if err == nil {
		// Đã blacklist rồi
		return nil
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	// Tạo blacklisted token
	blacklistedToken := &model.BlacklistedToken{
		TokenHash: tokenHash,
		UserID:    claims.UserID,
		ExpiresAt: expiresAt,
	}

	return database.Instance().Create(blacklistedToken).Error
}

// ValidateTokenWithoutBlacklistCheck validates token without checking blacklist
// Used internally to parse token for blacklisting
func (s *UserJWTService) ValidateTokenWithoutBlacklistCheck(tokenString string) (*UserClaims, error) {
	claims := &UserClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// IsTokenBlacklisted checks if a token is in the blacklist
func (s *UserJWTService) IsTokenBlacklisted(tokenString string) (bool, error) {
	tokenHash := s.hashToken(tokenString)
	
	var count int64
	err := database.Instance().
		Model(&model.BlacklistedToken{}).
		Where("token_hash = ? AND expires_at > ?", tokenHash, time.Now()).
		Count(&count).Error
	
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// hashToken creates a SHA256 hash of the token
func (s *UserJWTService) hashToken(tokenString string) string {
	hash := sha256.Sum256([]byte(tokenString))
	return hex.EncodeToString(hash[:])
}

// BlacklistTokenByClaims blacklists a token using its claims (for calculating remaining expiration)
func (s *UserJWTService) BlacklistTokenByClaims(tokenString string, claims *UserClaims) error {
	if claims == nil {
		// Nếu không có claims, parse lại
		var err error
		claims, err = s.ValidateToken(tokenString)
		if err != nil {
			return err
		}
	}

	// Tính thời gian còn lại của token
	now := time.Now()
	expiresAt := claims.ExpiresAt.Time
	remainingTime := expiresAt.Sub(now)

	// Nếu token đã hết hạn, không cần blacklist
	if remainingTime <= 0 {
		return nil
	}

	return s.BlacklistToken(tokenString, remainingTime)
}

