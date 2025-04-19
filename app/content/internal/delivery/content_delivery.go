package delivery

import (
	"context"
	"github.com/linkbox-group/linkbox-server/rpc-gen/common/pagination"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/linkbox-group/linkbox-server/model"
	"github.com/linkbox-group/linkbox-server/rpc-gen/common/cError"
	"github.com/linkbox-group/linkbox-server/rpc-gen/content"
)

func (d *ContentDelivery) CreateItem(ctx context.Context, req *content.CreateItemRequest) (resp *content.CreateItemResponse, err error) {
	item := model.Item{
		UserID: req.UserId,
		Title:  req.Title,
		URL:    req.Url,
	}
	err = d.s.CreateItem(ctx, &item)

	if err != nil {
		return &content.CreateItemResponse{
			Result: &content.CreateItemResponse_Error{
				Error: &cError.Error{
					Code:    40000,
					Message: err.Error(),
				},
			},
		}, err
	}
	return &content.CreateItemResponse{
		Result: &content.CreateItemResponse_Item{
			Item: &content.Item{
				Id:          item.ID,
				UserId:      item.UserID,
				Title:       item.Title,
				Description: "",
				Url:         item.URL,
				CreatedAt:   timestamppb.New(item.CreatedAt),
				UpdatedAt:   timestamppb.New(item.UpdatedAt),
			},
		},
	}, nil

}

// GetItem implements the ContentDelivery interface.
func (d *ContentDelivery) GetItem(ctx context.Context, req *content.GetItemRequest) (resp *content.GetItemResponse, err error) {
	item := model.Item{
		BaseModel: model.BaseModel{
			ID: req.Id,
		},
		UserID: req.UserId,
	}
	err = d.s.GetItem(ctx, &item)

	if err != nil {
		return &content.GetItemResponse{
			Result: &content.GetItemResponse_Error{
				Error: &cError.Error{
					Code:    40001,
					Message: err.Error(),
				},
			},
		}, err
	}

	return &content.GetItemResponse{
		Result: &content.GetItemResponse_Item{
			Item: &content.Item{
				Id:          item.ID,
				UserId:      item.UserID,
				Title:       item.Title,
				Description: "",
				Url:         item.URL,
				CreatedAt:   timestamppb.New(item.CreatedAt),
				UpdatedAt:   timestamppb.New(item.UpdatedAt),
			},
		},
	}, nil
}

// UpdateItem implements the ContentDelivery interface.
func (d *ContentDelivery) UpdateItem(ctx context.Context, req *content.UpdateItemRequest) (resp *content.UpdateItemResponse, err error) {
	item := model.Item{
		BaseModel: model.BaseModel{
			ID: req.Id,
		},
		UserID: req.UserId,
		Title:  req.Title,
	}

	err = d.s.UpdateItem(ctx, &item)

	if err != nil {
		return &content.UpdateItemResponse{
			Result: &content.UpdateItemResponse_Error{
				Error: &cError.Error{
					Code:    40000,
					Message: err.Error(),
				},
			},
		}, err
	}

	return &content.UpdateItemResponse{
		Result: &content.UpdateItemResponse_Item{
			Item: &content.Item{
				Id:          item.ID,
				UserId:      item.UserID,
				Title:       item.Title,
				Description: "",
				Url:         item.URL,
				CreatedAt:   timestamppb.New(item.CreatedAt),
				UpdatedAt:   timestamppb.New(item.UpdatedAt),
			},
		},
	}, nil
}

// DeleteItem implements the ContentDelivery interface.
func (d *ContentDelivery) DeleteItem(ctx context.Context, req *content.DeleteItemRequest) (resp *content.DeleteItemResponse, err error) {
	// TODO: Your code here...
	item := model.Item{
		BaseModel: model.BaseModel{
			ID: req.Id,
		},
		UserID: req.UserId,
	}

	err = d.s.DeleteItem(ctx, &item)

	if err != nil {
		return &content.DeleteItemResponse{
			Result: &content.DeleteItemResponse_Success{
				Success: false,
			},
		}, err
	}

	return &content.DeleteItemResponse{
		Result: &content.DeleteItemResponse_Success{
			Success: true,
		},
	}, nil
}

// GetItems implements the ContentDelivery interface.
func (s *ContentDelivery) GetItems(ctx context.Context, req *content.GetItemsRequest) (resp *content.GetItemsResponse, err error) {
	// TODO: Your code here...
	return
}

// BatchDeleteItems implements the ContentDelivery interface.
func (s *ContentDelivery) BatchDeleteItems(ctx context.Context, req *content.BatchDeleteItemsRequest) (resp *content.BatchDeleteItemsResponse, err error) {
	// TODO: Your code here...
	return
}

