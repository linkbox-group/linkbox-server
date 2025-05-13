package api

import (
	"context"
	"github.com/linkbox-group/linkbox-server/common/ecode"
	"github.com/linkbox-group/linkbox-server/gateway/pkg/log"
	"strconv"

	"github.com/linkbox-group/linkbox-server/rpc-gen/common/pagination"
	"github.com/linkbox-group/linkbox-server/rpc-gen/tag"

	"github.com/gin-gonic/gin"
	"github.com/linkbox-group/linkbox-server/gateway/internal/domain"
	"github.com/linkbox-group/linkbox-server/gateway/internal/infra/rpc"
)

type TagAPI struct{}

func NewTagAPI() *TagAPI {
	return &TagAPI{}
}

// 错误码定义
const (
	ErrInvalidReq    = 10002 // 请求参数错误
	ErrTagNotFound   = 40000 // 标签不存在
	ErrTagNameExists = 40001 // 标签名已存在
	ErrNotLoggedIn   = 30000 // 未登录或登录已过期
	ErrNoPermission  = 30001 // 没有操作权限
	ErrItemNotFound  = 40000 // 内容项不存在
)

// 错误码定义

// CreateTag 创建标签
func (a *TagAPI) CreateTag(c *gin.Context) {
	UserID, err := domain.GetUserIDFromContext(c)
	if err != nil {
		domain.Error(c, ecode.ErrAuthFailed, err.Error())
		return
	}
	log.Log().Info("[a *TagAPI] CreateTag ", "user_id", UserID)

	var req tag.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		domain.Error(c, ErrInvalidReq, "请求参数错误")
		return
	}
	req.UserId = UserID

	resp, err := rpc.TagClient.CreateTag(context.Background(), &req)
	if err != nil {
		domain.Error(c, ErrTagNameExists, "创建标签失败")
		return
	}

	tagResp := domain.TagResponse{}

	domain.Success(c, tagResp.Convert(resp.GetTag()))
}

// GetTag 获取标签
func (a *TagAPI) GetTag(c *gin.Context) {
	tagID := c.Param("id")
	userID, err := domain.GetUserIDFromContext(c)
	if err != nil {
		domain.Error(c, ecode.ErrAuthFailed, "用户认证失败")
		return
	}
	log.Log().Info("[a *TagAPI] GetTag ", "user_id", userID)
	resp, err := rpc.TagClient.GetTag(context.Background(), &tag.GetTagRequest{
		Id:     tagID,
		UserId: userID,
	})
	if err != nil {
		domain.Error(c, ErrTagNotFound, "获取标签失败")
		return
	}

	tagResp := domain.TagResponse{}
	domain.Success(c, tagResp.Convert(resp.GetTag()))
}

// UpdateTag 更新标签
func (a *TagAPI) UpdateTag(c *gin.Context) {
	userId, err := domain.GetUserIDFromContext(c)
	if err != nil {
		domain.Error(c, ecode.ErrAuthFailed, "用户认证失败")
		return
	}
	log.Log().Info("[a *TagAPI] UpdateTag ", "user_id", userId)
	tagID := c.Param("id")
	var req tag.UpdateTagRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		domain.Error(c, ErrInvalidReq, "请求参数错误")
		return
	}

	req.UserId = userId

	req.Id = tagID
	resp, err := rpc.TagClient.UpdateTag(context.Background(), &req)
	if err != nil {
		domain.Error(c, ErrTagNameExists, "更新标签失败")
		return
	}
	tagResp := domain.TagResponse{}
	domain.Success(c, tagResp.Convert(resp.GetTag()))
}

