package repository

import (
	"context"
	"github.com/linkbox-group/linkbox-server/model"
)

func (r *TagRepository) CreateTag(ctx context.Context, tag *model.Tag) (err error) {
	err = r.db.Create(tag).Error
	return err
}
