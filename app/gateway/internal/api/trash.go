package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/linkbox-group/linkbox-server/common/ecode"
	"github.com/linkbox-group/linkbox-server/gateway/internal/domain"
	"github.com/linkbox-group/linkbox-server/gateway/internal/infra/rpc"
	"github.com/linkbox-group/linkbox-server/gateway/pkg/log"
	"github.com/linkbox-group/linkbox-server/rpc-gen/common/pagination"
	"github.com/linkbox-group/linkbox-server/rpc-gen/item"
	"strconv"
)

type TrashApi struct {
}

func (a *TrashApi) ListTrashItem(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	userID, err := domain.GetUserIDFromContext(c)
	if err != nil {
		domain.Error(c, ecode.ErrAuthFailed, "用户认证失败")
		return
	}
	log.Log().Info("[a *TrashAPI] ListTrashItem ", "user_id", userID)
	req := &item.GetDeletedItemsRequest{
		UserId: userID,
		Pagination: &pagination.PaginationRequest{
			Page:     int32(page),
			PageSize: int32(pageSize),
		},
	}
	resp, err := rpc.ItemClient.GetDeletedItems(context.Background(), req)
	if err != nil {
		domain.Error(c, ErrTagNotFound, "获取回收站数据失败")
		return
	}

	itemRes := resp.GetItemsPage()
	trashItems := make([]*domain.TrashItem, 0)
	for _, i := range itemRes.Items {
		trashItem := &domain.TrashItem{}
		trashItem.Convert(i)
		trashItems = append(trashItems, trashItem)
	}
	res := domain.ListTrashItemResp{
		Items:      trashItems,
		Page:       itemRes.GetPagination().Page,
		PageSize:   itemRes.GetPagination().PageSize,
		TotalPages: itemRes.GetPagination().TotalPages,
		Total:      itemRes.GetPagination().TotalItems,
	}
	domain.Success(c, res)
}
func (a *TrashApi) RecoveryTrashItem(c *gin.Context) {
	userID, err := domain.GetUserIDFromContext(c)
	if err != nil {
		domain.Error(c, ecode.ErrAuthFailed, "用户认证失败")
		return
	}
	log.Log().Info("[a *TrashAPI] RecoveryTrashItem ", "user_id", userID)
	req := domain.RecoveryTrashItemReq{}
	err = c.ShouldBind(&req)
	if err != nil {
		domain.Error(c, ecode.ErrInvalidParam, "请求参数错误")
		return
	}

	resp, err := rpc.ItemClient.RecoverItemsBatch(context.Background(), &item.RecoverItemsBatchRequest{
		UserId: userID,
		Ids:    []string{req.ItemID},
	})
	if err != nil {
		domain.Error(c, ErrTagNotFound, "恢复数据失败")
		return
	}

	success := resp.GetSuccess()
	domain.Success(c, domain.RecoveryTrashItemResp{
		Success: success,
	})
}
func (a *TrashApi) DeleteTrashItem(c *gin.Context) {
	itemID := c.Param("item_id")
	userID, err := domain.GetUserIDFromContext(c)
	if err != nil {
		domain.Error(c, ecode.ErrAuthFailed, "用户认证失败")
		return
	}
	log.Log().Info("[a *TrashApi] DeleteTrashItem ", "user_id", userID)

	resp, err := rpc.ItemClient.DeleteItemsBatch(context.Background(), &item.DeleteItemsBatchRequest{
		UserId: userID,
		Ids:    []string{itemID},
	})
	if err != nil {
		domain.Error(c, ErrTagNotFound, "删除标签失败")
		return
	}

	success := resp.GetSuccess()
	tagResp := domain.TagSuccessResponse{
		Success: success,
	}
	domain.Success(c, tagResp)
}
