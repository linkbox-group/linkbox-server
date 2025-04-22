package api

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/linkbox-group/linkbox-server/gateway/internal/domain"
	"github.com/linkbox-group/linkbox-server/gateway/internal/infra/rpc"
	"github.com/linkbox-group/linkbox-server/rpc-gen/common/pagination"
	"github.com/linkbox-group/linkbox-server/rpc-gen/item"
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
type ItemResponse struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	URL             string    `json:"url"`
	Title           string    `json:"title"`
	ThumbnailURL    string    `json:"thumbnail_url"`
	Tags            []string  `json:"tags"`
	OrganizationIDs []string  `json:"organization_ids"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type ItemListResponse struct {
	Items      []ItemResponse `json:"items"`
	Total      int32          `json:"total"`
	Page       int32          `json:"page"`
	PageSize   int32          `json:"page_size"`
	TotalPages int32          `json:"total_pages"`
}

type ItemSuccessResponse struct {
	Success bool `json:"success"`
}

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
		domain.ErrorMsg(c, ErrAuthFailedCode, err.Error())
		return
	}
	req.UserId = userId

	resp, err := rpc.ItemClient.CreateItem(context.Background(), &req)
	if err != nil {
		domain.Error(c, ErrTitleExists, err.Error())
		return
	}

	item := resp.GetItem()
	itemResp := ItemResponse{
		ID:              item.Id,
		UserID:          item.UserId,
		URL:             item.Url,
		Title:           item.Title,
		ThumbnailURL:    item.ThumbnailUrl,
		Tags:            item.Tags,
		OrganizationIDs: item.CollectionIds,
		CreatedAt:       item.CreatedAt.AsTime(),
		UpdatedAt:       item.UpdatedAt.AsTime(),
	}
	domain.Success(c, itemResp)
}

// GetItem 获取内容
func (a *ItemAPI) GetItem(c *gin.Context) {
	itemID := c.Param("id")
	userId, err := domain.GetUserIdFromContext(c)
	if err != nil {
		domain.ErrorMsg(c, ErrAuthFailedCode, err.Error())
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

	item := resp.GetItem()
	itemResp := ItemResponse{
		ID:              item.Id,
		UserID:          item.UserId,
		URL:             item.Url,
		Title:           item.Title,
		ThumbnailURL:    item.ThumbnailUrl,
		Tags:            item.Tags,
		OrganizationIDs: item.CollectionIds,
		CreatedAt:       item.CreatedAt.AsTime(),
		UpdatedAt:       item.UpdatedAt.AsTime(),
	}
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
		domain.ErrorMsg(c, ErrAuthFailedCode, err.Error())
		return
	}

	req.Id = itemID
	req.UserId = userId
	resp, err := rpc.ItemClient.UpdateItem(context.Background(), &req)
	if err != nil {
		domain.Error(c, ErrTitleExists, err.Error())
		return
	}

	item := resp.GetItem()
	itemResp := ItemResponse{
		ID:              item.Id,
		UserID:          item.UserId,
		URL:             item.Url,
		Title:           item.Title,
		ThumbnailURL:    item.ThumbnailUrl,
		Tags:            item.Tags,
		OrganizationIDs: item.CollectionIds,
		CreatedAt:       item.CreatedAt.AsTime(),
		UpdatedAt:       item.UpdatedAt.AsTime(),
	}
	domain.Success(c, itemResp)
}

// DeleteItem 删除内容
func (a *ItemAPI) DeleteItem(c *gin.Context) {
	itemID := c.Param("id")
	userId, err := domain.GetUserIdFromContext(c)
	if err != nil {
		domain.ErrorMsg(c, ErrAuthFailedCode, err.Error())
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

	success := resp.GetSuccess()
	itemResp := ItemSuccessResponse{
		Success: success,
	}
	domain.Success(c, itemResp)
}

// GetItemsByTags 按标签获取内容
func (a *ItemAPI) GetItemsByTags(c *gin.Context) {
	userId, err := domain.GetUserIdFromContext(c)
	if err != nil {
		domain.ErrorMsg(c, ErrAuthFailedCode, err.Error())
		return
	}

	req := item.GetItemsByTagsRequest{
		Pagination: &pagination.PaginationRequest{},
	}
	req.UserId = userId
	err = c.ShouldBind(&req)
	if err != nil {
		logrus.Infoln(err)
		domain.ErrorMsg(c, ErrInvalidParamCode, err.Error())
		return
	}
	logrus.Infoln(req.Pagination)
	resp, err := rpc.ItemClient.GetItemsByTags(context.Background(), &req)
	if err != nil {
		domain.Error(c, ErrRpcFailedCode, "rpc调用失败")
		return
	}

	var items []ItemResponse
	for _, ite := range resp.GetItemsPage().Items {
		items = append(items, ItemResponse{
			ID:              ite.Id,
			UserID:          ite.UserId,
			URL:             ite.Url,
			Title:           ite.Title,
			ThumbnailURL:    ite.ThumbnailUrl,
			Tags:            ite.Tags,
			OrganizationIDs: ite.CollectionIds,
			CreatedAt:       ite.CreatedAt.AsTime(),
			UpdatedAt:       ite.UpdatedAt.AsTime(),
		})
	}

	itemListResp := ItemListResponse{
		Items:      items,
		Total:      resp.GetItemsPage().Pagination.TotalItems,
		Page:       resp.GetItemsPage().Pagination.Page,
		PageSize:   resp.GetItemsPage().Pagination.PageSize,
		TotalPages: resp.GetItemsPage().Pagination.TotalPages,
	}
	domain.Success(c, itemListResp)
}

// GetItemsByOrganization 按组织获取内容
func (a *ItemAPI) GetItemsByOrganization(c *gin.Context) {
	userId, err := domain.GetUserIdFromContext(c)
	if err != nil {
		domain.ErrorMsg(c, ErrAuthFailedCode, err.Error())
		return
	}
	req := item.GetItemsByOrganizationRequest{
		UserId: userId,
	}
	err = c.ShouldBind(&req)
	if err != nil {
		domain.Error(c, ErrInvalidReq, err.Error())
	}
	logrus.Infoln(req.OrganizationId)
	resp, err := rpc.ItemClient.GetItemsByOrganization(context.Background(), &req)
	if err != nil {
		domain.Error(c, ErrNoPermission, "没有操作权限")
		return
	}

	var items []ItemResponse
	for _, ite := range resp.GetItemsPage().Items {
		items = append(items, ItemResponse{
			ID:              ite.Id,
			UserID:          ite.UserId,
			URL:             ite.Url,
			Title:           ite.Title,
			ThumbnailURL:    ite.ThumbnailUrl,
			Tags:            ite.Tags,
			OrganizationIDs: ite.CollectionIds,
			CreatedAt:       ite.CreatedAt.AsTime(),
			UpdatedAt:       ite.UpdatedAt.AsTime(),
		})
	}

	itemListResp := ItemListResponse{
		Items:      items,
		Total:      resp.GetItemsPage().Pagination.TotalItems,
		Page:       resp.GetItemsPage().Pagination.Page,
		PageSize:   resp.GetItemsPage().Pagination.PageSize,
		TotalPages: resp.GetItemsPage().Pagination.TotalPages,
	}
	domain.Success(c, itemListResp)
}
func (a *ItemAPI) SearchItems(c *gin.Context) {
	userId, err := domain.GetUserIdFromContext(c)
	if err!= nil {
		domain.ErrorMsg(c, ErrAuthFailedCode, err.Error())
		return
	}
	req := item.SearchItemsRequest{
		Pagination: &pagination.PaginationRequest{},
	}
	req.UserId = userId
	err = c.ShouldBind(&req)
	if err!= nil {
		domain.Error(c, ErrInvalidReq, err.Error())
	}
	resp, err := rpc.ItemClient.SearchItems(context.Background(), &req)
	if err!= nil {
		domain.Error(c, ErrRpcFailedCode, "rpc调用失败")
		return
	}
	var items []ItemResponse
	for _, ite := range resp.GetData().GetItems() {
		items = append(items, ItemResponse{
			ID:              ite.Id,
			UserID:          ite.UserId,
			URL:             ite.Url,
			Title:           ite.Title,
			ThumbnailURL:    ite.ThumbnailUrl,
			Tags:            ite.Tags,
			OrganizationIDs: ite.CollectionIds,
			CreatedAt:       ite.CreatedAt.AsTime(),
			UpdatedAt:       ite.UpdatedAt.AsTime(),	
		})
	}
	itemListResp := ItemListResponse{
		Items:      items,
		Total:      resp.GetData().GetPagination().TotalItems,
		Page:       resp.GetData().GetPagination().Page,
		PageSize:   resp.GetData().GetPagination().PageSize,
		TotalPages: resp.GetData().GetPagination().TotalPages,	
	}
	domain.Success(c, itemListResp)
}