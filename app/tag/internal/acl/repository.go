package acl

import (
	"context"

	"github.com/linkbox-group/linkbox-server/model"
	"github.com/linkbox-group/linkbox-server/rpc-gen/common/pagination"
)

type TagRepositoryItf interface {
	CreateTag(ctx context.Context, tag *model.Tag) (err error)
	GetTag(ctx context.Context, tag *model.Tag) (err error)
	UpdateTag(ctx context.Context, tag *model.Tag) (err error)
	DeleteTag(ctx context.Context, tag *model.Tag) (err error)
	GetUserTag(ctx context.Context, tag *model.Tag, paginationReq *pagination.PaginationRequest, searchQuery *string) (tags []model.Tag, err error)
	AddTagsToItems(ctx context.Context, tag *model.Tag, tagIds []string, itemIds []string) (successCount int32, failedItemIds []string, err error)
}
