package acl

import (
	"context"
	"time"
)

type AuthRepoItf interface {
	StoreToken(ctx context.Context, userId string, jti string, expiration time.Duration) error
	VerifyToken(ctx context.Context, userId string, jti string) (bool, error)
	DestroyToken(ctx context.Context, userId string, jti string) error
}
