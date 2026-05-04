package jwt

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type claimsKey struct{}

// JWTAuth возвращает middleware для проверки Bearer-токена
func JWTAuth(tm TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing_authorization"})
			return
		}

		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid_authorization_format"})
			return
		}

		claims, err := tm.ValidateToken(c.Request.Context(), parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid_or_expired_token"})
			return
		}

		// запись в контекст
		ctx := context.WithValue(c.Request.Context(), claimsKey{}, claims)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func GetClaims(c *gin.Context) *AuthClaims {
	if v, ok := c.Request.Context().Value(claimsKey{}).(*AuthClaims); ok {
		return v
	}
	return nil
}
