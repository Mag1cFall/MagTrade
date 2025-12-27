// JWT Token 工具
//
// 本檔案處理 JWT Token 生成和驗證
// 支援 Access Token 和 Refresh Token 雙 Token 機制
// 使用 HS256 簽名演算法
package utils

import (
	"errors"
	"time"

	"github.com/Mag1cFall/magtrade/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

// TokenType Token 類型
type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

// Claims 自訂 JWT Claims
type Claims struct {
	UserID   int64     `json:"user_id"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
	Type     TokenType `json:"type"` // access/refresh
	jwt.RegisteredClaims
}

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// GenerateTokenPair 生成 Access + Refresh Token 對
func GenerateTokenPair(userID int64, username, role string, cfg *config.JWTConfig) (accessToken, refreshToken string, err error) {
	accessToken, err = generateToken(userID, username, role, AccessToken, cfg.AccessTokenExpire, cfg.Secret)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = generateToken(userID, username, role, RefreshToken, cfg.RefreshTokenExpire, cfg.Secret)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// generateToken 生成單個 Token
func generateToken(userID int64, username, role string, tokenType TokenType, expire time.Duration, secret string) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		Type:     tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "magtrade",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken 解析並驗證 Token
func ParseToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// ValidateAccessToken 驗證 Access Token
func ValidateAccessToken(tokenString, secret string) (*Claims, error) {
	claims, err := ParseToken(tokenString, secret)
	if err != nil {
		return nil, err
	}

	if claims.Type != AccessToken {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// ValidateRefreshToken 驗證 Refresh Token
func ValidateRefreshToken(tokenString, secret string) (*Claims, error) {
	claims, err := ParseToken(tokenString, secret)
	if err != nil {
		return nil, err
	}

	if claims.Type != RefreshToken {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
