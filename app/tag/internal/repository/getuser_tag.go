package repository

import (
	"context"
	"github.com/linkbox-group/linkbox-server/rpc-gen/common/pagination"

	"github.com/linkbox-group/linkbox-server/model"
)

func (r *TagRepository) GetUserTag(ctx context.Context, tag *model.Tag, paginationReq *pagination.PaginationRequest, searchQuery *string) (tags []model.Tag, err error) {

	query := r.db.WithContext(ctx).Model(&model.Tag{}).Where("user_id = ?", tag.UserID)

	var total int64
	if err = query.Count(&total).Error; err != nil {
		return nil, err
	}
	offset := (paginationReq.GetPage() - 1) * paginationReq.GetPageSize()
	tags = make([]model.Tag, 0)
	if err = query.Offset(int(offset)).Limit(int(paginationReq.GetPageSize())).Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}
