package api

import (
	"context"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/linkbox-group/linkbox-server/gateway/internal/domain"
	"github.com/linkbox-group/linkbox-server/gateway/internal/infra/rpc"
	"github.com/linkbox-group/linkbox-server/rpc-gen/common/pagination"
	"github.com/linkbox-group/linkbox-server/rpc-gen/item"
)

type ItemAPI struct{}

func NewItemAPI() *ItemAPI {
	return &ItemAPI{}
}

// 错误码定义
const (
	ErrTitleExists = 40001 // 标题已存在
)

// CreateItem 创建内容
func (a *ItemAPI) CreateItem(c *gin.Context) {
	var req item.CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
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

	domain.Success(c, resp)
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

	domain.Success(c, resp)
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

	domain.Success(c, resp)
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

	domain.Success(c, resp)
}

// GetItemsByTags 按标签获取内容
func (a *ItemAPI) GetItemsByTags(c *gin.Context) {
	userId, err := domain.GetUserIdFromContext(c)
	if err != nil {
		domain.ErrorMsg(c, ErrAuthFailedCode, err.Error())
		return
	}

	tagsStr := c.Query("tags")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if tagsStr == "" {
		domain.Error(c, ErrInvalidReq, "请求参数错误")
		return
	}

	tags := strings.Split(tagsStr, ",")
	resp, err := rpc.ItemClient.GetItemsByTags(context.Background(), &item.GetItemsByTagsRequest{
		UserId: userId,
		Tags:   tags,
		Pagination: &pagination.PaginationRequest{
			Page:     int32(page),
			PageSize: int32(pageSize),
		},
	})
	if err != nil {
		domain.Error(c, ErrNoPermission, "没有操作权限")
		return
	}

	domain.Success(c, resp)
}
