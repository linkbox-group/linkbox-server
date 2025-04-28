package api

import (
	"context"
	"github.com/linkbox-group/linkbox-server/common/ecode"
	"github.com/linkbox-group/linkbox-server/model/treemodel"
	"github.com/linkbox-group/linkbox-server/rpc-gen/item"

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
	var req domain.CreateOrganizationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Error(err)
		domain.Error(c, ErrInvalidReq, "请求参数错误")
		return
	}
	userId, err := domain.GetUserIdFromContext(c)
	if err != nil {
		domain.Error(c, ecode.ErrAuthFailed, err.Error())
		return
	}
	orgReq := organization.CreateOrganizationRequest{
		UserId:      userId,
		Name:        req.Name,
		ParentCode:  &req.ParentCode,
		Description: &req.Description,
		SortOrder:   &req.SortOrder,
	}
	resp, err := rpc.OrganizationClient.CreateOrganization(context.Background(), &orgReq)
	if err != nil {
		domain.Error(c, ErrOrganizationNameExists, err.Error())
		return
	}
	org := resp.GetOrganization()
	orgResp := domain.Organization{
		ID:          org.Id,
		Code:        org.Code,
		ParentCode:  org.ParentCode,
		Name:        org.Name,
		UserID:      org.UserId,
		Description: org.Description,
		SortOrder:   int(org.SortOrder),
		CreatedAt:   org.CreatedAt.AsTime(),
		UpdatedAt:   org.UpdatedAt.AsTime(),
	}
	domain.Success(c, orgResp)
}

// GetOrganization 获取组织详情
func (a *OrganizationAPI) GetOrganization(c *gin.Context) {
	orgID := c.Param("id")
	userId, err := domain.GetUserIdFromContext(c)
	if err != nil {
		domain.Error(c, ecode.ErrAuthFailed, err.Error())
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

	org := resp.GetOrganization()
	orgResp := domain.Organization{
		ID:          org.Id,
		Code:        org.Code,
		ParentCode:  org.ParentCode,
		Name:        org.Name,
		UserID:      org.UserId,
		Description: org.Description,
		SortOrder:   int(org.SortOrder),
		CreatedAt:   org.CreatedAt.AsTime(),
		UpdatedAt:   org.UpdatedAt.AsTime(),
	}
	domain.Success(c, orgResp)
}

// UpdateOrganization 更新组织
func (a *OrganizationAPI) UpdateOrganization(c *gin.Context) {
	var req organization.UpdateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		domain.Error(c, ErrInvalidReq, "请求参数错误")
		return
	}
	userId, err := domain.GetUserIdFromContext(c)
	if err != nil {
		domain.Error(c, ecode.ErrAuthFailed, err.Error())
		return
	}
	req.UserId = userId

	resp, err := rpc.OrganizationClient.UpdateOrganization(context.Background(), &req)
	if err != nil {
		domain.Error(c, ErrOrganizationNameExists, err.Error())
		return
	}

	org := resp.GetOrganization()
	orgResp := domain.Organization{
		ID:          org.Id,
		Code:        org.Code,
		ParentCode:  org.ParentCode,
		Name:        org.Name,
		UserID:      org.UserId,
		Description: org.Description,
		SortOrder:   int(org.SortOrder),
		CreatedAt:   org.CreatedAt.AsTime(),
		UpdatedAt:   org.UpdatedAt.AsTime(),
	}
	domain.Success(c, orgResp)
}

// DeleteOrganization 删除组织
func (a *OrganizationAPI) DeleteOrganization(c *gin.Context) {
	orgID := c.Param("id")
	cascade := c.DefaultQuery("cascade", "false")
	userId, err := domain.GetUserIdFromContext(c)
	if err != nil {
		domain.Error(c, ecode.ErrAuthFailed, err.Error())
		return
	}
	req := &organization.DeleteOrganizationRequest{
		UserId:  userId,
		Id:      orgID,
		Cascade: cascade == "ture",
	}
	resp, err := rpc.OrganizationClient.DeleteOrganization(context.Background(), req)
	if err != nil {
		domain.Error(c, ErrOrganizationNotFound, "组织不存在")
		return
	}
	orgResp := domain.OrganizationSuccessResponse{
		Success: resp.GetSuccess(),
	}

	domain.Success(c, orgResp)
}

// GetUserOrganizations 获取用户组织列表
func (a *OrganizationAPI) GetUserOrganizations(c *gin.Context) {
	userId, err := domain.GetUserIdFromContext(c)
	if err != nil {
		domain.Error(c, ecode.ErrAuthFailed, err.Error())
		return
	}

	resp, err := rpc.OrganizationClient.GetUserOrganizations(context.Background(), &organization.GetUserOrganizationsRequest{
		UserId: userId,
	})
	if err != nil {
		domain.Error(c, ecode.ErrRpcServiceError, "rpc服务错误"+err.Error())
		return
	}

	var orgs []*domain.Organization
	for _, org := range resp.GetOrganizations().Organizations {
		orgs = append(orgs, &domain.Organization{
			ID:          org.Id,
			Code:        org.Code,
			ParentCode:  org.ParentCode,
			Name:        org.Name,
			UserID:      org.UserId,
			Description: org.Description,
			SortOrder:   int(org.SortOrder),
			CreatedAt:   org.CreatedAt.AsTime(),
			UpdatedAt:   org.UpdatedAt.AsTime(),
		})
	}
	orgsResp := domain.ListOrganizationsResponse{
		Organizations: orgs,
	}
	domain.Success(c, orgsResp)

}

