package repository

import (
	"context"
	"github.com/linkbox-group/linkbox-server/model"
)

func (r *TagRepository) GetTag(ctx context.Context, tag *model.Tag) (err error) {
	if tag.Name != "" {
		return r.db.Where("name = ? and user_id = ?", tag.Name, tag.UserID).First(&tag).Error
	}
	err = r.db.Where("id = ? and user_id = ?", tag.ID, tag.UserID).First(tag).Error
	return err
}
