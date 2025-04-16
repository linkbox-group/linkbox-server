package repository

import (
	"context"
	"errors"
	"github.com/linkbox-group/linkbox-server/model"
	"gorm.io/gorm"
)

func (r *TagRepository) UpdateTag(ctx context.Context, tag *model.Tag) (err error) {
	if tag == nil {
		return errors.New("tag cannot be nil")
	}
	if tag.ID == "" {
		return errors.New("tag id cannot be nil")
	}
	if tag.Name == "" {
		return errors.New("tag name cannot be nil")
	}
	updateData := model.Tag{
		BaseModel: model.BaseModel{
			ID: tag.ID,
		},
		Name:   tag.Name,
		UserID: tag.UserID,
		Color:  tag.Color,
	}

	result := r.db.WithContext(ctx).
		Model(&model.Tag{}).
		Where("id = ?", tag.ID).
		Updates(updateData)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
