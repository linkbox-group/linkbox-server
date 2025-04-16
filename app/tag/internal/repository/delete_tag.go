package repository

import (
	"context"
	"errors"

	"github.com/linkbox-group/linkbox-server/model"
)

func (r *TagRepository) DeleteTag(ctx context.Context, tag *model.Tag) (err error) {
	if tag == nil {
		return errors.New("tag cannot be nil")
	}
	if tag.ID == "" {
		return errors.New("tag id cannot be nil")
	}
	if tag.UserID == "" {
		return errors.New("tag user_id cannot be nil")
	}
	result := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", tag.ID, tag.UserID).Delete(tag)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("tag not found")
	}
	return nil
}
