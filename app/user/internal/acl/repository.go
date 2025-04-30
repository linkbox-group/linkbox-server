package acl

import (
	"context"
	"github.com/linkbox-group/linkbox-server/model"
)

type UserRepositoryItf interface {
	CreateUser(ctx context.Context, user *model.User) (err error)
	FindUserByEmail(ctx context.Context, email string) (user *model.User, err error)
	FindUserByID(ctx context.Context, id string) (user *model.User, err error)
	UpdateUser(ctx context.Context, user *model.User) (err error)
	DeleteUser(ctx context.Context, id string) (err error)
	RegisterUser(ctx context.Context, user *model.User) (err error)
}
