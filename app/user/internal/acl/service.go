package acl

import (
	"context"
	"github.com/linkbox-group/linkbox-server/rpc-gen/user"
)

type UserServiceItf interface {
	SendCode(ctx context.Context, email string) (err error)
	RegisterUser(ctx context.Context, email, code, password string) (resp *user.RegisterResp, err error)
	LoginUser(ctx context.Context, email, password string) (resp *user.LoginResp, err error)
	GetUser(ctx context.Context, id string) (resp *user.GetUserInfoResp, err error)
	UpdatePassword(ctx context.Context, id string, oldPassword, newPassword string) (err error)
	UpdateUserInfo(ctx context.Context, req *user.UpdateUserInfoReq) (err error)
	DeleteUser(ctx context.Context, id string) (err error)
}
