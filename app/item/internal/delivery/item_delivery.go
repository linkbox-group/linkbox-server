package delivery

import (
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/linkbox-group/linkbox-server/model"
	"github.com/linkbox-group/linkbox-server/rpc-gen/common/cError"
	"github.com/linkbox-group/linkbox-server/rpc-gen/item"
)

func (d *ItemDelivery) CreateItem(ctx context.Context, req *item.CreateItemRequest) (resp *item.CreateItemResponse, err error) {
	itemModel := model.Item{
		UserID: req.UserId,
		Title:  req.Title,
		URL:    req.Url,
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
			Item: &item.Item{
				Id:          itemModel.ID,
				UserId:      itemModel.UserID,
				Title:       itemModel.Title,
				Description: "",
				Url:         itemModel.URL,
				CreatedAt:   timestamppb.New(itemModel.CreatedAt),
				UpdatedAt:   timestamppb.New(itemModel.UpdatedAt),
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
			Item: &item.Item{
				Id:          itemModel.ID,
				UserId:      itemModel.UserID,
				Title:       itemModel.Title,
				Description: "",
				Url:         itemModel.URL,
				CreatedAt:   timestamppb.New(itemModel.CreatedAt),
				UpdatedAt:   timestamppb.New(itemModel.UpdatedAt),
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
		UserID: req.UserId,
		Title:  req.Title,
	}

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
			Item: &item.Item{
				Id:          itemModel.ID,
				UserId:      itemModel.UserID,
				Title:       itemModel.Title,
				Description: "",
				Url:         itemModel.URL,
				CreatedAt:   timestamppb.New(itemModel.CreatedAt),
				UpdatedAt:   timestamppb.New(itemModel.UpdatedAt),
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

// BatchDeleteItems implements the ItemDelivery interface.
func (d *ItemDelivery) BatchDeleteItems(ctx context.Context, req *item.BatchDeleteItemsRequest) (resp *item.BatchDeleteItemsResponse, err error) {
	// TODO: Your code here...
	return
}

// GetItemsByTags implements the ItemDelivery interface.
func (d *ItemDelivery) GetItemsByTags(ctx context.Context, req *item.GetItemsByTagsRequest) (resp *item.GetItemsByTagsResponse, err error) {
	// TODO: Your code here...
	return
}

// ExtractMetadata implements the ItemDelivery interface.
func (d *ItemDelivery) ExtractMetadata(ctx context.Context, req *item.ExtractMetadataRequest) (resp *item.ExtractMetadataResponse, err error) {
	// TODO: Your code here...
	return
}

// GetRecentItems implements the ItemDelivery interface.
func (d *ItemDelivery) GetRecentItems(ctx context.Context, req *item.GetRecentItemsRequest) (resp *item.GetRecentItemsResponse, err error) {
	// TODO: Your code here...
	return
}

// BatchUpdateItems implements the ItemDelivery interface.
func (d *ItemDelivery) BatchUpdateItems(ctx context.Context, req *item.BatchUpdateItemsRequest) (resp *item.BatchUpdateItemsResponse, err error) {
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
	// TODO: Your code here...
	return
}

// AddItemNote implements the ItemDelivery interface.
func (d *ItemDelivery) AddItemNote(ctx context.Context, req *item.AddItemNoteRequest) (resp *item.AddItemNoteResponse, err error) {
	// TODO: Your code here...
	return
}

// UpdateItemNote implements the ItemDelivery interface.
func (d *ItemDelivery) UpdateItemNote(ctx context.Context, req *item.UpdateItemNoteRequest) (resp *item.UpdateItemNoteResponse, err error) {
	// TODO: Your code here...
	return
}

// GetItemNote implements the ItemDelivery interface.
func (d *ItemDelivery) GetItemNote(ctx context.Context, req *item.GetItemNoteRequest) (resp *item.GetItemNoteResponse, err error) {
	// TODO: Your code here...
	return
}
