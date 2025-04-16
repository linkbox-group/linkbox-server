package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/wire"
	"github.com/linkbox-group/linkbox-server/auth/internal/acl"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	tokenKey = "uid:%s:jti:%s"
)

var (
	ProviderSet = wire.NewSet(wire.Bind(new(acl.AuthRepoItf), new(*AuthRepository)), NewTokenRepo)
)

// AuthRepository 实现了 TokenRepo 接口，使用 RedisConfig 存储 token
type AuthRepository struct {
	rdb *redis.Client
}

// NewTokenRepo 创建 TokenRepo 实例
func NewTokenRepo(client *redis.Client) *AuthRepository {
	return &AuthRepository{client}
}

// StoreToken 存储 token
func (r *AuthRepository) StoreToken(ctx context.Context, userId string, jti string, expiration time.Duration) error {
	key := fmt.Sprintf(tokenKey, userId, jti)
	err := r.rdb.Set(ctx, key, true, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to store token: %w", err)
	}
	return nil
}

// VerifyToken 验证 token
func (r *AuthRepository) VerifyToken(ctx context.Context, userId string, jti string) (bool, error) {
	key := fmt.Sprintf(tokenKey, userId, jti)
	isVerified, err := r.rdb.Get(ctx, key).Bool()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil // Token not found
		}
		return false, fmt.Errorf("failed to verify token: %w", err)
	}

	return isVerified, nil
}

// DestroyToken 销毁 token
func (r *AuthRepository) DestroyToken(ctx context.Context, userId string, jti string) error {
	key := fmt.Sprintf(tokenKey, userId, jti)

	err := r.rdb.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to destroy token: %w", err)
	}
	return nil

}
