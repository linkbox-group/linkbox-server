package service

import (
	"context"
	"fmt"
	"github.com/google/wire"
	"github.com/spf13/viper"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/linkbox-group/linkbox-server/model"

	"github.com/linkbox-group/linkbox-server/auth/internal/acl"
)

const tokenKeyFormat = "uid:%s:jti:%s"

var (
	ProviderSet = wire.NewSet(wire.Bind(new(acl.AuthServiceItf), new(*AuthService)), NewAuthService)
)

type AuthService struct {
	repo acl.AuthRepoItf
}

// GenerateTokenKey 根据 userId 生成 Redis 的 token key
func GenerateTokenKey(userId string) string {
	return fmt.Sprintf(tokenKeyFormat, userId)
}

func NewAuthService(repo acl.AuthRepoItf) *AuthService {
	return &AuthService{repo: repo}

}

func (s *AuthService) VerifyTokenService(ctx context.Context, userId string, jti string) (bool, error) {
	isValid, err := s.repo.VerifyToken(ctx, userId, jti)
	if err != nil {
		return false, fmt.Errorf("verify token failed: %w", err)
	}
	return isValid, nil

}

func (s *AuthService) DestroyTokenService(ctx context.Context, userId string, jti string) error {
	err := s.repo.DestroyToken(ctx, userId, jti)
	if err != nil {
		return fmt.Errorf("destroy token failed: %w", err)
	}
	return nil

}

func (s *AuthService) GenerateToken(ctx context.Context, userId string, tokenType model.TokenType) (tokenString string, err error) {
	//  生成 JTI
	jti := uuid.New().String()
	expiredAt := viper.GetDuration("jwt.access_duration") * time.Second
	if tokenType == model.TokenTypeRefresh {
		expiredAt = viper.GetDuration("jwt.refresh_duration") * time.Second
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, model.TokenClaims{
		UserId: userId,
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    viper.GetString("jwt.Issuer"),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiredAt)),
			ID:        jti,
		},
	})
	if tokenType == model.TokenTypeRefresh {
		err = s.repo.StoreToken(ctx, userId, jti, expiredAt)
	}
	if err != nil {
		return "", fmt.Errorf("store token failed: %w", err)
	}
	return token.SignedString([]byte(viper.GetString("JWT_SECRET")))
}

// ValidateToken 验证 token
func (s *AuthService) ValidateToken(ctx context.Context, tokenString string, tokenType model.TokenType) (*model.TokenClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &model.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(viper.GetString("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, fmt.Errorf("parse token failed: %w", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token content")
	}
	//查看claim
	fmt.Println(token.Claims)

	claims, ok := token.Claims.(*model.TokenClaims)

	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	if claims.Type != tokenType {
		return nil, fmt.Errorf("invalid token type")
	}

	return claims, nil
}
