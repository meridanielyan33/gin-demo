package middleware

import (
	"context"
	"errors"
	"gin-demo/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Claims struct {
	Email     string `json:"email"`
	SessionID string `json:"session_id"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(email string) (string, error) {
	cfg := config.GetConfig()

	secretKey := cfg.Secret
	expirationTime := time.Now().Add(8 * time.Hour)
	sessionID := uuid.New().String()

	claims := &Claims{
		Email:     email,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func ValidateAccessToken(tokenString string, redisClient *redis.Client) (*Claims, error) {
	cfg := config.GetConfig()

	secretKey := cfg.Secret
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	}, jwt.WithoutClaimsValidation())
	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	ctx := context.Background()
	redisKey := claims.Email
	storedToken, err := redisClient.Get(ctx, redisKey).Result()
	if errors.Is(err, redis.Nil) {
		return nil, errors.New("token not found in Redis")
	} else if err != nil {
		return nil, err
	}

	if storedToken != tokenString {
		return nil, errors.New("token mismatch")
	}

	if time.Now().After(claims.ExpiresAt.Time) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}
