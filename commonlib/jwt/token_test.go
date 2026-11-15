package jwt

import (
	"strings"
	"testing"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v4"
)

func TestTokenPair(t *testing.T) {
	m, err := NewManager(Config{
		AccessSecret: strings.Repeat("a", 32), RefreshSecret: strings.Repeat("b", 32),
		AccessTTL: time.Minute, RefreshTTL: time.Hour,
	})
	if err != nil {
		t.Fatal(err)
	}
	sid, err := NewSessionID()
	if err != nil {
		t.Fatal(err)
	}
	const authVersion = "3"
	pair, err := m.GenerateTokenPair(42, sid, "web", authVersion)
	if err != nil {
		t.Fatal(err)
	}
	access, err := m.ParseAccessToken(pair.AccessToken)
	if err != nil {
		t.Fatal(err)
	}
	refresh, err := m.ParseRefreshToken(pair.RefreshToken)
	if err != nil {
		t.Fatal(err)
	}
	for _, claims := range []*Claims{access, refresh} {
		if claims.UserID != 42 || claims.SessionID != sid || claims.Platform != "web" || claims.AuthVersion != authVersion {
			t.Fatalf("unexpected identity: %+v", claims)
		}
	}
	if access.ID == refresh.ID || pair.AccessJTI != access.ID ||
		!pair.AccessExpireAt.Equal(access.ExpiresAt.Time) || !pair.RefreshExpireAt.Equal(refresh.ExpiresAt.Time) ||
		!pair.RefreshExpireAt.After(pair.AccessExpireAt) {
		t.Fatal("unexpected token IDs or expiration")
	}
	if _, err := m.ParseAccessToken(pair.RefreshToken); err == nil {
		t.Error("accepted refresh token as access token")
	}
	if _, err := m.ParseRefreshToken(pair.AccessToken); err == nil {
		t.Error("accepted access token as refresh token")
	}
	if _, err := m.GenerateTokenPair(42, sid, "desktop", authVersion); err == nil {
		t.Error("accepted unsupported platform")
	}
	if _, err := m.GenerateTokenPair(42, sid, "web", ""); err == nil {
		t.Error("accepted empty auth version")
	}

	expired := *access
	expired.IssuedAt = jwtlib.NewNumericDate(time.Now().Add(-2 * time.Hour))
	expired.ExpiresAt = jwtlib.NewNumericDate(time.Now().Add(-time.Hour))
	noExpiry := *access
	noExpiry.ExpiresAt = nil
	noAuthVersion := *access
	noAuthVersion.AuthVersion = ""
	wrongType := *access
	wrongType.TokenType = string(refreshToken)
	for _, tc := range []struct {
		name   string
		alg    jwtlib.SigningMethod
		claims Claims
		key    []byte
	}{
		{"expired", jwtlib.SigningMethodHS256, expired, m.accessKey},
		{"missing expiration", jwtlib.SigningMethodHS256, noExpiry, m.accessKey},
		{"empty auth version", jwtlib.SigningMethodHS256, noAuthVersion, m.accessKey},
		{"wrong type with valid signature", jwtlib.SigningMethodHS256, wrongType, m.accessKey},
		{"wrong signature", jwtlib.SigningMethodHS256, *access, m.refreshKey},
		{"wrong algorithm", jwtlib.SigningMethodHS384, *access, m.accessKey},
	} {
		t.Run(tc.name, func(t *testing.T) {
			raw, err := jwtlib.NewWithClaims(tc.alg, &tc.claims).SignedString(tc.key)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := m.ParseAccessToken(raw); err == nil {
				t.Fatal("accepted invalid token")
			}
		})
	}
}
