package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/linkbox-group/linkbox-server/common/ecode"
	"github.com/linkbox-group/linkbox-server/gateway/internal/domain"
	"github.com/linkbox-group/linkbox-server/gateway/internal/infra/rpc"
	"github.com/linkbox-group/linkbox-server/rpc-gen/auth"
	"github.com/linkbox-group/linkbox-server/rpc-gen/user"
)

type UserApi struct {
}

// SendCode 发送验证码
func (api *UserApi) SendCode(ctx *gin.Context) {
	var req user.SendCodeReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fmt.Println(err)
		domain.ErrorMsg(ctx, ecode.ErrInvalidParam, err.Error())
		return
	}

	resp, err := rpc.UserClient.SendCode(ctx, &req)
	if err != nil {
		fmt.Println(err)
		domain.ErrorMsg(ctx, ecode.ErrRpcServiceError, err.Error())
		return
	}

	domain.Success(ctx, resp)
}

// Login 用户登录
func (api *UserApi) Login(ctx *gin.Context) {
	var req user.LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		domain.ErrorMsg(ctx, ecode.ErrInvalidParam, err.Error())
	}

	resp, err := rpc.UserClient.Login(ctx, &req)
	if err != nil {
		domain.ErrorMsg(ctx, ecode.ErrRpcServiceError, err.Error())
		return
	}
	vo := domain.UserLoginResp{
		UserId:       resp.UserId,
		Username:     resp.Username,
		Email:        resp.Email,
		Avatar:       resp.Avatar,
		Bio:          resp.Bio,
		Theme:        resp.Theme,
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}

	domain.Success(ctx, vo)
}

// Register 用户注册
func (api *UserApi) Register(ctx *gin.Context) {
	var req user.RegisterReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		domain.ErrorMsg(ctx, ecode.ErrInvalidParam, err.Error())
	}

	resp, err := rpc.UserClient.Register(ctx, &req)
	if err != nil {
		domain.ErrorMsg(ctx, ecode.ErrRpcServiceError, err.Error())
		return
	}
	vo := domain.UserRegisterResp{
		UserId:       resp.UserId,
		Username:     resp.Username,
		Email:        resp.Email,
		Avatar:       resp.Avatar,
		Bio:          resp.Bio,
		Theme:        resp.Theme,
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}

	domain.Success(ctx, vo)
}

// ChangePassword 修改密码
func (api *UserApi) ChangePassword(ctx *gin.Context) {
	var req user.ChangePasswordReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		domain.ErrorMsg(ctx, ecode.ErrInvalidParam, err.Error())
		return
	}

	// 设置 userId
	userId, err := domain.GetUserIdFromContext(ctx)
	if err != nil {
		domain.ErrorMsg(ctx, ecode.ErrAuthFailed, err.Error())
		return
	}
	req.UserId = userId

	resp, err := rpc.UserClient.ChangePassword(ctx, &req)
	if err != nil {
		domain.ErrorMsg(ctx, ecode.ErrRpcServiceError, err.Error())
		return
	}
	vo := domain.UserChangePasswordResp{
		Success: resp.Success,
		Message: resp.Message,
	}

	domain.Success(ctx, vo)
}

// UpdateUserInfo 更新用户信息
func (api *UserApi) UpdateUserInfo(ctx *gin.Context) {
	var req user.UpdateUserInfoReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		domain.ErrorMsg(ctx, ecode.ErrInvalidParam, err.Error())
		return
	}

	// 设置 userId
	userId, err := domain.GetUserIdFromContext(ctx)
	if err != nil {
		domain.ErrorMsg(ctx, ecode.ErrAuthFailed, err.Error())
		return
	}
	req.UserId = userId

	resp, err := rpc.UserClient.UpdateUserInfo(ctx, &req)
	if err != nil {
		domain.ErrorMsg(ctx, ecode.ErrRpcServiceError, err.Error())
		return
	}
	vo := domain.UserUpdateResp{
		Success: resp.Success,
		Message: resp.Message,
	}
	domain.Success(ctx, vo)
}

// GetUserInfo 获取用户信息
func (api *UserApi) GetUserInfo(ctx *gin.Context) {
	var req user.GetUserInfoReq

	// 设置 userId
	userId, err := domain.GetUserIdFromContext(ctx)

	req.UserId = userId
	fmt.Println("userID", req.UserId)
	if err != nil {
		domain.ErrorMsg(ctx, ecode.ErrAuthFailed, err.Error())
		return
	}

	resp, err := rpc.UserClient.GetUserInfo(ctx, &req)
	if err != nil {
		domain.ErrorMsg(ctx, ecode.ErrRpcServiceError, err.Error())
		return
	}
	vo := domain.UserGetInfoResp{
		UserId:   resp.UserId,
		Username: resp.Username,
		Email:    resp.Email,
		Avatar:   resp.Avatar,
		Bio:      resp.Bio,
		Theme:    resp.Theme,
	}

	domain.Success(ctx, vo)
}

// DeleteUser 删除用户
func (api *UserApi) DeleteUser(ctx *gin.Context) {
	var req user.DeleteUserReq

	// 设置 userId
	userId, err := domain.GetUserIdFromContext(ctx)
	if err != nil {
		domain.ErrorMsg(ctx, ecode.ErrAuthFailed, err.Error())
		return
	}
	req.UserId = userId

	resp, err := rpc.UserClient.DeleteUser(ctx, &req)
	if err != nil {
		domain.ErrorMsg(ctx, ecode.ErrRpcServiceError, err.Error())
		return
	}
	vo := domain.UserLogoutResp{
		Success: resp.Success,
		Message: resp.Message,
	}

	domain.Success(ctx, vo)
}

// RefreshToken 刷新访问令牌
func (api *UserApi) RefreshToken(ctx *gin.Context) {
	var req auth.TokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		domain.ErrorMsg(ctx, ecode.ErrInvalidParam, err.Error())
		return
	}

	resp, err := rpc.AuthClient.RefreshToken(ctx, &req)
	if err != nil {
		domain.ErrorMsg(ctx, ecode.ErrRpcServiceError, err.Error())
		return
	}

	domain.Success(ctx, resp.Token)
}
