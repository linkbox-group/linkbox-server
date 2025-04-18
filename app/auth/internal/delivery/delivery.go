package delivery

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/wire"
	"github.com/linkbox-group/linkbox-server/auth/internal/acl"
	"github.com/linkbox-group/linkbox-server/model"
	"github.com/linkbox-group/linkbox-server/rpc-gen/auth"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrInvalidToken   = errors.New("invalid token")
	ErrInternalError  = errors.New("internal error")
)
var (
	ProviderSet = wire.NewSet(NewAuthDelivery)
)

type AuthDelivery struct {
	service acl.AuthServiceItf
}

func NewAuthDelivery(service acl.AuthServiceItf) *AuthDelivery {
	return &AuthDelivery{
		service: service,
	}
}

// RefreshToken implements the TokenImpl interface.
func (d *AuthDelivery) RefreshToken(ctx context.Context, req *auth.TokenRequest) (resp *auth.TokenString, err error) {
	// 检验参数
	if req == nil || req.Token == "" {
		return nil, fmt.Errorf("%w:%w", ErrInvalidRequest, err)
	}

	// 解析 refresh token
	claim, err := d.service.ValidateToken(ctx, req.Token, model.TokenTypeRefresh)
	if err != nil {
		return nil, fmt.Errorf("%w:%w", ErrInvalidToken, err)
	}

	// 查询 refresh token

	valid, err := d.service.VerifyTokenService(ctx, claim.UserId, claim.ID)
	if err != nil {
		return nil, fmt.Errorf("%w:%w", ErrInternalError, err)
	}
	if !valid {
		return nil, fmt.Errorf("%w:%w", ErrInvalidToken, err)
	}

	// 生成新的 access token
	accessToken, err := d.service.GenerateToken(ctx, claim.UserId, model.TokenTypeAccess)
	if err != nil {
		return nil, fmt.Errorf("%w:%w", ErrInternalError, err)
	}

	return &auth.TokenString{
		Token: accessToken,
	}, nil
}

// DestroyRefreshToken implements the TokenImpl interface.
func (d *AuthDelivery) DestroyRefreshToken(ctx context.Context, req *auth.TokenRequest) (resp *emptypb.Empty, err error) {
	// 检验参数
	if req == nil || req.Token == "" {
		return nil, fmt.Errorf("%w:%w", ErrInvalidRequest, err)
	}

	// 解析 refresh token
	claim, err := d.service.ValidateToken(ctx, req.Token, model.TokenTypeRefresh)
	if err != nil {
		return nil, fmt.Errorf("%w:%w", ErrInvalidToken, err)
	}

	// 销毁 refresh token
	err = d.service.DestroyTokenService(ctx, claim.UserId, claim.ID)
	if err != nil {
		return nil, fmt.Errorf("%w:%w", ErrInternalError, err)
	}

	return &emptypb.Empty{}, nil
}

// VerifyAccessToken implements the TokenImpl interface.
func (d *AuthDelivery) VerifyAccessToken(ctx context.Context, req *auth.TokenRequest) (resp *auth.AccessTokenInfo, err error) {
	// 检验参数
	if req == nil || req.Token == "" {
		return nil, fmt.Errorf("%w:%w", ErrInvalidRequest, err)
	}

	// 解析 access token
	claim, err := d.service.ValidateToken(ctx, req.Token, model.TokenTypeAccess)
	if err != nil {
		return nil, fmt.Errorf("%w:%w", ErrInvalidToken, err)
	}

	return &auth.AccessTokenInfo{
		Uid: claim.UserId,
	}, nil
}

// GenerateRefreshToken implements the AuthServiceImpl interface.
func (d *AuthDelivery) GenerateRefreshToken(ctx context.Context, req *auth.GenerateTokenReq) (resp *auth.TokenString, err error) {
	token, err := d.service.GenerateToken(ctx, req.Uid, model.TokenTypeRefresh)
	if err != nil {
		return nil, fmt.Errorf("%w:%w", ErrInternalError, err)
	}
	resp = &auth.TokenString{
		Token: token,
	}
	return resp, nil
}

// GenerateAccessToken implements the AuthServiceImpl interface.
func (d *AuthDelivery) GenerateAccessToken(ctx context.Context, req *auth.GenerateTokenReq) (resp *auth.TokenString, err error) {
	token, err := d.service.GenerateToken(ctx, req.Uid, model.TokenTypeAccess)
	if err != nil {
		return nil, fmt.Errorf("%w:%w", ErrInternalError, err)
	}
	resp = &auth.TokenString{
		Token: token,
	}
	return resp, nil
}
