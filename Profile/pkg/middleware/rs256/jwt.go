package rs256

import (
	"crypto/rsa"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	UserIDKey = "user_id"
	RoleKey   = "role"
)

type JWTAuthMiddleware struct {
	publicKey *rsa.PublicKey
	issuer    string
	audience  string
}

func New(pubKey *rsa.PublicKey, issuer, audience string) *JWTAuthMiddleware {
	return &JWTAuthMiddleware{
		publicKey: pubKey,
		issuer:    issuer,
		audience:  audience,
	}
}

func (m *JWTAuthMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			unauthorized(c)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			unauthorized(c)
			return
		}

		claims, err := m.verifyToken(parts[1])
		if err != nil {
			unauthorized(c)
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Set(RoleKey, claims.Role)
		c.Next()
	}
}

type claims struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (m *JWTAuthMiddleware) verifyToken(tokenStr string) (*claims, error) {
	if m.issuer == "" || m.audience == "" {
		return nil, errors.New("auth middleware is not properly configured")
	}

	tokenClaims := &claims{}

	token, err := jwt.ParseWithClaims(
		tokenStr,
		tokenClaims,
		func(t *jwt.Token) (any, error) {

			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return m.publicKey, nil
		},
		jwt.WithIssuer(m.issuer),
		jwt.WithAudience(m.audience),
		jwt.WithValidMethods([]string{"RS256"}),
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	
	if tokenClaims.UserID <= 0 {
		return nil, errors.New("invalid user id in token")
	}

	return tokenClaims, nil
}

func unauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(
		http.StatusUnauthorized,
		gin.H{"error": "unauthorized"},
	)
}
