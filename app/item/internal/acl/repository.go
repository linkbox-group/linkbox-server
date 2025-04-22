package acl

import (
	"context"
	"github.com/linkbox-group/linkbox-server/model"
	"github.com/linkbox-group/linkbox-server/rpc-gen/common/pagination"
)

type UserRepositoryItf interface {
	CreateItem(ctx context.Context, req *model.Item) (err error)
	GetItem(ctx context.Context, req *model.Item) (err error)
	UpdateItem(ctx context.Context, req *model.Item) (err error)
	DeleteItem(ctx context.Context, req *model.Item) (err error)
	GetItemsByTags(context.Context, string, []string, *pagination.PaginationRequest) ([]model.Item, int, error)
	GetItemsByOrganization(context.Context, string,string,int,int) ([]model.Item, int, error)
	SearchItemsByTitle(context.Context, string, string, int, int) ([]model.Item, int, error)
}
