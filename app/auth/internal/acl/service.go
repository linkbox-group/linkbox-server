package acl

import (
	"context"
	"github.com/linkbox-group/linkbox-server/model"
)

type AuthServiceItf interface {
	VerifyTokenService(ctx context.Context, userId string, jti string) (bool, error)
	DestroyTokenService(ctx context.Context, userId string, jti string) error
	GenerateToken(ctx context.Context, userId string, tokenType model.TokenType) (string, error)
	ValidateToken(ctx context.Context, tokenString string, tokenType model.TokenType) (*model.TokenClaims, error)
}