// DeleteTag 删除标签
func (a *TagAPI) DeleteTag(c *gin.Context) {
	tagID := c.Param("id")
	userID, err := domain.GetUserIDFromContext(c)
	if err != nil {
		domain.Error(c, ecode.ErrAuthFailed, "用户认证失败")
		return
	}
	log.Log().Info("[a *TagAPI] DeleteTag ", "user_id", userID)

	resp, err := rpc.TagClient.DeleteTag(context.Background(), &tag.DeleteTagRequest{
		Id:     tagID,
		UserId: userID,
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

// GetUserTags 获取用户所有标签
func (a *TagAPI) GetUserTags(c *gin.Context) {
	userID, err := domain.GetUserIDFromContext(c)
	if err != nil {
		log.Log().Error(err.Error())
		domain.Error(c, ecode.ErrAuthFailed, "用户认证失败")
		return
	}
	log.Log().Info("[a *TagAPI] GetUserTags ", "user_id", userID)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	resp, err := rpc.TagClient.GetUserTags(context.Background(), &tag.GetUserTagsRequest{
		UserId: userID,
		Pagination: &pagination.PaginationRequest{
			Page:     int32(page),
			PageSize: int32(pageSize),
		},
	})
	if err != nil {
		domain.Error(c, ecode.ErrRpcServiceError, "获取标签失败")
		return
	}

	var tags []domain.TagResponse
	for _, t := range resp.GetTags().Tags {
		tagResp := domain.TagResponse{}
		tagResp.Convert(t)
		tags = append(tags, tagResp)
	}

	tagListResp := domain.TagListResponse{
		Tags: tags,
		Pagination: domain.Pagination{
			Total:      resp.GetTags().GetPagination().GetTotalItems(),
			Page:       resp.GetTags().GetPagination().GetPage(),
			PageSize:   resp.GetTags().GetPagination().GetPageSize(),
			TotalPages: resp.GetTags().GetPagination().GetTotalPages(),
		},
	}
	domain.Success(c, tagListResp)
}

// AddTagsToItems 添加标签到内容项
func (a *TagAPI) AddTagsToItems(c *gin.Context) {
	userID, err := domain.GetUserIDFromContext(c)
	if err != nil {
		domain.Error(c, ecode.ErrAuthFailed, "用户认证失败")
		return
	}
	log.Log().Info("[a *TagAPI] AddTagsToItems ", "user_id", userID)

	var req tag.AddTagsToItemsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		domain.Error(c, ErrInvalidReq, "请求参数错误")
		return
	}
	req.UserId = userID

	resp, err := rpc.TagClient.AddTagsToItems(context.Background(), &req)
	if err != nil {
		log.Log().Error(err.Error())
		domain.Error(c, ecode.ErrRpcServiceError, "添加标签失败")
		return
	}

	// 使用domain定义的结构体响应
	responseData := domain.AddTagsToItemsResponse{
		SuccessCount:  resp.GetData().SuccessCount,
		FailureCount:  resp.GetData().FailureCount,
		FailedItemIDs: resp.GetData().FailedItemIds,
	}
	domain.Success(c, responseData)
}

// RemoveTagsFromItems 从内容项移除标签
func (a *TagAPI) RemoveTagsFromItems(c *gin.Context) {
	userID, err := domain.GetUserIDFromContext(c)
	if err != nil {
		domain.Error(c, ecode.ErrAuthFailed, "用户认证失败")
		return
	}
	log.Log().Info("[a *TagAPI] RemoveTagsFromItems ", "user_id", userID)

	var req tag.RemoveTagsFromItemsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		domain.Error(c, ErrInvalidReq, "请求参数错误")
		return
	}
	req.UserId = userID

	resp, err := rpc.TagClient.RemoveTagsFromItems(context.Background(), &req)
	if err != nil {
		domain.Error(c, ErrNoPermission, "移除标签失败")
		return
	}

	// 使用domain中定义的结构体响应
	responseData := domain.RemoveTagsFromItemsResponse{
		SuccessCount:  resp.GetData().SuccessCount,
		FailureCount:  resp.GetData().FailureCount,
		FailedItemIDs: resp.GetData().FailedItemIds,
	}
	domain.Success(c, responseData)
}

// GetItemTags 获取内容项的标签
func (a *TagAPI) GetItemTags(c *gin.Context) {
	itemID := c.Param("item_id")
	userID, err := domain.GetUserIDFromContext(c)
	if err != nil {
		domain.Error(c, ecode.ErrRpcServiceError, "rpc调用失败"+err.Error())
		return
	}
	log.Log().Info("[a *TagAPI] GetItemTags ", "user_id", userID)

	resp, err := rpc.TagClient.GetItemTags(context.Background(), &tag.GetItemTagsRequest{
		ItemId: itemID,
		UserId: userID,
	})
	if err != nil {
		domain.Error(c, ErrItemNotFound, "获取标签失败")
		return
	}

	var tags []domain.TagResponse
	for _, t := range resp.GetTags().Tags {
		tagResp := domain.TagResponse{}
		tagResp.Convert(t)
		tags = append(tags, tagResp)
	}

	// 使用domain中定义的ItemTagsResponse结构体
	responseData := domain.ItemTagsResponse{
		ItemID: itemID,
		Tags:   tags,
	}
	domain.Success(c, responseData)
}
