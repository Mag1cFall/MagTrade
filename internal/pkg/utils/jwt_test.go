// JWT Token 工具單元測試
//
// 測試覆蓋：
// - GenerateTokenPair: Access Token 和 Refresh Token 生成
// - ValidateAccessToken: Access Token 驗證（有效、無效、過期、錯誤密鑰）
// - ValidateRefreshToken: Refresh Token 驗證
// - Claims 結構體欄位驗證（UserID, Username, Role）
package utils

import (
	"testing"
	"time"

	"github.com/Mag1cFall/magtrade/internal/config"
)

func TestGenerateAndValidateTokens(t *testing.T) {
	cfg := &config.JWTConfig{
		Secret:             "test-secret-key-12345678",
		AccessTokenExpire:  2 * time.Hour,
		RefreshTokenExpire: 168 * time.Hour,
	}

	accessToken, refreshToken, err := GenerateTokenPair(123, "testuser", "admin", cfg)
	if err != nil {
		t.Fatalf("GenerateTokenPair() error = %v", err)
	}

	if accessToken == "" {
		t.Error("GenerateTokenPair() returned empty access token")
	}

	if refreshToken == "" {
		t.Error("GenerateTokenPair() returned empty refresh token")
	}

	claims, err := ValidateAccessToken(accessToken, cfg.Secret)
	if err != nil {
		t.Fatalf("ValidateAccessToken() error = %v", err)
	}

	if claims.UserID != 123 {
		t.Errorf("claims.UserID = %v, want %v", claims.UserID, 123)
	}

	if claims.Username != "testuser" {
		t.Errorf("claims.Username = %v, want %v", claims.Username, "testuser")
	}

	if claims.Role != "admin" {
		t.Errorf("claims.Role = %v, want %v", claims.Role, "admin")
	}

	refreshClaims, err := ValidateRefreshToken(refreshToken, cfg.Secret)
	if err != nil {
		t.Fatalf("ValidateRefreshToken() error = %v", err)
	}

	if refreshClaims.UserID != 123 {
		t.Errorf("refreshClaims.UserID = %v, want %v", refreshClaims.UserID, 123)
	}
}

func TestValidateAccessToken_InvalidToken(t *testing.T) {
	_, err := ValidateAccessToken("invalid.token.here", "test-secret")
	if err == nil {
		t.Error("ValidateAccessToken() should return error for invalid token")
	}
}

func TestValidateAccessToken_WrongSecret(t *testing.T) {
	cfg := &config.JWTConfig{
		Secret:             "correct-secret",
		AccessTokenExpire:  2 * time.Hour,
		RefreshTokenExpire: 168 * time.Hour,
	}

	accessToken, _, _ := GenerateTokenPair(1, "testuser", "user", cfg)

	_, err := ValidateAccessToken(accessToken, "wrong-secret")
	if err == nil {
		t.Error("ValidateAccessToken() should return error for wrong secret")
	}
}

func TestValidateAccessToken_ExpiredToken(t *testing.T) {
	cfg := &config.JWTConfig{
		Secret:             "test-secret",
		AccessTokenExpire:  -1 * time.Hour,
		RefreshTokenExpire: 168 * time.Hour,
	}

	accessToken, _, _ := GenerateTokenPair(1, "testuser", "user", cfg)

	_, err := ValidateAccessToken(accessToken, cfg.Secret)
	if err != ErrExpiredToken {
		t.Errorf("Expected ErrExpiredToken, got %v", err)
	}
}
