package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/linkbox-group/linkbox-server/common/ecode"
	"github.com/linkbox-group/linkbox-server/gateway/internal/domain"
	"github.com/linkbox-group/linkbox-server/gateway/internal/infra/rpc"
	"github.com/linkbox-group/linkbox-server/model/treemodel"
	"github.com/linkbox-group/linkbox-server/rpc-gen/common/pagination"
	"github.com/linkbox-group/linkbox-server/rpc-gen/item"
	itemmodel "github.com/linkbox-group/linkbox-server/rpc-gen/model"
	"github.com/linkbox-group/linkbox-server/rpc-gen/organization"
	"github.com/sirupsen/logrus"
)

type ItemAPI struct{}

func NewItemAPI() *ItemAPI {
	return &ItemAPI{}
}

// 错误码定义
const (
	ErrTitleExists = 40001 // 标题已存在
)

// 响应结构体定义
// 已移动到domain/item.go中

// CreateItem 创建内容
func (a *ItemAPI) CreateItem(c *gin.Context) {
	var req item.CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Infoln(err)
		domain.Error(c, ErrInvalidReq, "请求参数错误")
		return
	}

	userId, err := domain.GetUserIdFromContext(c)
	if err != nil {
		domain.ErrorMsg(c, ecode.ErrAuthFailed, err.Error())
		return
	}
	req.UserId = userId

	resp, err := rpc.ItemClient.CreateItem(context.Background(), &req)
	if err != nil {
		domain.Error(c, ErrTitleExists, err.Error())
		return
	}

	itemData := resp.GetItem()
	itemResp := &domain.Item{}
	itemResp.Convert(itemData)
	domain.Success(c, itemResp)
}

// GetItem 获取内容
func (a *ItemAPI) GetItem(c *gin.Context) {
	itemID := c.Param("id")
	userId, err := domain.GetUserIdFromContext(c)
	if err != nil {
		domain.ErrorMsg(c, ecode.ErrAuthFailed, err.Error())
		return
	}

	resp, err := rpc.ItemClient.GetItem(context.Background(), &item.GetItemRequest{
		Id:     itemID,
		UserId: userId,
	})
	if err != nil {
		domain.Error(c, ErrItemNotFound, "内容不存在")
		return
	}

	itemData := resp.GetItem()
	itemResp := &domain.Item{}
	itemResp.Convert(itemData)
	domain.Success(c, itemResp)
}

// UpdateItem 更新内容
func (a *ItemAPI) UpdateItem(c *gin.Context) {
	itemID := c.Param("id")
	var req item.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		domain.Error(c, ErrInvalidReq, "请求参数错误")
		return
	}

	userId, err := domain.GetUserIdFromContext(c)
	if err != nil {
		domain.ErrorMsg(c, ecode.ErrAuthFailed, err.Error())
		return
	}

	req.Id = itemID
	req.UserId = userId
	resp, err := rpc.ItemClient.UpdateItem(context.Background(), &req)
	if err != nil {
		domain.Error(c, ErrTitleExists, err.Error())
		return
	}

	itemData := resp.GetItem()
	itemResp := &domain.Item{}
	itemResp.Convert(itemData)
	domain.Success(c, itemResp)
}

// DeleteItem 删除内容
func (a *ItemAPI) DeleteItem(c *gin.Context) {
	itemID := c.Param("id")
	userId, err := domain.GetUserIdFromContext(c)
	if err != nil {
		domain.ErrorMsg(c, ecode.ErrAuthFailed, err.Error())
		return
	}

	resp, err := rpc.ItemClient.DeleteItem(context.Background(), &item.DeleteItemRequest{
		Id:     itemID,
		UserId: userId,
	})
	if err != nil {
		domain.Error(c, ErrItemNotFound, "内容不存在")
		return
	}
	itemResp := domain.ItemSuccessResponse{
		Success: resp.GetSuccess(),
	}
	domain.Success(c, itemResp)
}

// GetItemsByTags 按标签获取内容
func (a *ItemAPI) GetItemsByTags(c *gin.Context) {
	userId, err := domain.GetUserIdFromContext(c)
	if err != nil {
		domain.ErrorMsg(c, ecode.ErrAuthFailed, err.Error())
		return
	}

	req := item.GetItemsByTagsRequest{
		Pagination: &pagination.PaginationRequest{},
	}
	req.UserId = userId
	err = c.ShouldBind(&req)
	if err != nil {
		logrus.Infoln(err)
		domain.ErrorMsg(c, ecode.ErrInvalidParam, err.Error())
		return
	}
	logrus.Infoln(req.Pagination)
	resp, err := rpc.ItemClient.GetItemsByTags(context.Background(), &req)
	if err != nil {
		domain.Error(c, ecode.ErrRpcServiceError, "rpc调用失败")
		return
	}

	var items []*domain.Item
	for _, ite := range resp.GetItemsPage().Items {
		itemResp := &domain.Item{}
		itemResp.Convert(ite)
		items = append(items, itemResp)
	}

	itemListResp := domain.ItemListResponse{
		Items: items,
		Pagination: domain.Pagination{
			Total:      resp.GetItemsPage().Pagination.TotalItems,
			Page:       resp.GetItemsPage().Pagination.Page,
			PageSize:   resp.GetItemsPage().Pagination.PageSize,
			TotalPages: resp.GetItemsPage().Pagination.TotalPages,
		},
	}
	domain.Success(c, itemListResp)
}

