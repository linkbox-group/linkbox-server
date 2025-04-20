package api

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/linkbox-group/linkbox-server/gateway/internal/domain"
	"github.com/linkbox-group/linkbox-server/gateway/internal/infra/rpc"
	"github.com/linkbox-group/linkbox-server/rpc-gen/organization"
	"github.com/sirupsen/logrus"
)

type OrganizationAPI struct{}

func NewOrganizationAPI() *OrganizationAPI {
	return &OrganizationAPI{}
}

// 错误码定义
const (
	ErrOrganizationNotFound   = 40000 // 组织不存在
	ErrOrganizationNameExists = 40001 // 组织名已存在
)

// CreateOrganization 创建组织
func (a *OrganizationAPI) CreateOrganization(c *gin.Context) {
	var req organization.CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Infoln(err)
		domain.Error(c, ErrInvalidReq, "请求参数错误")
		return
	}

	resp, err := rpc.OrganizationClient.CreateOrganization(context.Background(), &req)
	if err != nil {
		domain.Error(c, ErrOrganizationNameExists, err.Error())
		return
	}

	domain.Success(c, resp)
}

// GetOrganization 获取组织详情
func (a *OrganizationAPI) GetOrganization(c *gin.Context) {
	orgID := c.Param("id")
	userId, err := domain.GetUserIdFromContext(c)
	if err != nil {
		domain.Error(c, ErrAuthFailedCode, err.Error())
		return
	}

	resp, err := rpc.OrganizationClient.GetOrganization(context.Background(), &organization.GetOrganizationRequest{
		Id:     orgID,
		UserId: userId,
	})
	if err != nil {
		domain.Error(c, ErrOrganizationNotFound, "组织不存在")
		return
	}

	domain.Success(c, resp)
}

// UpdateOrganization 更新组织
func (a *OrganizationAPI) UpdateOrganization(c *gin.Context) {
	var req organization.UpdateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		domain.Error(c, ErrInvalidReq, "请求参数错误")
		return
	}

	resp, err := rpc.OrganizationClient.UpdateOrganization(context.Background(), &req)
	if err != nil {
		domain.Error(c, ErrOrganizationNameExists, err.Error())
		return
	}

	domain.Success(c, resp)
}

// DeleteOrganization 删除组织
func (a *OrganizationAPI) DeleteOrganization(c *gin.Context) {
	orgID := c.Param("id")
	userID := c.Query("user_id")
	cascade := c.DefaultQuery("cascade", "false")

	if userID == "" {
		domain.Error(c, ErrInvalidReq, "请求参数错误")
		return
	}

	resp, err := rpc.OrganizationClient.DeleteOrganization(context.Background(), &organization.DeleteOrganizationRequest{
		Id:      orgID,
		UserId:  userID,
		Cascade: cascade == "true",
	})
	if err != nil {
		domain.Error(c, ErrOrganizationNotFound, "组织不存在")
		return
	}

	domain.Success(c, resp)
}

// GetUserOrganizations 获取用户组织列表
func (a *OrganizationAPI) GetUserOrganizations(c *gin.Context) {
	userId, err := domain.GetUserIdFromContext(c)
	if err != nil {
		domain.Error(c, ErrAuthFailedCode, err.Error())
		return
	}

	resp, err := rpc.OrganizationClient.GetUserOrganizations(context.Background(), &organization.GetUserOrganizationsRequest{
		UserId: userId,
	})
	if err != nil {
		domain.Error(c, ErrNoPermission, "没有操作权限")
		return
	}

	domain.Success(c, resp)
}

// GetOrganizationTree 获取组织树
func (a *OrganizationAPI) GetOrganizationTree(c *gin.Context) {
	userId, err := domain.GetUserIdFromContext(c)
	if err != nil {
		domain.Error(c, ErrAuthFailedCode, err.Error())
		return
	}
	rootCode := c.Query("root_code")

	resp, err := rpc.OrganizationClient.GetOrganizationTree(context.Background(), &organization.GetOrganizationTreeRequest{
		UserId:   userId,
		RootCode: &rootCode,
	})
	if err != nil {
		domain.Error(c, ErrOrganizationNotFound, "组织不存在")
		return
	}

	domain.Success(c, resp)
}

// GetOrganizationChildren 获取组织子节点
func (a *OrganizationAPI) GetOrganizationChildren(c *gin.Context) {
	userId, err := domain.GetUserIdFromContext(c)
	if err != nil {
		domain.Error(c, ErrAuthFailedCode, err.Error())
		return
	}
	parentCode := c.Query("parent_code")
	recursive := c.DefaultQuery("recursive", "false")

	if parentCode == "" {
		domain.Error(c, ErrInvalidReq, "请求参数错误")
		return
	}

	resp, err := rpc.OrganizationClient.GetOrganizationChildren(context.Background(), &organization.GetOrganizationChildrenRequest{
		UserId:     userId,
		ParentCode: parentCode,
		Recursive:  recursive == "true",
	})
	if err != nil {
		domain.Error(c, ErrOrganizationNotFound, "组织不存在")
		return
	}

	domain.Success(c, resp)
}

// MoveOrganization 移动组织
func (a *OrganizationAPI) MoveOrganization(c *gin.Context) {
	var req organization.MoveOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		domain.Error(c, ErrInvalidReq, "请求参数错误")
		return
	}

	resp, err := rpc.OrganizationClient.MoveOrganization(context.Background(), &req)
	if err != nil {
		domain.Error(c, ErrOrganizationNotFound, "组织不存在")
		return
	}

	domain.Success(c, resp)
}

// AddItemsToOrganization 添加内容项到组织
func (a *OrganizationAPI) AddItemsToOrganization(c *gin.Context) {
	var req organization.AddItemsToOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		domain.Error(c, ErrInvalidReq, "请求参数错误")
		return
	}
	userId, err := domain.GetUserIdFromContext(c)
	if err != nil {
		domain.Error(c, ErrAuthFailedCode, err.Error())
		return
	}
	req.UserId = userId

	resp, err := rpc.OrganizationClient.AddItemsToOrganization(context.Background(), &req)
	if err != nil {
		domain.Error(c, ErrItemNotFound, "内容项不存在")
		return
	}

	domain.Success(c, resp)
}

// RemoveItemsFromOrganization 从组织移除内容项
func (a *OrganizationAPI) RemoveItemsFromOrganization(c *gin.Context) {
	var req organization.RemoveItemsFromOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		domain.Error(c, ErrInvalidReq, "请求参数错误")
		return
	}
	userId, err := domain.GetUserIdFromContext(c)
	if err != nil {
		domain.Error(c, ErrAuthFailedCode, err.Error())
		return
	}
	req.UserId = userId
	resp, err := rpc.OrganizationClient.RemoveItemsFromOrganization(context.Background(), &req)
	if err != nil {
		domain.Error(c, ErrItemNotFound, "内容项不存在")
		return
	}

	domain.Success(c, resp)
}
