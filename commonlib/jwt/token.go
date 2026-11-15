package jwt

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v4"
)

type tokenType string

const (
	accessToken  tokenType = "access"
	refreshToken tokenType = "refresh"
)

type Claims struct {
	UserID      int64  `json:"uid"`
	SessionID   string `json:"sid"`
	Platform    string `json:"platform"`
	TokenType   string `json:"tokenType"`
	AuthVersion string `json:"authVersion"`
	jwtlib.RegisteredClaims
}

type TokenPair struct {
	AccessToken     string
	RefreshToken    string
	AccessExpireAt  time.Time
	RefreshExpireAt time.Time
	AccessJTI       string
}

type Config struct {
	AccessSecret  string
	RefreshSecret string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

type Manager struct {
	accessKey  []byte
	refreshKey []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewManager(cfg Config) (*Manager, error) {
	if len(cfg.AccessSecret) < 32 || len(cfg.RefreshSecret) < 32 {
		return nil, errors.New("JWT HS256 secrets must each contain at least 32 bytes")
	}
	if cfg.AccessSecret == cfg.RefreshSecret {
		return nil, errors.New("JWT access and refresh secrets must differ")
	}
	if cfg.AccessTTL < time.Second || cfg.RefreshTTL/time.Second <= cfg.AccessTTL/time.Second {
		return nil, errors.New("JWT refresh TTL must exceed access TTL at second precision")
	}
	return &Manager{
		accessKey:  []byte(cfg.AccessSecret),
		refreshKey: []byte(cfg.RefreshSecret),
		accessTTL:  cfg.AccessTTL,
		refreshTTL: cfg.RefreshTTL,
	}, nil
}

// NewSessionID creates the ID for a new login. Refresh reuses the existing ID.
func NewSessionID() (string, error) {
	return randomID()
}

func randomID() (string, error) {
	var id [16]byte
	if _, err := rand.Read(id[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(id[:]), nil
}

func (m *Manager) GenerateTokenPair(userID int64, sessionID, platform string, authVersion string) (TokenPair, error) {
	identity := Claims{
		UserID:      userID,
		SessionID:   sessionID,
		Platform:    platform,
		AuthVersion: authVersion,
	}
	if !identity.validIdentity() {
		return TokenPair{}, errors.New("invalid JWT session identity")
	}
	now := time.Now().UTC().Truncate(time.Second)
	access, accessClaims, err := sign(identity, accessToken, m.accessKey, now, m.accessTTL)
	if err != nil {
		return TokenPair{}, err
	}
	refresh, refreshClaims, err := sign(identity, refreshToken, m.refreshKey, now, m.refreshTTL)
	if err != nil {
		return TokenPair{}, err
	}
	return TokenPair{
		AccessToken:     access,
		RefreshToken:    refresh,
		AccessExpireAt:  accessClaims.ExpiresAt.Time,
		RefreshExpireAt: refreshClaims.ExpiresAt.Time,
		AccessJTI:       accessClaims.ID,
	}, nil
}

func sign(claims Claims, kind tokenType, secret []byte, now time.Time, ttl time.Duration) (string, Claims, error) {
	jti, err := randomID()
	if err != nil {
		return "", Claims{}, err
	}
	claims.TokenType = string(kind)
	claims.RegisteredClaims = jwtlib.RegisteredClaims{
		ID:        jti,
		IssuedAt:  jwtlib.NewNumericDate(now),
		ExpiresAt: jwtlib.NewNumericDate(now.Add(ttl)),
	}
	raw, err := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, &claims).SignedString(secret)
	return raw, claims, err
}

func (m *Manager) ParseAccessToken(raw string) (*Claims, error) {
	return parse(raw, accessToken, m.accessKey)
}

func (m *Manager) ParseRefreshToken(raw string) (*Claims, error) {
	return parse(raw, refreshToken, m.refreshKey)
}

func parse(raw string, kind tokenType, secret []byte) (*Claims, error) {
	var claims Claims
	token, err := jwtlib.ParseWithClaims(raw, &claims, func(_ *jwtlib.Token) (any, error) {
		return secret, nil
	}, jwtlib.WithValidMethods([]string{jwtlib.SigningMethodHS256.Alg()}))
	if err != nil || !token.Valid || !claims.validIdentity() || claims.TokenType != string(kind) ||
		claims.ID == "" || claims.IssuedAt == nil || claims.ExpiresAt == nil ||
		!claims.ExpiresAt.After(claims.IssuedAt.Time) {
		return nil, errors.New("invalid JWT")
	}
	return &claims, nil
}

func (c Claims) validIdentity() bool {
	switch c.Platform {
	case "web", "android", "ios":
		return c.UserID > 0 && c.SessionID != "" && c.AuthVersion != ""
	default:
		return false
	}
}

func TokenHash(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
