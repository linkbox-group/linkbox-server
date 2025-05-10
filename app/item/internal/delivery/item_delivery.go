package delivery

import (
	"context"
	"github.com/linkbox-group/linkbox-server/item/internal/infra/rpc"
	"github.com/linkbox-group/linkbox-server/item/pkg/log"
	"github.com/linkbox-group/linkbox-server/model/treemodel"
	"github.com/linkbox-group/linkbox-server/rpc-gen/organization"
	"time"

	"github.com/linkbox-group/linkbox-server/rpc-gen/common/pagination"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/linkbox-group/linkbox-server/model"
	"github.com/linkbox-group/linkbox-server/rpc-gen/common/cError"
	"github.com/linkbox-group/linkbox-server/rpc-gen/item"
	itemmodel "github.com/linkbox-group/linkbox-server/rpc-gen/model"
)

func (d *ItemDelivery) CreateItem(ctx context.Context, req *item.CreateItemRequest) (resp *item.CreateItemResponse, err error) {

	itemModel := model.Item{
		UserID:         req.UserId,
		ItemType:       req.Type,
		Title:          req.Title,
		URL:            req.Url,
		TagNames:       req.Tags,
		OrganizationID: req.OrganizationId,
		Note:           req.Note,
		CreatedAt:      model.CustomTime(time.Now()),
		UpdatedAt:      model.CustomTime(time.Now()),
	}
	if req.OrganizationId == "" || req.OrganizationId == treemodel.ROOT_ID {
		orgID, err := rpc.OrganizationClient.GetDefaultOrgID(ctx, &organization.GetDefaultOrgIDReq{
			UserId: req.UserId,
			Code:   treemodel.ROOT_ID,
		})
		if err != nil {
			return nil, err
		}
		itemModel.OrganizationID = orgID.Id
	}
	err = d.s.CreateItem(ctx, &itemModel)

	if err != nil {
		return &item.CreateItemResponse{
			Result: &item.CreateItemResponse_Error{
				Error: &cError.Error{
					Code:    40000,
					Message: err.Error(),
				},
			},
		}, err
	}
	return &item.CreateItemResponse{
		Result: &item.CreateItemResponse_Item{
			Item: &itemmodel.Item{
				Id:               itemModel.ID,
				UserId:           itemModel.UserID,
				Title:            itemModel.Title,
				Type:             itemModel.ItemType,
				Description:      "",
				Url:              itemModel.URL,
				Note:             itemModel.Note,
				OrganizationPath: itemModel.OrganizationPath,
				TagNames:         itemModel.TagNames,
				CreatedAt:        timestamppb.New(itemModel.CreatedAt.Time()),
				UpdatedAt:        timestamppb.New(itemModel.UpdatedAt.Time()),
			},
		},
	}, nil

}

// GetItem implements the ItemDelivery interface.
func (d *ItemDelivery) GetItem(ctx context.Context, req *item.GetItemRequest) (resp *item.GetItemResponse, err error) {
	itemModel := model.Item{
		BaseModel: model.BaseModel{
			ID: req.Id,
		},
		UserID: req.UserId,
	}
	err = d.s.GetItem(ctx, &itemModel)

	if err != nil {
		return &item.GetItemResponse{
			Result: &item.GetItemResponse_Error{
				Error: &cError.Error{
					Code:    40001,
					Message: err.Error(),
				},
			},
		}, err
	}

	return &item.GetItemResponse{
		Result: &item.GetItemResponse_Item{
			Item: &itemmodel.Item{
				Id:          itemModel.ID,
				UserId:      itemModel.UserID,
				Title:       itemModel.Title,
				Type:        itemModel.ItemType,
				Description: "",
				Url:         itemModel.URL,
				Note:        itemModel.Note,
				TagNames:    itemModel.TagNames,
				CreatedAt:   timestamppb.New(itemModel.CreatedAt.Time()),
				UpdatedAt:   timestamppb.New(itemModel.UpdatedAt.Time()),
			},
		},
	}, nil
}