// GetItemsByOrganization 按组织获取内容
func (a *ItemAPI) GetItemsByOrganization(c *gin.Context) {
	req := &item.GetItemsByOrganizationRequest{}
	err := c.ShouldBind(req)
	if err != nil {
		domain.Error(c, ErrInvalidReq, err.Error())
	}
	userId, err := domain.GetUserIdFromContext(c)
	if err != nil {
		domain.ErrorMsg(c, ecode.ErrAuthFailed, err.Error())
		return
	}
	req.UserId = userId
	if req.OrganizationId == "" || req.OrganizationId == treemodel.ROOT_ID {
		orgID, err := rpc.OrganizationClient.GetDefaultOrgID(c, &organization.GetDefaultOrgIDReq{
			UserId: req.UserId,
			Code:   treemodel.ROOT_ID,
		})

		if err != nil {
			domain.ErrorMsg(c, ecode.ErrRpcServiceError, err.Error())
			return
		}
		req.OrganizationId = orgID.GetId()
	}

	resp, err := rpc.ItemClient.GetItemsByOrganization(context.Background(), req)
	if err != nil {
		domain.Error(c, ecode.ErrRpcServiceError, err.Error())
		return
	}

	var items []*domain.Item
	for _, ite := range resp.GetItemsPage().Items {
		itemResp := &domain.Item{}
		itemResp.Convert(ite)
		items = append(items, itemResp)
	}

	itemListResp := domain.ItemListResponse{
		Items: items,
		Pagination: domain.Pagination{
			Total:      resp.GetItemsPage().Pagination.TotalItems,
			Page:       resp.GetItemsPage().Pagination.Page,
			PageSize:   resp.GetItemsPage().Pagination.PageSize,
			TotalPages: resp.GetItemsPage().Pagination.TotalPages,
		},
	}
	domain.Success(c, itemListResp)
}
func (a *ItemAPI) SearchItems(c *gin.Context) {
	userId, err := domain.GetUserIdFromContext(c)
	if err != nil {
		domain.ErrorMsg(c, ecode.ErrAuthFailed, err.Error())
		return
	}
	req := domain.SearchItemsReq{}
	err = c.ShouldBind(&req)

	if err != nil {
		domain.Error(c, ErrInvalidReq, err.Error())
	}
	resp, err := rpc.ItemClient.SearchItems(context.Background(), &item.SearchItemsRequest{
		UserId:     userId,
		Query:      req.Query,
		Pagination: req.Pagination,
		Type:       itemmodel.ItemType(itemmodel.ItemType_value[req.ItemType]),
	})
	if err != nil {
		domain.Error(c, ecode.ErrRpcServiceError, "rpc调用失败")
		return
	}
	var items []*domain.Item
	for _, ite := range resp.GetData().Items {
		itemResp := &domain.Item{}
		itemResp.Convert(ite)

		items = append(items, itemResp)
	}

	itemListResp := domain.ItemListResponse{

		Items: items,
		Pagination: domain.Pagination{
			Total:      resp.GetData().GetPagination().TotalItems,
			Page:       resp.GetData().GetPagination().Page,
			PageSize:   resp.GetData().GetPagination().PageSize,
			TotalPages: resp.GetData().GetPagination().TotalPages,
		},
	}
	domain.Success(c, itemListResp)
}
func (a *ItemAPI) RecoverItemBatch(c *gin.Context) {
	var reqBody struct {
		Ids []string `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		domain.Error(c, ErrInvalidReq, "请求参数错误")
		return
	}

	userId, err := domain.GetUserIdFromContext(c)
	if err != nil {
		domain.ErrorMsg(c, ecode.ErrAuthFailed, err.Error())
		return
	}

	resp, err := rpc.ItemClient.RecoverItemsBatch(context.Background(), &item.RecoverItemsBatchRequest{
		Ids:    reqBody.Ids,
		UserId: userId,
	})
	if err != nil {
		domain.Error(c, ErrItemNotFound, "内容不存在")
		return
	}
	itemResp := domain.ItemSuccessResponse{
		Success: resp.GetSuccess(),
	}
	domain.Success(c, itemResp)
}
func (a *ItemAPI) DeleteItemBatch(c *gin.Context) {
	var reqBody struct {
		Ids []string `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		domain.Error(c, ErrInvalidReq, "请求参数错误")
		return
	}

	userId, err := domain.GetUserIdFromContext(c)
	if err != nil {
		domain.ErrorMsg(c, ecode.ErrAuthFailed, err.Error())
		return
	}

	resp, err := rpc.ItemClient.DeleteItemsBatch(context.Background(), &item.DeleteItemsBatchRequest{
		Ids:    reqBody.Ids,
		UserId: userId,
	})
	if err != nil {
		domain.Error(c, ErrItemNotFound, "内容不存在")
		return
	}
	itemResp := domain.ItemSuccessResponse{
		Success: resp.GetSuccess(),
	}
	domain.Success(c, itemResp)
}
