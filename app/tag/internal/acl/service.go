package acl

import (
	"context"

	"github.com/linkbox-group/linkbox-server/model"
	"github.com/linkbox-group/linkbox-server/rpc-gen/common/pagination"
)

type TagServiceItf interface {
	CreateTagService(ctx context.Context, tag *model.Tag) (err error)
	GetTagService(ctx context.Context, tag *model.Tag) (err error)
	UpdateTagService(ctx context.Context, tag *model.Tag) (err error)
	DeleteTagService(ctx context.Context, tag *model.Tag) (err error)
	GetUserTagService(ctx context.Context, tag *model.Tag, paginationReq *pagination.PaginationRequest, searchQuery *string) (tags []model.Tag, err error)
	AddTagsToItemsService(ctx context.Context, tag *model.Tag, tagIds []string, itemIds []string) (successCount int32, failedItemIds []string, err error)
}
