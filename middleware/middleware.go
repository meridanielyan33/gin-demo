package middleware

import (
	"context"
	errMessage "gin-demo/errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func AuthMiddleware(redisClient *redis.Client) gin.HandlerFunc {
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

		claims, err := ValidateAccessToken(tokenString, redisClient)
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

		c.Set("claims", claims)
		c.Next()
	}
}

func InvalidateToken(c *gin.Context, email string, redisClient *redis.Client) {
	val, err := redisClient.Get(context.Background(), email).Result()
	if err != nil && err != redis.Nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMessage.RedisTokenFail})
		return
	}

	if val != "" {
		err := redisClient.Del(context.Background(), email).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errMessage.TokenDeleteFailRedis})
			return
		}
		log.Printf("User with email %v logged out successfully. Session invalidated.", email)
	} else {
		log.Printf("No active session found for user with email %v.", email)
	}
}
