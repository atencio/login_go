package security

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secret      []byte
	ttl         time.Duration
	rememberTTL time.Duration
}

func NewJWTManager(secret string, ttl, rememberTTL time.Duration) *JWTManager {
	return &JWTManager{secret: []byte(secret), ttl: ttl, rememberTTL: rememberTTL}
}

type claims struct {
	UserID int64 `json:"uid"`
	jwt.RegisteredClaims
}

func (m *JWTManager) Generate(userID int64, remember bool) (string, time.Time, error) {
	ttl := m.ttl
	if remember {
		ttl = m.rememberTTL
	}
	exp := time.Now().Add(ttl)
	c := claims{UserID: userID, RegisteredClaims: jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(exp),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}}
	tok, err := jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(m.secret)
	return tok, exp, err
}

func (m *JWTManager) Parse(token string) (int64, error) {
	c := &claims{}
	t, err := jwt.ParseWithClaims(token, c, func(t *jwt.Token) (any, error) { return m.secret, nil },
		jwt.WithValidMethods([]string{"HS256"}))
	if err != nil || !t.Valid {
		return 0, errors.New("token inválido")
	}
	return c.UserID, nil
}
