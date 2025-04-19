package repository

import (
	"context"
	"errors"

	"github.com/linkbox-group/linkbox-server/model"
	"gorm.io/gorm"
)

func (s *TagRepository) AddTagsToItems(ctx context.Context, tag *model.Tag, tagsIds []string, itemIds []string) (successCount int32, failedItemIds []string, err error) {
	if tag == nil {
		err = errors.New("tag is nil")
		return 0, nil, err
	}
	if tag.UserID == "" {
		err = errors.New("user id is empty")
		return 0, nil, err
	}
	if len(tagsIds) == 0 {
		return 0, nil, nil
	}
	if len(itemIds) == 0 {
		return 0, nil, nil
	}
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return 0, nil, tx.Error
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var count int64
	if err = tx.Model(&model.Tag{}).Where("id IN ? AND user_id = ?", tagsIds, tag.UserID).Count(&count).Error; err != nil {
		return 0, nil, err
	}
	if int(count) != len(tagsIds) {
		return 0, nil, errors.New("tag not found")
	}

	for _, item := range itemIds {
		var itemEntity model.Item
		if err = tx.Where("id = ? AND user_id = ?", item, tag.UserID).First(&itemEntity).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				failedItemIds = append(failedItemIds, item)
				continue
			}
			return 0, failedItemIds, err
		}

		for _, tagItem := range tagsIds {
			itemTag := model.ItemTag{
				ItemID: item,
				TagID:  tagItem,
			}
			if err = tx.Create(&itemTag).Error; err != nil {
				failedItemIds = append(failedItemIds, item)
				break
			}
		}
	}
	if err = tx.Commit().Error; err != nil {
		return 0, failedItemIds, err
	}

	successCount = int32(len(itemIds) - len(failedItemIds))
	return successCount, failedItemIds, nil
}