// UpdateItem implements the ItemDelivery interface.
func (d *ItemDelivery) UpdateItem(ctx context.Context, req *item.UpdateItemRequest) (resp *item.UpdateItemResponse, err error) {

	itemModel := model.Item{
		BaseModel: model.BaseModel{
			ID: req.Id,
		},
		UserID:         req.UserId,
		Title:          req.Title,
		URL:            req.Url,
		Note:           req.Note,
		OrganizationID: req.OrganizationId,
		UpdatedAt:      model.CustomTime(time.Now()),
		TagNames:       req.GetTags(),
	}
	if req.OrganizationId == "" || req.OrganizationId == treemodel.ROOT_ID {
		orgID, err := rpc.OrganizationClient.GetDefaultOrgID(ctx, &organization.GetDefaultOrgIDReq{
			UserId: req.UserId,
			Code:   treemodel.ROOT_ID,
		})
		if err != nil {
			log.Log().Error(err.Error(), "req", req)
			return nil, err
		}
		itemModel.OrganizationID = orgID.Id
	}
	org, err := rpc.OrganizationClient.GetOrganization(ctx, &organization.GetOrganizationRequest{
		Id:     itemModel.OrganizationID,
		UserId: itemModel.UserID,
	})
	if err != nil {
		log.Log().Error(err.Error(), "req", req)
		return nil, err
	}
	itemModel.OrganizationPath = org.GetOrganization().TreeNames

	err = d.s.UpdateItem(ctx, &itemModel)

	if err != nil {
		return &item.UpdateItemResponse{
			Result: &item.UpdateItemResponse_Error{
				Error: &cError.Error{
					Code:    40000,
					Message: err.Error(),
				},
			},
		}, err
	}

	return &item.UpdateItemResponse{
		Result: &item.UpdateItemResponse_Item{
			Item: &itemmodel.Item{
				Id:        itemModel.ID,
				UserId:    itemModel.UserID,
				Title:     itemModel.Title,
				Url:       itemModel.URL,
				CreatedAt: timestamppb.New(itemModel.CreatedAt.Time()),
				UpdatedAt: timestamppb.New(itemModel.UpdatedAt.Time()),
			},
		},
	}, nil
}

// DeleteItem implements the ItemDelivery interface.
func (d *ItemDelivery) DeleteItem(ctx context.Context, req *item.DeleteItemRequest) (resp *item.DeleteItemResponse, err error) {
	// TODO: Your code here...
	itemModel := model.Item{
		BaseModel: model.BaseModel{
			ID: req.Id,
		},
		UserID: req.UserId,
	}

	err = d.s.DeleteItem(ctx, &itemModel)

	if err != nil {
		return &item.DeleteItemResponse{
			Result: &item.DeleteItemResponse_Success{
				Success: false,
			},
		}, err
	}

	return &item.DeleteItemResponse{
		Result: &item.DeleteItemResponse_Success{
			Success: true,
		},
	}, nil
}

// GetItems implements the ItemDelivery interface.
func (d *ItemDelivery) GetItems(ctx context.Context, req *item.GetItemsRequest) (resp *item.GetItemsResponse, err error) {
	// TODO: Your code here...
	return
}

// GetItemsByTags implements the ContentDelivery interface.
func (d *ItemDelivery) GetItemsByTags(ctx context.Context, req *item.GetItemsByTagsRequest) (resp *item.GetItemsByTagsResponse, err error) {
	userID := req.UserId
	tagNames := req.Tags
	paginationReq := req.Pagination // 保存请求中的分页信息
	// 注意：传递给 service 层的 paginationReq 类型应为 *commonPagination.PaginationRequest
	items, total, err := d.s.GetItemsByTags(ctx, userID, tagNames, paginationReq)
	if err != nil {
		return &item.GetItemsByTagsResponse{
			Result: &item.GetItemsByTagsResponse_Error{
				Error: &cError.Error{
					Code:    40000,
					Message: err.Error(),
				},
			},
		}, err
	}

	// 转换结果为响应格式
	respItems := make([]*itemmodel.Item, 0, len(items))
	for _, dbItem := range items {
		tagStrings := make([]string, 0, len(dbItem.Tags))
		for _, tag := range dbItem.Tags {
			tagStrings = append(tagStrings, tag.Name)
		}

		respItems = append(respItems, &itemmodel.Item{
			Id:     dbItem.ID,
			UserId: dbItem.UserID,
			Title:  dbItem.Title,
			Note:   dbItem.Note,
			Type:   dbItem.ItemType,
			//Description: dbItem.Description,
			Url:       dbItem.URL,
			Tags:      tagStrings,
			CreatedAt: timestamppb.New(dbItem.CreatedAt.Time()),
			UpdatedAt: timestamppb.New(dbItem.UpdatedAt.Time()),
		})
	}

	// 构造分页响应信息
	// 使用正确的类型 commonPagination.PaginationResponse 和字段名
	paginationResp := &pagination.PaginationMeta{
		TotalItems: int32(total), // 从 int 转为 int64
		TotalPages: int32(total) / paginationReq.PageSize,
		Page:       int32(paginationReq.Page),
		PageSize:   int32(paginationReq.PageSize),
	}

	// 构造成功响应，使用 ItemsPage 结构
	return &item.GetItemsByTagsResponse{
		Result: &item.GetItemsByTagsResponse_ItemsPage{
			ItemsPage: &item.ItemsPage{
				Items:      respItems,
				Pagination: paginationResp,
			},
		},
	}, nil
}

