package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/linkbox-group/linkbox-server/rpc-gen/item"
	"github.com/linkbox-group/linkbox-server/tag/internal/infra/rpc"

	"github.com/linkbox-group/linkbox-server/model"
	"github.com/linkbox-group/linkbox-server/rpc-gen/common/pagination"
)

var (
	ErrDbCreateTagFailed      = errors.New("create tag failed")
	ErrDbGetTagFailed         = errors.New("get tag failed")
	ErrDbUpdataTagFailed      = errors.New("updata tag failed")
	ErrDbDeleteTagFailed      = errors.New("delete tag failed")
	ErrDbAddTagsToItemsFailed = errors.New("add tags to items failed")
)

func (s *TagService) CreateTagService(ctx context.Context, tag *model.Tag) (err error) {
	err = s.repo.CreateTag(ctx, tag)
	if err != nil {
		return fmt.Errorf("%w:%w", ErrDbCreateTagFailed, err)
	}
	return nil
}

func (s *TagService) GetTagService(ctx context.Context, tag *model.Tag) (err error) {
	err = s.repo.GetTag(ctx, tag)
	if err != nil {
		return fmt.Errorf("%w:%w", ErrDbGetTagFailed, err)
	}
	return nil
}

func (s *TagService) UpdateTagService(ctx context.Context, tag *model.Tag) (err error) {
	err = s.repo.UpdateTag(ctx, tag)
	if err != nil {
		return fmt.Errorf("%w:%w", ErrDbUpdataTagFailed, err)
	}
	return nil
}
func (s *TagService) DeleteTagService(ctx context.Context, tag *model.Tag) (err error) {
	err = s.repo.DeleteTag(ctx, tag)
	if err != nil {
		return fmt.Errorf("%w:%w", ErrDbDeleteTagFailed, err)
	}
	return nil
}
func (s *TagService) GetUserTagService(ctx context.Context, tag *model.Tag, paginationReq *pagination.PaginationRequest, searchQuery *string) (tags []model.Tag, err error) {
	tags, err = s.repo.GetUserTag(ctx, tag, paginationReq, searchQuery)
	if err != nil {
		return nil, fmt.Errorf("%w:%w", ErrDbGetTagFailed, err)
	}
	return tags, nil
}
func (s *TagService) AddTagsToItemsService(ctx context.Context, tag *model.Tag, tagNames []string, itemIds []string) (successCount int32, failedItemIds []string, err error) {
	//TODO !!使用事务
	tagIds := make([]string, 0)
	for _, tagName := range tagNames {
		tag := &model.Tag{
			Name:   tagName,
			UserID: tag.UserID,
		}
		err := s.repo.CreateTag(ctx, tag)
		if err != nil {
			return 0, nil, fmt.Errorf("%w:%w", ErrDbAddTagsToItemsFailed, err)
		}
		tagIds = append(tagIds, tag.ID)
	}
	successCount, failedItemIds, err = s.repo.AddTagsToItems(ctx, tag, tagIds, itemIds)
	if err != nil {
		return 0, nil, fmt.Errorf("%w:%w", ErrDbAddTagsToItemsFailed, err)
	}

	for _, itemID := range itemIds {
		req := item.UpdateItemRequest{
			Id:     itemID,
			UserId: tag.UserID,
			Tags:   tagNames,
		}

		_, err := rpc.ItemClient.UpdateItem(
			ctx,
			&req,
		)
		if err != nil {
			return 0, nil, fmt.Errorf("update item %w:%w", ErrDbAddTagsToItemsFailed, err)
		}

	}

	return successCount, failedItemIds, nil
}
