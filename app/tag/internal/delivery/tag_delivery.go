package delivery

import (
	"context"
	"log"

	"github.com/linkbox-group/linkbox-server/model"
	"github.com/linkbox-group/linkbox-server/rpc-gen/tag"
)

// TagDelivery implements the last service interface defined in the IDL.

// CreateTag implements the TagDelivery interface.
func (s *TagDelivery) CreateTag(ctx context.Context, req *tag.CreateTagRequest) (resp *tag.CreateTagResponse, err error) {
	tagEntity := model.Tag{
		Name:   req.Name,
		Color:  req.GetColor(),
		UserID: req.UserId,
		Icon:   req.GetIcon(),
	}
	err = s.service.CreateTagService(ctx, &tagEntity)
	if err != nil {
		log.Println(err)
		return nil, err

	}
	resp = &tag.CreateTagResponse{
		Result: &tag.CreateTagResponse_Tag{
			Tag: tagEntity.ConvertTo()}}
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
			Tag: tagEntity.ConvertTo()}}
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
		Name:   req.GetName(),
		Color:  req.GetColor(),
	}
	err = s.service.UpdateTagService(ctx, &tagEntity)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	resp = &tag.UpdateTagResponse{
		Result: &tag.UpdateTagResponse_Tag{
			Tag: tagEntity.ConvertTo()}}
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
		pbTags = append(pbTags, t.ConvertTo())
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
