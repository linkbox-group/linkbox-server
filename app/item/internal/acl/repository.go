package acl

import (
	"context"
	"github.com/linkbox-group/linkbox-server/model"
	"github.com/linkbox-group/linkbox-server/rpc-gen/common/pagination"
	itemmodel "github.com/linkbox-group/linkbox-server/rpc-gen/model"
)

type UserRepositoryItf interface {
	CreateItem(ctx context.Context, req *model.Item) (err error)
	GetItem(ctx context.Context, req *model.Item) (err error)
	UpdateItem(ctx context.Context, req *model.Item) (err error)
	DeleteItem(ctx context.Context, req *model.Item) (err error)
	GetItemsByTags(context.Context, string, []string, *pagination.PaginationRequest) ([]model.Item, int, error)
	GetItemsByOrganization(context.Context, string, string, int, int) ([]model.Item, int, error)
	SearchItems(ctx context.Context, userID string, query string, pageNum int, pageSize int) ([]model.Item, int, error)
	RecoverItemsBatch(context.Context, string, []string) (err error)
	DeleteItemsBatch(context.Context, string, []string) (err error)
	GetDeletedItems(ctx context.Context, userID string, pagination *pagination.PaginationRequest) ([]*model.Item, int, error)
}
type EsRepositoryItf interface {
	SearchItems(ctx context.Context, UserID string, query string, itemType itemmodel.ItemType, pageNum int, pageSize int) (items []model.Item, count int, err error)
}
