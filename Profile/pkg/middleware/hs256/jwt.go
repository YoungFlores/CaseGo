package hs256

import (
	"fmt"

	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	UserIDKey = "user_id"
	RoleKey   = "role"
)

type JWTAuthMiddleware struct {
	secret string
}

func New(secret string) *JWTAuthMiddleware {
	return &JWTAuthMiddleware{secret: secret}
}

func (m *JWTAuthMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization header required"})
			return
		}

		tokenStr, ok := extractBearerToken(auth)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid Authorization format, expected Bearer <token>"})
			return
		}

		claims, err := m.parseAndValidateToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Set(RoleKey, claims.Role)
		c.Next()
	}
}

// extractBearerToken возвращает токен и true, если заголовок вида "Bearer xxx"
func extractBearerToken(authHeader string) (string, bool) {
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 {
		return "", false
	}
	if !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}
	return parts[1], true
}

type CustomClaims struct {
	UserID int `json:"user_id"`
	Role   int `json:"role"`
	jwt.RegisteredClaims
}

func (m *JWTAuthMiddleware) parseAndValidateToken(tokenStr string) (*CustomClaims, error) {
	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(m.secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("token parse error: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid or expired token")
	}

	// Проверяем, что токен не из будущего (iat)
	if claims.IssuedAt != nil && claims.IssuedAt.After(time.Now()) {
		return nil, fmt.Errorf("token issued in future")
	}

	return claims, nil
}