// GetItemsByTags implements the ContentDelivery interface.
func (d *ContentDelivery) GetItemsByTags(ctx context.Context, req *content.GetItemsByTagsRequest) (resp *content.GetItemsByTagsResponse, err error) {
	userID := req.UserId
	tagIDs := req.Tags
	paginationReq := req.Pagination // 保存请求中的分页信息
	// 注意：传递给 service 层的 paginationReq 类型应为 *commonPagination.PaginationRequest
	items, total, err := d.s.GetItemsByTags(ctx, userID, tagIDs, paginationReq)
	if err != nil {
		return &content.GetItemsByTagsResponse{
			Result: &content.GetItemsByTagsResponse_Error{
				Error: &cError.Error{
					Code:    40000,
					Message: err.Error(),
				},
			},
		}, err
	}

	// 转换结果为响应格式
	respItems := make([]*content.Item, 0, len(items))
	for _, dbItem := range items {
		tagStrings := make([]string, 0, len(dbItem.Tags))
		for _, tag := range dbItem.Tags {
			tagStrings = append(tagStrings, tag.Name)
		}

		respItems = append(respItems, &content.Item{
			Id:     dbItem.ID,
			UserId: dbItem.UserID,
			Title:  dbItem.Title,
			//Description: dbItem.Description,
			Url:       dbItem.URL,
			Tags:      tagStrings,
			CreatedAt: timestamppb.New(dbItem.CreatedAt),
			UpdatedAt: timestamppb.New(dbItem.UpdatedAt),
		})
	}

	// 构造分页响应信息
	var currentPage, currentPageSize int64 = 1, 10 // 默认值 (int64)
	if paginationReq != nil {                      // 先检查 nil
		if paginationReq.Page > 0 {
			currentPage = int64(paginationReq.Page) // 从 int32 转为 int64
		}
		if paginationReq.PageSize > 0 {
			currentPageSize = int64(paginationReq.PageSize) // 从 int32 转为 int64
			// 可选：添加最大页面大小限制
			// const maxPageSize = 100
			// if currentPageSize > maxPageSize {
			// 	currentPageSize = maxPageSize
			// }
		}
	}
	// 使用正确的类型 commonPagination.PaginationResponse 和字段名
	paginationResp := &pagination.PaginationMeta{
		TotalPages: int32(total), // 从 int 转为 int64
		Page:       int32(currentPage),
		PageSize:   int32(currentPageSize),
	}

	// 构造成功响应，使用 ItemsPage 结构
	return &content.GetItemsByTagsResponse{
		Result: &content.GetItemsByTagsResponse_ItemsPage{
			ItemsPage: &content.ItemsPage{
				Items:      respItems,
				Pagination: paginationResp,
			},
		},
	}, nil
}

// ExtractMetadata implements the ContentDelivery interface.
func (s *ContentDelivery) ExtractMetadata(ctx context.Context, req *content.ExtractMetadataRequest) (resp *content.ExtractMetadataResponse, err error) {
	// TODO: Your code here...
	return
}

// GetRecentItems implements the ContentDelivery interface.
func (s *ContentDelivery) GetRecentItems(ctx context.Context, req *content.GetRecentItemsRequest) (resp *content.GetRecentItemsResponse, err error) {
	// TODO: Your code here...
	return
}

// BatchUpdateItems implements the ContentDelivery interface.
func (s *ContentDelivery) BatchUpdateItems(ctx context.Context, req *content.BatchUpdateItemsRequest) (resp *content.BatchUpdateItemsResponse, err error) {
	// TODO: Your code here...
	return
}

// ImportFromFile implements the ContentDelivery interface.
func (s *ContentDelivery) ImportFromFile(ctx context.Context, req *content.ImportFromFileRequest) (resp *content.ImportFromFileResponse, err error) {
	// TODO: Your code here...
	return
}

// ExportToFile implements the ContentDelivery interface.
func (s *ContentDelivery) ExportToFile(ctx context.Context, req *content.ExportToFileRequest) (resp *content.ExportToFileResponse, err error) {
	// TODO: Your code here...
	return
}

// SearchItems implements the ContentDelivery interface.
func (s *ContentDelivery) SearchItems(ctx context.Context, req *content.SearchItemsRequest) (resp *content.SearchItemsResponse, err error) {
	// TODO: Your code here...
	return
}

// AddItemNote implements the ContentDelivery interface.
func (s *ContentDelivery) AddItemNote(ctx context.Context, req *content.AddItemNoteRequest) (resp *content.AddItemNoteResponse, err error) {
	// TODO: Your code here...
	return
}

// UpdateItemNote implements the ContentDelivery interface.
func (s *ContentDelivery) UpdateItemNote(ctx context.Context, req *content.UpdateItemNoteRequest) (resp *content.UpdateItemNoteResponse, err error) {
	// TODO: Your code here...
	return
}

// GetItemNote implements the ContentDelivery interface.
func (s *ContentDelivery) GetItemNote(ctx context.Context, req *content.GetItemNoteRequest) (resp *content.GetItemNoteResponse, err error) {
	// TODO: Your code here...
	return
}
