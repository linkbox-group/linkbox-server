package delivery

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/wire"
	"github.com/linkbox-group/linkbox-server/rpc-gen/user"
	"github.com/linkbox-group/linkbox-server/user/internal/acl"
	"github.com/linkbox-group/linkbox-server/user/pkg/log"
	"github.com/sirupsen/logrus"
)

var (
	errPasswordNotMatch = errors.New("password does not match")
)

var ProviderSet = wire.NewSet(NewUserDelivery)

type UserDelivery struct {
	service acl.UserServiceItf
}

func NewUserDelivery(service acl.UserServiceItf) *UserDelivery {
	return &UserDelivery{
		service: service,
	}
}

// SendCode implements the UserServiceImpl interface.
func (d *UserDelivery) SendCode(ctx context.Context, req *user.SendCodeReq) (resp *user.SendCodeResp, err error) {
	err = d.service.SendCode(ctx, req.Email)
	if err != nil {
		return &user.SendCodeResp{
			Success: false,
			Message: err.Error(),
		}, err
	}

	return &user.SendCodeResp{
		Success: true,
		Message: "send code success",
	}, nil
}

// Register implements the UserServiceImpl interface.
func (d *UserDelivery) Register(ctx context.Context, req *user.RegisterReq) (resp *user.RegisterResp, err error) {

	registerUser, err := d.service.RegisterUser(ctx, req.Email, req.Code, req.Password)
	if err != nil {
		log.Log().Error(err.Error(), "req", req)
		return nil, err
	}
	return registerUser, nil
}

// Login implements the UserServiceImpl interface.
func (d *UserDelivery) Login(ctx context.Context, req *user.LoginReq) (resp *user.LoginResp, err error) {
	loginUser, err := d.service.LoginUser(ctx, req.Email, req.Password)
	if err != nil {
		log.Log().Error(err.Error(), "req", req)
		return nil, fmt.Errorf("login:%w", err)
	}

	return loginUser, err
}

// UpdateUserInfo implements the UserServiceImpl interface.
func (d *UserDelivery) UpdateUserInfo(ctx context.Context, req *user.UpdateUserInfoReq) (res *user.UpdateUserInfoResp, err error) {
	err = d.service.UpdateUserInfo(ctx, req)
	if err != nil {
		logrus.Errorln(err)
		return &user.UpdateUserInfoResp{
			Success: false,
			Message: err.Error(),
		}, fmt.Errorf("update user info:%w", err)
	}

	return &user.UpdateUserInfoResp{
		Success: true,
		Message: "update user info success",
	}, nil
}

// GetUserInfo implements the UserServiceImpl interface.
func (d *UserDelivery) GetUserInfo(ctx context.Context, req *user.GetUserInfoReq) (res *user.GetUserInfoResp, err error) {
	resp, err := d.service.GetUser(ctx, req.UserId)
	if err != nil {
		logrus.Errorln(err)
		return nil, fmt.Errorf("get user info:%w", err)
	}
	return resp, nil
}

// DeleteUser implements the UserServiceImpl interface.
func (d *UserDelivery) DeleteUser(ctx context.Context, req *user.DeleteUserReq) (res *user.DeleteUserResp, err error) {
	err = d.service.DeleteUser(ctx, req.UserId)
	if err != nil {
		logrus.Errorln(err)
		return &user.DeleteUserResp{
			Success: false,
			Message: err.Error(),
		}, fmt.Errorf("delete user:%w", err)
	}

	return &user.DeleteUserResp{
		Success: true,
		Message: "delete user success",
	}, nil
}

// ChangePassword implements the UserServiceImpl interface.
func (d *UserDelivery) ChangePassword(ctx context.Context, req *user.ChangePasswordReq) (res *user.ChangePasswordResp, err error) {
	err = d.service.UpdatePassword(ctx, req.UserId, req.OldPassword, req.NewPassword)
	if err != nil {
		logrus.Errorln(err)
		return &user.ChangePasswordResp{
			Success: false,
			Message: err.Error(),
		}, fmt.Errorf("change password:%w", err)
	}

	return &user.ChangePasswordResp{
		Success: true,
		Message: "change password success",
	}, nil
}
