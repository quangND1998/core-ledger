package jwt

import (
	config "core-ledger/configs"
	"core-ledger/model/core-ledger"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

