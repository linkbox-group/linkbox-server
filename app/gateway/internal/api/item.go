package api

import (
	"context"
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

	domain.Success(c, resp)
}
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
	domain.Success(c, resp)
}
