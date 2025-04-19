package acl

import (
	"context"
	"github.com/linkbox-group/linkbox-server/model"
)

type UserServiceItf interface {
	CreateItem(ctx context.Context, req *model.Item) (err error)
	GetItem(ctx context.Context, req *model.Item) (err error)
}