// GetItemsByOrganization implements the ItemServiceImpl interface.
func (d *ItemDelivery) GetItemsByOrganization(ctx context.Context, req *item.GetItemsByOrganizationRequest) (resp *item.GetItemsByOrganizationResponse, err error) {
	userID := req.UserId
	orgID := req.OrganizationId
	paginationReq := req.Pagination
	// 构造分页响应信息
	var currentPage, currentPageSize int = 1, 10 // 默认值
	if paginationReq != nil {
		if paginationReq.Page > 0 {
			currentPage = int(paginationReq.Page)
		}
		if paginationReq.PageSize > 0 {
			currentPageSize = int(paginationReq.PageSize)
		}
	}

	items, total, err := d.s.GetItemsByOrganization(ctx, userID, orgID, currentPage, currentPageSize)
	if err != nil {
		return &item.GetItemsByOrganizationResponse{
			Result: &item.GetItemsByOrganizationResponse_Error{
				Error: &cError.Error{
					Code:    40000,
					Message: err.Error(),
				},
			},
		}, err
	}

	// 转换结果为响应格式
	respItems := make([]*itemmodel.Item, 0, len(items))
	for _, dbItem := range items {
		tagStrings := make([]string, 0, len(dbItem.Tags))
		for _, tag := range dbItem.Tags {
			tagStrings = append(tagStrings, tag.Name)
		}
		//TODO: use convert method
		respItems = append(respItems, &itemmodel.Item{
			Id:               dbItem.ID,
			UserId:           dbItem.UserID,
			Type:             dbItem.ItemType,
			Title:            dbItem.Title,
			Note:             dbItem.Note,
			Description:      "",
			TagNames:         dbItem.TagNames,
			OrganizationPath: dbItem.OrganizationPath,
			Url:              dbItem.URL,
			Tags:             tagStrings,
			CreatedAt:        timestamppb.New(dbItem.CreatedAt.Time()),
			UpdatedAt:        timestamppb.New(dbItem.UpdatedAt.Time()),
		})
	}

	paginationResp := &pagination.PaginationMeta{
		TotalItems: int32(total),
		Page:       int32(currentPage),
		PageSize:   int32(currentPageSize),
	}

	return &item.GetItemsByOrganizationResponse{
		Result: &item.GetItemsByOrganizationResponse_ItemsPage{
			ItemsPage: &item.ItemsPage{
				Items:      respItems,
				Pagination: paginationResp,
			},
		},
	}, nil
}

// GetRecentItems implements the ItemDelivery interface.
func (d *ItemDelivery) GetRecentItems(ctx context.Context, req *item.GetRecentItemsRequest) (resp *item.GetRecentItemsResponse, err error) {
	// TODO: Your code here...
	return
}

// ImportFromFile implements the ItemDelivery interface.
func (d *ItemDelivery) ImportFromFile(ctx context.Context, req *item.ImportFromFileRequest) (resp *item.ImportFromFileResponse, err error) {
	// TODO: Your code here...
	return
}

// ExportToFile implements the ItemDelivery interface.
func (d *ItemDelivery) ExportToFile(ctx context.Context, req *item.ExportToFileRequest) (resp *item.ExportToFileResponse, err error) {
	// TODO: Your code here...
	return
}

