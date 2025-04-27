package service

import (
	"context"
	"fmt"
	"github.com/linkbox-group/linkbox-server/model"
	"github.com/linkbox-group/linkbox-server/rpc-gen/item"
	"github.com/linkbox-group/linkbox-server/tag/internal/infra/rpc"
)

func (s *TagService) AddTagsToItemsService(ctx context.Context, tag *model.Tag, tagNames []string, itemIds []string) (successCount int32, failedItemIds []string, err error) {
	//TODO !!使用事务
	tagIds := make([]string, 0)
	for _, tagName := range tagNames {
		tag := &model.Tag{
			Name:   tagName,
			UserID: tag.UserID,
		}
		if s.repo.GetTag(ctx, tag) == nil {
			tagIds = append(tagIds, tag.ID)
			continue
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
		itemResp, err := rpc.ItemClient.GetItem(ctx, &item.GetItemRequest{
			Id:     itemID,
			UserId: tag.UserID,
		})
		if err != nil {
			return 0, nil, fmt.Errorf("get item %w:%w", ErrDbAddTagsToItemsFailed, err)
		}
		itemData := itemResp.GetItem()
		req := item.UpdateItemRequest{
			Id:     itemData.Id,
			UserId: itemData.UserId,
			Tags:   append(itemData.Tags, tagIds...),
		}

		_, err = rpc.ItemClient.UpdateItem(
			ctx,
			&req,
		)
		if err != nil {
			return 0, nil, fmt.Errorf("update item %w:%w", ErrDbAddTagsToItemsFailed, err)
		}

	}

	return successCount, failedItemIds, nil
}
