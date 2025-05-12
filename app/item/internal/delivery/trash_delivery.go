package delivery

import (
	"context"
	"github.com/linkbox-group/linkbox-server/item/pkg/log"
	"github.com/linkbox-group/linkbox-server/rpc-gen/common/cError"
	"github.com/linkbox-group/linkbox-server/rpc-gen/common/pagination"
	"github.com/linkbox-group/linkbox-server/rpc-gen/item"
	"github.com/linkbox-group/linkbox-server/rpc-gen/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (d *ItemDelivery) GetDeletedItems(ctx context.Context, req *item.GetDeletedItemsRequest) (res *item.GetDeletedItemsResponse, err error) {

	items, total, err := d.s.GetDeletedItems(ctx, req.GetUserId(), req.GetPagination())

	if err != nil {
		log.Log().Error("get deleted items  failed" + err.Error())
		return &item.GetDeletedItemsResponse{
			Result: &item.GetDeletedItemsResponse_Error{
				Error: &cError.Error{
					Message: err.Error(),
				},
			},
		}, err
	}
	trashItems := make([]*model.TrashItem, 0)
	for _, i := range items {
		trashItems = append(trashItems, &model.TrashItem{
			Id:        i.ID,
			UserId:    i.UserID,
			Title:     i.Title,
			DeletedAt: timestamppb.New(i.DeletedAt.Time),
			ExpiredAt: timestamppb.New(i.DeletedAt.Time.AddDate(0, 0, 30)),
		})
	}

	return &item.GetDeletedItemsResponse{
		Result: &item.GetDeletedItemsResponse_ItemsPage{
			ItemsPage: &item.TrashItemsPage{
				Items: trashItems,
				Pagination: &pagination.PaginationMeta{
					TotalPages: int32(total) / req.GetPagination().GetPageSize(),
					TotalItems: int32(total),
					PageSize:   req.GetPagination().GetPageSize(),
					Page:       req.GetPagination().GetPage(),
				},
			},
		},
	}, nil
}

func (d *ItemDelivery) RecoverItemsBatch(ctx context.Context, req *item.RecoverItemsBatchRequest) (res *item.RecoverItemsBatchResponse, err error) {
	userID := req.UserId

	err = d.s.RecoverItemsBatch(ctx, userID, req.Ids)

	if err != nil {
		log.Log().Error("recover items batch failed" + err.Error())
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
	userID := req.UserId

	err = d.s.DeleteItemsBatch(ctx, userID, req.Ids)

	if err != nil {
		log.Log().Error("delete items batch failed" + err.Error())
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