// SearchItems implements the ItemDelivery interface.
func (d *ItemDelivery) SearchItems(ctx context.Context, req *item.SearchItemsRequest) (resp *item.SearchItemsResponse, err error) {
	userID := req.UserId
	paginationReq := req.Pagination
	// 构造分页响应信息
	var currentPage, currentPageSize int = 1, 10 // 默认值
	if paginationReq != nil {
		if paginationReq.Page > 0 {
			currentPage = int(paginationReq.Page)
		}
		if paginationReq.PageSize > 0 {
			currentPageSize = int(paginationReq.PageSize)
		}
	}
	items, total, err := d.s.SearchItems(ctx, userID, req.Query, req.Type, currentPage, currentPageSize)
	if err != nil {
		return &item.SearchItemsResponse{
			Result: &item.SearchItemsResponse_Error{
				Error: &cError.Error{
					Code:    40000,
					Message: err.Error(),
				},
			},
		}, err
	}
	// 转换结果为响应格式
	respItems := make([]*itemmodel.Item, 0, len(items))
	for _, dbItem := range items {
		tagStrings := make([]string, 0, len(dbItem.Tags))
		for _, tag := range dbItem.Tags {
			tagStrings = append(tagStrings, tag.Name)
		}
		respItems = append(respItems, &itemmodel.Item{
			Id:               dbItem.ID,
			UserId:           dbItem.UserID,
			Title:            dbItem.Title,
			Type:             dbItem.ItemType,
			Description:      "",
			Url:              dbItem.URL,
			TagNames:         dbItem.TagNames,
			OrganizationPath: dbItem.OrganizationPath,
			ThumbnailUrl:     dbItem.ThumbnailURL,
			Note:             dbItem.Note,
			OrganizationId:   dbItem.OrganizationID,
			Tags:             tagStrings,
			CreatedAt:        timestamppb.New(dbItem.CreatedAt.Time()),
			UpdatedAt:        timestamppb.New(dbItem.UpdatedAt.Time()),
		})
	}

	return &item.SearchItemsResponse{
		Result: &item.SearchItemsResponse_Data{
			Data: &item.SearchResult{
				Items: respItems,
				Pagination: &pagination.PaginationMeta{
					TotalItems: int32(total),
					Page:       int32(currentPage),
					PageSize:   int32(currentPageSize),
				},
			},
		},
	}, nil
}
func (d *ItemDelivery) RecoverItem(ctx context.Context, req *item.RecoverItemRequest) (res *item.RecoverItemREsponse, err error) {
	return
}

func (d *ItemDelivery) GetDeletedItems(ctx context.Context, req *item.GetDeletedItemsRequest) (res *item.GetDeletedItemsResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *ItemDelivery) RecoverItemsBatch(ctx context.Context, req *item.RecoverItemsBatchRequest) (res *item.RecoverItemsBatchResponse, err error) {
	userID := req.UserId

	err = d.s.RecoverItemsBatch(ctx, userID, req.Ids)

	if err != nil {
		return &item.RecoverItemsBatchResponse{
			Result: &item.RecoverItemsBatchResponse_Error{
				Error: &cError.Error{
					Code:    40000,
					Message: err.Error(),
				},
			},
		}, err
	}

	return &item.RecoverItemsBatchResponse{
		Result: &item.RecoverItemsBatchResponse_Success{
			Success: true,
		},
	}, nil
}

func (d *ItemDelivery) DeleteItemsBatch(ctx context.Context, req *item.DeleteItemsBatchRequest) (res *item.DeleteItemsBatchResponse, err error) {
	//TODO implement me
	userID := req.UserId

	err = d.s.DeleteItemsBatch(ctx, userID, req.Ids)

	if err != nil {
		return &item.DeleteItemsBatchResponse{
			Result: &item.DeleteItemsBatchResponse_Error{
				Error: &cError.Error{
					Code:    40000,
					Message: err.Error(),
				},
			},
		}, err
	}

	return &item.DeleteItemsBatchResponse{
		Result: &item.DeleteItemsBatchResponse_Success{
			Success: true,
		},
	}, nil
}
