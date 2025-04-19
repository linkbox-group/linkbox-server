package acl

import (
	"context"
	"github.com/linkbox-group/linkbox-server/model"
)

type UserRepositoryItf interface {
	CreateItem(ctx context.Context, req *model.Item) (err error)
}
