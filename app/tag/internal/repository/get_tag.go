package repository

import (
	"context"
	"github.com/linkbox-group/linkbox-server/model"
)

func (r *TagRepository) GetTag(ctx context.Context, tag *model.Tag) (err error) {
	err = r.db.First(tag).Error
	return err
}
