package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/linkbox-group/linkbox-server/gateway/internal/domain"
	"github.com/linkbox-group/linkbox-server/gateway/internal/infra/rpc"
	"github.com/linkbox-group/linkbox-server/rpc-gen/user"
)

var (
	ErrInvalidParamCode = 40001
	ErrRpcFailed        = 50050
)

type UserApi struct {
}

// Login 用户登录
func (api *UserApi) Login(ctx *gin.Context) {
	var req user.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		domain.ErrorMsg(ctx, ErrInvalidParamCode, err.Error())
		return
	}

	resp, err := rpc.UserClient.Login(ctx, &req)
	if err != nil {
		domain.ErrorMsg(ctx, ErrRpcFailed, err.Error())
		return
	}
	if resp.GetError() != nil {
		domain.ErrorMsg(ctx, int(resp.GetError().Code), resp.GetError().Message)
		return
	}
	res := resp.GetAuth()
	vo := domain.UserLoginResp{
		UserData: domain.UserData{
			Id:          res.User.Id,
			Username:    res.User.Username,
			Email:       res.User.Email,
			AvatarUrl:   res.User.AvatarUrl,
			DisplayName: res.User.DisplayName,
			Roles:       res.User.Roles,
		},
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}

	domain.Success(ctx, vo)
}

// Register 用户注册
func (api *UserApi) Register(ctx *gin.Context) {
	var req user.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		domain.ErrorMsg(ctx, ErrInvalidParamCode, err.Error())
		return
	}

	resp, err := rpc.UserClient.Register(ctx, &req)
	if err != nil {
		domain.ErrorMsg(ctx, ErrRpcFailed, err.Error())
		return
	}
	if resp.GetError() != nil {
		domain.ErrorMsg(ctx, int(resp.GetError().Code), resp.GetError().Message)
		return
	}
	res := resp.GetUser()
	vo := domain.UserRegisterResp{
		UserData: domain.UserData{
			Id:          res.Id,
			Username:    res.Username,
			Email:       res.Email,
			AvatarUrl:   res.AvatarUrl,
			DisplayName: res.DisplayName,
			Roles:       res.Roles,
		},
	}
	domain.Success(ctx, vo)

}

// ChangePassword 修改密码
func (api *UserApi) ChangePassword(ctx *gin.Context) {
	var req user.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		domain.ErrorMsg(ctx, ErrInvalidParamCode, err.Error())
		return
	}

	// 设置 userId
	userId, err := domain.GetUserIdFromContext(ctx)
	if err != nil {
		domain.ErrorMsg(ctx, ErrInvalidParamCode, err.Error())
		return
	}
	req.UserId = userId

	resp, err := rpc.UserClient.ChangePassword(ctx, &req)
	if err != nil {
		domain.ErrorMsg(ctx, ErrRpcFailed, err.Error())
		return
	}
	if resp.GetError() != nil {
		domain.ErrorMsg(ctx, int(resp.GetError().Code), resp.GetError().Message)
		return
	}
	res := resp.GetSuccess()
	vo := domain.UserChangePasswordResp{
		Success: res,
	}

	domain.Success(ctx, vo)
}

// UpdateUserInfo 更新用户信息
func (api *UserApi) UpdateUserInfo(ctx *gin.Context) {
	var req user.UpdateUserProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		domain.ErrorMsg(ctx, ErrInvalidParamCode, err.Error())
		return
	}

	// 设置 userId
	userId, err := domain.GetUserIdFromContext(ctx)
	if err != nil {
		domain.ErrorMsg(ctx, ErrInvalidParamCode, err.Error())
		return
	}
	req.UserId = userId

	resp, err := rpc.UserClient.UpdateUserProfile(ctx, &req)
	if err != nil {
		domain.ErrorMsg(ctx, ErrRpcFailed, err.Error())
		return
	}
	if resp.GetError() != nil {
		domain.ErrorMsg(ctx, int(resp.GetError().Code), resp.GetError().Message)
		return
	}
	res := resp.GetProfile()
	vo := domain.UserUpdateResp{
		Id:          res.Id,
		Username:    res.Username,
		Email:       res.Email,
		AvatarUrl:   res.AvatarUrl,
		DisplayName: res.DisplayName,
		Bio:         res.Bio,
		CreatedAt:   res.CreatedAt.AsTime(),
		UpdatedAt:   res.UpdatedAt.AsTime(),
	}
	domain.Success(ctx, vo)
}

// GetUserInfo 获取用户信息
func (api *UserApi) GetUserInfo(ctx *gin.Context) {
	var req user.UserIdRequest

	// 设置 userId
	userId, err := domain.GetUserIdFromContext(ctx)

	req.UserId = userId
	fmt.Println("userID", req.UserId)
	if err != nil {
		domain.ErrorMsg(ctx, ErrInvalidParamCode, err.Error())
		return
	}

	resp, err := rpc.UserClient.GetUserProfile(ctx, &req)
	if err != nil {
		domain.ErrorMsg(ctx, ErrRpcFailed, err.Error())
		return
	}
	if resp.GetError() != nil {
		domain.ErrorMsg(ctx, int(resp.GetError().Code), resp.GetError().Message)
		return
	}
	res := resp.GetProfile()
	vo := domain.UserGetInfoResp{
		Id:          res.Id,
		Username:    res.Username,
		Email:       res.Email,
		AvatarUrl:   res.AvatarUrl,
		DisplayName: res.DisplayName,
		Bio:         res.Bio,
		CreatedAt:   res.CreatedAt.AsTime(),
		UpdatedAt:   res.UpdatedAt.AsTime(),
	}

	domain.Success(ctx, vo)
}

// RefreshToken 刷新访问令牌
func (api *UserApi) RefreshToken(ctx *gin.Context) {
	var req user.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		domain.ErrorMsg(ctx, ErrInvalidParamCode, err.Error())
		return
	}

	resp, err := rpc.UserClient.RefreshToken(ctx, &req)
	if err != nil {
		domain.ErrorMsg(ctx, ErrRpcFailed, err.Error())
		return
	}
	if resp.GetError() != nil {
		domain.ErrorMsg(ctx, int(resp.GetError().Code), resp.GetError().Message)
		return
	}

	vo := domain.TokenRefreshResp{
		AccessToken: resp.GetAuth().AccessToken,
		Success:     true,
	}

	domain.Success(ctx, vo)
}

func (api *UserApi) DeleteUser(ctx *gin.Context) {
	var req user.DeleteAccountRequest

	// 设置 userId
	userId, err := domain.GetUserIdFromContext(ctx)
	if err != nil {
		domain.ErrorMsg(ctx, ErrInvalidParamCode, err.Error())
		return
	}
	req.UserId = userId

	resp, err := rpc.UserClient.DeleteAccount(ctx, &req)
	if err != nil {
		domain.ErrorMsg(ctx, ErrRpcFailed, err.Error())
		return
	}
	if resp.GetError() != nil {
		domain.ErrorMsg(ctx, int(resp.GetError().Code), resp.GetError().Message)
		return
	}
	vo := domain.UserLogoutResp{
		Success: resp.GetSuccess(),
	}

	domain.Success(ctx, vo)
}
