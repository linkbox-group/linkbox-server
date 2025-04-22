package acl

import (
	"context"
	"github.com/linkbox-group/linkbox-server/model"
	"github.com/linkbox-group/linkbox-server/rpc-gen/common/pagination"
)

type UserServiceItf interface {
	CreateItem(ctx context.Context, req *model.Item) (err error)
	GetItem(ctx context.Context, req *model.Item) (err error)
	UpdateItem(ctx context.Context, req *model.Item) (err error)
	DeleteItem(ctx context.Context, req *model.Item) (err error)
	GetItemsByTags(ctx context.Context, userID string, tagIDs []string, pagination *pagination.PaginationRequest) (items []model.Item, total int, err error)
	GetItemsByOrganization(ctx context.Context, userID string, organizationID string, pageNum int, pageSize int) (items []model.Item, total int, err error)
	SearchItems(ctx context.Context, userID string, query string, pageNum int, pageSize int) (items []model.Item, total int, err error)
}
