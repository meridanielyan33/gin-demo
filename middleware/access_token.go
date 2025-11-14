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

type JWTStrategy struct {
	redis  *redis.Client
	secret string
	ttl    time.Duration
}

func NewJWTStrategy(redisClient *redis.Client) *JWTStrategy {
	cfg := config.GetConfig()
	return &JWTStrategy{
		redis:  redisClient,
		secret: cfg.Secret,
		ttl:    8 * time.Hour,
	}
}

func (s *JWTStrategy) GenerateAccessToken(ctx context.Context, email string) (string, error) {
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
	signed, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", err
	}

	if err := s.redis.Set(ctx, email, signed, s.ttl).Err(); err != nil {
		return "", err
	}
	return signed, nil
}

func (s *JWTStrategy) ValidateAccessToken(ctx context.Context, tokenString string) (*TokenData, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	}, jwt.WithoutClaimsValidation())
	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	storedToken, err := s.redis.Get(ctx, claims.Email).Result()
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

	return &TokenData{
		Email:     claims.Email,
		SessionID: claims.SessionID,
	}, nil
}

func (s *JWTStrategy) InvalidateToken(ctx context.Context, email string) error {
	return s.redis.Del(ctx, email).Err()
}
