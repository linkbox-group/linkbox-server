package delivery

import (
	"context"
	"log"

	"github.com/linkbox-group/linkbox-server/model"
	"github.com/linkbox-group/linkbox-server/rpc-gen/tag"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TagDelivery implements the last service interface defined in the IDL.

// CreateTag implements the TagDelivery interface.
func (s *TagDelivery) CreateTag(ctx context.Context, req *tag.CreateTagRequest) (resp *tag.CreateTagResponse, err error) {
	tagEntity := model.Tag{
		Name:   req.Name,
		Color:  req.Color,
		UserID: req.UserId,
	}
	err = s.service.CreateTagService(ctx, &tagEntity)
	if err != nil {
		log.Println(err)
		return nil, err

	}
	resp = &tag.CreateTagResponse{
		Result: &tag.CreateTagResponse_Tag{
			Tag: &tag.Tag{
				Id:          tagEntity.ID,
				UserId:      tagEntity.UserID,
				Name:        tagEntity.Name,
				Description: "",
				Color:       *tagEntity.Color,
				ItemCount:   int32(len(tagEntity.Items)),
				CreatedAt:   timestamppb.New(tagEntity.CreatedAt),
				UpdatedAt:   timestamppb.New(tagEntity.UpdatedAt),
			}}}
	return resp, nil
}

// GetTag implements the TagDelivery interface.
func (s *TagDelivery) GetTag(ctx context.Context, req *tag.GetTagRequest) (resp *tag.GetTagResponse, err error) {
	// TODO: Your code here...
	tagEntity := model.Tag{
		BaseModel: model.BaseModel{
			ID: req.Id,
		},
		UserID: req.UserId,
	}
	err = s.service.GetTagService(ctx, &tagEntity)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	resp = &tag.GetTagResponse{
		Result: &tag.GetTagResponse_Tag{
			Tag: &tag.Tag{
				Id:          tagEntity.ID,
				UserId:      tagEntity.UserID,
				Name:        tagEntity.Name,
				Description: "",
				Color:       *tagEntity.Color,
				ItemCount:   int32(len(tagEntity.Items)),
				CreatedAt:   timestamppb.New(tagEntity.CreatedAt),
				UpdatedAt:   timestamppb.New(tagEntity.UpdatedAt),
			}}}
	return resp, nil
	//return
}

// UpdateTag implements the TagDelivery interface.
func (s *TagDelivery) UpdateTag(ctx context.Context, req *tag.UpdateTagRequest) (resp *tag.UpdateTagResponse, err error) {
	// TODO: Your code here...
	tagEntity := model.Tag{
		BaseModel: model.BaseModel{
			ID: req.Id,
		},
		UserID: req.UserId,
		Name:   *req.Name,
		Color:  req.Color,
	}
	err = s.service.UpdateTagService(ctx, &tagEntity)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	resp = &tag.UpdateTagResponse{
		Result: &tag.UpdateTagResponse_Tag{
			Tag: &tag.Tag{
				Id:          tagEntity.ID,
				UserId:      tagEntity.UserID,
				Name:        tagEntity.Name,
				Description: "",
				Color:       *tagEntity.Color,
				ItemCount:   int32(len(tagEntity.Items)),
				CreatedAt:   timestamppb.New(tagEntity.CreatedAt),
				UpdatedAt:   timestamppb.New(tagEntity.UpdatedAt)}}}
	return resp, nil
}

// DeleteTag implements the TagDelivery interface.
func (s *TagDelivery) DeleteTag(ctx context.Context, req *tag.DeleteTagRequest) (resp *tag.DeleteTagResponse, err error) {
	// TODO: Your code here...
	tagEntity := model.Tag{
		BaseModel: model.BaseModel{
			ID: req.Id,
		},
		UserID: req.UserId,
	}
	err = s.service.DeleteTagService(ctx, &tagEntity)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	resp = &tag.DeleteTagResponse{
		Result: &tag.DeleteTagResponse_Success{
			Success: true,
		},
	}
	return resp, nil
}

// GetUserTags implements the TagDelivery interface.
func (s *TagDelivery) GetUserTags(ctx context.Context, req *tag.GetUserTagsRequest) (resp *tag.GetUserTagsResponse, err error) {
	tagEntity := model.Tag{
		UserID: req.UserId,
	}
	var tags []model.Tag
	tags, err = s.service.GetUserTagService(ctx, &tagEntity, req.Pagination, req.SearchQuery)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var pbTags []*tag.Tag
	for _, t := range tags {
		pbTags = append(pbTags, &tag.Tag{
			Id:          t.ID,
			UserId:      t.UserID,
			Name:        t.Name,
			Description: "",
			Color:       *t.Color,
			ItemCount:   int32(len(t.Items)),
			CreatedAt:   timestamppb.New(t.CreatedAt),
			UpdatedAt:   timestamppb.New(t.UpdatedAt),
		})
	}
	resp = &tag.GetUserTagsResponse{
		Result: &tag.
			GetUserTagsResponse_Tags{
			Tags: &tag.TagsData{
				Tags: pbTags,
			},
		},
	}
	return resp, nil
	// TODO: Your code here...
}

// AddTagsToItems implements the TagDelivery interface.
func (s *TagDelivery) AddTagsToItems(ctx context.Context, req *tag.AddTagsToItemsRequest) (resp *tag.AddTagsToItemsResponse, err error) {
	// TODO: Your code here...
	tagEntity := model.Tag{
		UserID: req.UserId,
	}

	successCount, failedItemIds, err := s.service.AddTagsToItemsService(ctx, &tagEntity, req.Tags, req.ItemIds)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	resp = &tag.AddTagsToItemsResponse{
		Result: &tag.AddTagsToItemsResponse_Data{
			Data: &tag.TagOperationResult{
				SuccessCount:  successCount,
				FailureCount:  int32(len(req.ItemIds)) - successCount,
				FailedItemIds: failedItemIds,
			},
		},
	}
	return resp, nil
}

// RemoveTagsFromItems implements the TagDelivery interface.
func (s *TagDelivery) RemoveTagsFromItems(ctx context.Context, req *tag.RemoveTagsFromItemsRequest) (resp *tag.RemoveTagsFromItemsResponse, err error) {
	// TODO: Your code here...
	return
}

// GetItemTags implements the TagDelivery interface.
func (s *TagDelivery) GetItemTags(ctx context.Context, req *tag.GetItemTagsRequest) (resp *tag.GetItemTagsResponse, err error) {
	// TODO: Your code here...
	return
}

// MergeTags implements the TagDelivery interface.
func (s *TagDelivery) MergeTags(ctx context.Context, req *tag.MergeTagsRequest) (resp *tag.MergeTagsResponse, err error) {
	// TODO: Your code here...
	return
}

// GetTagStats implements the TagDelivery interface.
func (s *TagDelivery) GetTagStats(ctx context.Context, req *tag.GetTagStatsRequest) (resp *tag.GetTagStatsResponse, err error) {
	// TODO: Your code here...
	return
}

// GetRelatedTags implements the TagDelivery interface.
func (s *TagDelivery) GetRelatedTags(ctx context.Context, req *tag.GetRelatedTagsRequest) (resp *tag.GetRelatedTagsResponse, err error) {
	// TODO: Your code here...
	return
}

// SuggestTags implements the TagDelivery interface.
func (s *TagDelivery) SuggestTags(ctx context.Context, req *tag.SuggestTagsRequest) (resp *tag.SuggestTagsResponse, err error) {
	// TODO: Your code here...
	return
}
