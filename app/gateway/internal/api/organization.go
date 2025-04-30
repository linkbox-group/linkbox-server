package api

import (
	"context"
	"github.com/linkbox-group/linkbox-server/common/ecode"
	"github.com/linkbox-group/linkbox-server/gateway/pkg/log"
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
	UserID, err := domain.GetUserIDFromContext(c)
	if err != nil {
		domain.Error(c, ecode.ErrAuthFailed, "用户认证失败")
		return
	}
	orgReq := organization.CreateOrganizationRequest{
		UserId:      UserID,
		Name:        req.Name,
		ParentCode:  &req.ParentCode,
		Description: &req.Description,
		SortOrder:   &req.SortOrder,
	}
	resp, err := rpc.OrganizationClient.CreateOrganization(context.Background(), &orgReq)
	if err != nil {
		domain.Error(c, ErrOrganizationNameExists, "创建组织失败")
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
	log.Log().Info("[a *OrganizationAPI] CreateOrganization ", "user_id", UserID)
}

// GetOrganization 获取组织详情
func (a *OrganizationAPI) GetOrganization(c *gin.Context) {
	orgID := c.Param("id")
	UserID, err := domain.GetUserIDFromContext(c)
	if err != nil {
		domain.Error(c, ecode.ErrAuthFailed, "用户认证失败")
		return
	}

	resp, err := rpc.OrganizationClient.GetOrganization(context.Background(), &organization.GetOrganizationRequest{
		Id:     orgID,
		UserId: UserID,
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
	log.Log().Info("[a *OrganizationAPI] GetOrganization ", "user_id", UserID)

	domain.Success(c, orgResp)
}

// UpdateOrganization 更新组织
func (a *OrganizationAPI) UpdateOrganization(c *gin.Context) {
	var req organization.UpdateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		domain.Error(c, ErrInvalidReq, "请求参数错误")
		return
	}
	UserID, err := domain.GetUserIDFromContext(c)
	if err != nil {
		domain.Error(c, ecode.ErrAuthFailed, "用户认证失败")
		return
	}
	req.UserId = UserID

	resp, err := rpc.OrganizationClient.UpdateOrganization(context.Background(), &req)
	if err != nil {
		domain.Error(c, ErrOrganizationNameExists, "更新组织失败")
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
	log.Log().Info("[a *OrganizationAPI] UpdateOrganization ", "user_id", UserID)
	domain.Success(c, orgResp)
}

// DeleteOrganization 删除组织
func (a *OrganizationAPI) DeleteOrganization(c *gin.Context) {
	orgID := c.Param("id")
	cascade := c.DefaultQuery("cascade", "false")
	userId, err := domain.GetUserIDFromContext(c)
	if err != nil {
		domain.Error(c, ecode.ErrAuthFailed, err.Error())
		return
	}
	req := &organization.DeleteOrganizationRequest{
		UserId:  userId,
		Id:      orgID,
		Cascade: cascade == "true",
	}
	resp, err := rpc.OrganizationClient.DeleteOrganization(context.Background(), req)
	if err != nil {
		domain.Error(c, ErrOrganizationNotFound, "删除组织失败")
		return
	}
	orgResp := domain.OrganizationSuccessResponse{
		Success: resp.GetSuccess(),
	}
	log.Log().Info("[a *OrganizationAPI] DeleteOrganization ", "user_id", userId)
	domain.Success(c, orgResp)
}

// GetUserOrganizations 获取用户组织列表
func (a *OrganizationAPI) GetUserOrganizations(c *gin.Context) {
	UserID, err := domain.GetUserIDFromContext(c)
	if err != nil {
		domain.Error(c, ecode.ErrAuthFailed, "用户认证失败")
		return
	}

	resp, err := rpc.OrganizationClient.GetUserOrganizations(context.Background(), &organization.GetUserOrganizationsRequest{
		UserId: UserID,
	})
	if err != nil {
		domain.Error(c, ecode.ErrRpcServiceError, "获取组织列表失败")
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
	log.Log().Info("[a *OrganizationAPI] GetUserOrganizations ", "user_id", UserID)
	domain.Success(c, orgsResp)

}

// GetOrganizationTree 获取组织树
func (a *OrganizationAPI) GetOrganizationTree(c *gin.Context) {
	UserID, err := domain.GetUserIDFromContext(c)
	if err != nil {
		domain.Error(c, ecode.ErrAuthFailed, "用户认证失败")
		return
	}
	rootCode := c.Query("root_code")

	resp, err := rpc.OrganizationClient.GetOrganizationTree(context.Background(), &organization.GetOrganizationTreeRequest{
		UserId:   UserID,
		RootCode: &rootCode,
	})
	if err != nil {
		domain.Error(c, ErrOrganizationNotFound, "获取组织树失败")
		return
	}
	log.Log().Info("[a *OrganizationAPI] GetOrganizationTree ", "user_id", UserID)
	domain.Success(c, resp.GetRoot().GetData())
}

// GetOrganizationChildren 获取组织子节点
func (a *OrganizationAPI) GetOrganizationChildren(c *gin.Context) {
	UserID, err := domain.GetUserIDFromContext(c)
	if err != nil {
		domain.Error(c, ecode.ErrAuthFailed, "用户认证失败")
		return
	}
	parentCode := c.Query("parent_code")
	recursive := c.DefaultQuery("recursive", "false")

	if parentCode == "" {
		domain.Error(c, ErrInvalidReq, "请求参数错误")
		return
	}

	resp, err := rpc.OrganizationClient.GetOrganizationChildren(context.Background(), &organization.GetOrganizationChildrenRequest{
		UserId:     UserID,
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
	UserID, err := domain.GetUserIDFromContext(c)
	if err != nil {
		domain.Error(c, ecode.ErrAuthFailed, "用户认证失败")
		return
	}
	req.UserId = UserID
	resp, err := rpc.OrganizationClient.MoveOrganization(context.Background(), &req)
	if err != nil {
		domain.Error(c, ErrOrganizationNotFound, "移动组织失败")
		return
	}

	org := resp.GetSuccess()
	orgResp := domain.OrganizationSuccessResponse{
		Success: org,
	}
	log.Log().Info("[a *OrganizationAPI] MoveOrganization ", "user_id", UserID)
	domain.Success(c, orgResp)
}

// AddItemsToOrganization 添加内容项到组织
func (a *OrganizationAPI) AddItemsToOrganization(c *gin.Context) {
	var req domain.AddItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		domain.Error(c, ErrInvalidReq, "请求参数错误")
		return
	}
	UserID, err := domain.GetUserIDFromContext(c)
	if err != nil {
		domain.Error(c, ecode.ErrAuthFailed, "用户认证失败")
		return
	}
	successCount := 0
	FailureCount := 0
	FailedItemIDs := make([]string, 0)
	for _, itemID := range req.ItemID {
		ItemReq := item.UpdateItemRequest{
			Id:             itemID,
			OrganizationId: req.OrganizationID,
			UserId:         UserID,
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
	log.Log().Info("[a *OrganizationAPI] AddItemsToOrganization ", "user_id", UserID)
	domain.Success(c, orgResp)
}

// RemoveItemsFromOrganization 从组织移除内容项
func (a *OrganizationAPI) RemoveItemsFromOrganization(c *gin.Context) {
	var req domain.RemoveItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		domain.Error(c, ErrInvalidReq, "请求参数错误")
		return
	}
	userId, err := domain.GetUserIDFromContext(c)
	if err != nil {
		domain.Error(c, ecode.ErrAuthFailed, "添加内容失败")
		return
	}

	successCount := 0
	FailureCount := 0
	FailedItemIDs := make([]string, 0)
	orgID, err := rpc.OrganizationClient.GetDefaultOrgID(c, &organization.GetDefaultOrgIDReq{
		UserId: userId,
		Code:   treemodel.ROOT_ID,
	})
	if err != nil {
		domain.Error(c, ecode.ErrRpcServiceError, "移除内容项失败")
		return
	}
	for _, itemID := range req.ItemID {
		ItemReq := item.UpdateItemRequest{
			Id:             itemID,
			OrganizationId: orgID.GetId(),
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
