package middleware

import (
	"context"
)

type TokenData struct {
	Email     string
	SessionID string
}

// Token strategy is using strategy pattern,
// if a new type of token will be used,
// it will only implement the interface methods by adding its own
type TokenStrategy interface {
	GenerateToken(ctx context.Context, email string) (string, error)
	ValidateToken(ctx context.Context, token string) (*TokenData, error)
	InvalidateToken(ctx context.Context, email string) error
}