// GetOrganizationTree 获取组织树
func (a *OrganizationAPI) GetOrganizationTree(c *gin.Context) {
	userId, err := domain.GetUserIdFromContext(c)
	if err != nil {
		domain.Error(c, ecode.ErrAuthFailed, err.Error())
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

	domain.Success(c, resp.GetRoot().GetData())
}

// GetOrganizationChildren 获取组织子节点
func (a *OrganizationAPI) GetOrganizationChildren(c *gin.Context) {
	userId, err := domain.GetUserIdFromContext(c)
	if err != nil {
		domain.Error(c, ecode.ErrAuthFailed, err.Error())
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

	var orgs []*domain.Organization
	for _, org := range resp.GetChildren().Organizations {
		orgs = append(orgs, &domain.Organization{
			ID:          org.Id,
			Code:        org.Code,
			ParentCode:  org.ParentCode,
			Name:        org.Name,
			UserID:      org.UserId,
			Description: org.Description,
			SortOrder:   int(org.SortOrder),
			CreatedAt:   org.CreatedAt.AsTime(),
			UpdatedAt:   org.UpdatedAt.AsTime(),
		})
	}
	orgsResp := domain.ListOrganizationsResponse{
		Organizations: orgs}
	domain.Success(c, orgsResp)

}

// MoveOrganization 移动组织
func (a *OrganizationAPI) MoveOrganization(c *gin.Context) {
	var req organization.MoveOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		domain.Error(c, ErrInvalidReq, "请求参数错误")
		return
	}
	userId, err := domain.GetUserIdFromContext(c)
	if err != nil {
		domain.Error(c, ecode.ErrAuthFailed, err.Error())
		return
	}
	req.UserId = userId
	resp, err := rpc.OrganizationClient.MoveOrganization(context.Background(), &req)
	if err != nil {
		domain.Error(c, ErrOrganizationNotFound, "组织不存在")
		return
	}

	org := resp.GetSuccess()
	orgResp := domain.OrganizationSuccessResponse{
		Success: org,
	}
	domain.Success(c, orgResp)
}

// AddItemsToOrganization 添加内容项到组织
func (a *OrganizationAPI) AddItemsToOrganization(c *gin.Context) {
	var req domain.AddItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		domain.Error(c, ErrInvalidReq, "请求参数错误")
		return
	}
	userId, err := domain.GetUserIdFromContext(c)
	if err != nil {
		domain.Error(c, ecode.ErrAuthFailed, err.Error())
		return
	}
	successCount := 0
	FailureCount := 0
	FailedItemIDs := make([]string, 0)
	for _, itemID := range req.ItemID {
		ItemReq := item.UpdateItemRequest{
			Id:             itemID,
			OrganizationId: req.OrganizationID,
			UserId:         userId,
		}
		_, err = rpc.ItemClient.UpdateItem(c, &ItemReq)
		if err != nil {
			FailureCount++
			FailedItemIDs = append(FailedItemIDs, itemID)
			continue
		}
		successCount++
	}
	orgResp := domain.AddItemsToOrganizationResponse{
		SuccessCount:  int32(successCount),
		FailureCount:  int32(FailureCount),
		FailedItemIDs: FailedItemIDs,
	}
	domain.Success(c, orgResp)
}

// RemoveItemsFromOrganization 从组织移除内容项
func (a *OrganizationAPI) RemoveItemsFromOrganization(c *gin.Context) {
	var req domain.RemoveItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		domain.Error(c, ErrInvalidReq, "请求参数错误")
		return
	}
	userId, err := domain.GetUserIdFromContext(c)
	if err != nil {
		domain.Error(c, ecode.ErrAuthFailed, err.Error())
		return
	}

	successCount := 0
	FailureCount := 0
	FailedItemIDs := make([]string, 0)
	for _, itemID := range req.ItemID {
		ItemReq := item.UpdateItemRequest{
			Id:             itemID,
			OrganizationId: treemodel.DEFAULT_ID,
			UserId:         userId,
		}
		_, err = rpc.ItemClient.UpdateItem(c, &ItemReq)
		if err != nil {
			FailureCount++
			FailedItemIDs = append(FailedItemIDs, itemID)
			continue
		}
		successCount++
	}
	orgResp := domain.AddItemsToOrganizationResponse{
		SuccessCount:  int32(successCount),
		FailureCount:  int32(FailureCount),
		FailedItemIDs: FailedItemIDs,
	}
	domain.Success(c, orgResp)
}
