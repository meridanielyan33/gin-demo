package middleware

import (
	"context"
	errMessage "gin-demo/errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(strategy JWTStrategy) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		} else {
			tokenCookie, err := c.Cookie("token")
			if err == nil {
				tokenString = strings.TrimPrefix(tokenCookie, "Bearer ")
			}
		}

		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errMessage.AuthHeaderMissing})
			return
		}

		claims, err := strategy.ValidateAccessToken(context.Background(), tokenString)
		if err != nil {
			if err.Error() == "token not found in Redis" {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": errMessage.LoggedOut})
				return
			}
			if err.Error() == "invalid token" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errMessage.InvalidToken})
				return
			}
			if err.Error() == "token expired" {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": errMessage.ExpiredToken})
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.Set("email", claims.Email)
		c.Next()
	}
}
