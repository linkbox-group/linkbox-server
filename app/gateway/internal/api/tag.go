package api

import (
	"context"
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
	var req tag.CreateTagRequest
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

	resp, err := rpc.TagClient.CreateTag(context.Background(), &req)
	if err != nil {
		domain.Error(c, ErrTagNameExists, err.Error())
		return
	}

	tag := resp.GetTag()
	tagResp := domain.TagResponse{
		ID:          tag.Id,
		UserID:      tag.UserId,
		Name:        tag.Name,
		Description: tag.Description,
		Color:       tag.Color,
		ItemCount:   tag.ItemCount,
		CreatedAt:   tag.CreatedAt.AsTime(),
		UpdatedAt:   tag.UpdatedAt.AsTime(),
	}
	domain.Success(c, tagResp)
}

// GetTag 获取标签
func (a *TagAPI) GetTag(c *gin.Context) {
	tagID := c.Param("id")
	userId, err := domain.GetUserIdFromContext(c)
	if err != nil {
		domain.Error(c, ErrAuthFailedCode, err.Error())
		return
	}
	resp, err := rpc.TagClient.GetTag(context.Background(), &tag.GetTagRequest{
		Id:     tagID,
		UserId: userId,
	})
	if err != nil {
		domain.Error(c, ErrTagNotFound, "标签不存在")
		return
	}

	tag := resp.GetTag()
	tagResp := domain.TagResponse{
		ID:          tag.Id,
		UserID:      tag.UserId,
		Name:        tag.Name,
		Description: tag.Description,
		Color:       tag.Color,
		ItemCount:   tag.ItemCount,
		CreatedAt:   tag.CreatedAt.AsTime(),
		UpdatedAt:   tag.UpdatedAt.AsTime(),
	}
	domain.Success(c, tagResp)
}

// UpdateTag 更新标签
func (a *TagAPI) UpdateTag(c *gin.Context) {
	tagID := c.Param("id")
	var req tag.UpdateTagRequest

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

	req.Id = tagID
	resp, err := rpc.TagClient.UpdateTag(context.Background(), &req)
	if err != nil {
		domain.Error(c, ErrTagNameExists, err.Error())
		return
	}

	tag := resp.GetTag()
	tagResp := domain.TagResponse{
		ID:          tag.Id,
		UserID:      tag.UserId,
		Name:        tag.Name,
		Description: tag.Description,
		Color:       tag.Color,
		ItemCount:   tag.ItemCount,
		CreatedAt:   tag.CreatedAt.AsTime(),
		UpdatedAt:   tag.UpdatedAt.AsTime(),
	}
	domain.Success(c, tagResp)
}

// DeleteTag 删除标签
func (a *TagAPI) DeleteTag(c *gin.Context) {
	tagID := c.Param("id")
	userID, err := domain.GetUserIdFromContext(c)
	if err != nil {
		domain.Error(c, ErrAuthFailedCode, err.Error())
		return
	}

	resp, err := rpc.TagClient.DeleteTag(context.Background(), &tag.DeleteTagRequest{
		Id:     tagID,
		UserId: userID,
	})
	if err != nil {
		domain.Error(c, ErrTagNotFound, "标签不存在")
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
	userID, err := domain.GetUserIdFromContext(c)
	if err != nil {
		log.Log().Error(err.Error())
		domain.Error(c, ErrAuthFailedCode, err.Error())
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	searchQuery := c.Query("search_query")

	resp, err := rpc.TagClient.GetUserTags(context.Background(), &tag.GetUserTagsRequest{
		UserId: userID,
		Pagination: &pagination.PaginationRequest{
			Page:     int32(page),
			PageSize: int32(pageSize),
		},
		SearchQuery: &searchQuery,
	})
	if err != nil {
		domain.Error(c, ErrRpcFailedCode, "rpc调用失败")
		return
	}

	var tags []domain.TagResponse
	for _, tag := range resp.GetTags().Tags {
		tags = append(tags, domain.TagResponse{
			ID:          tag.Id,
			UserID:      tag.UserId,
			Name:        tag.Name,
			Description: tag.Description,
			Color:       tag.Color,
			ItemCount:   tag.ItemCount,
			CreatedAt:   tag.CreatedAt.AsTime(),
			UpdatedAt:   tag.UpdatedAt.AsTime(),
		})
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
	userID, err := domain.GetUserIdFromContext(c)
	if err != nil {
		domain.Error(c, ErrAuthFailedCode, err.Error())
		return
	}

	var req tag.AddTagsToItemsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		domain.Error(c, ErrInvalidReq, "请求参数错误")
		return
	}
	req.UserId = userID

	resp, err := rpc.TagClient.AddTagsToItems(context.Background(), &req)
	if err != nil {
		log.Log().Error(err.Error())
		domain.Error(c, ErrRpcFailedCode, "rpc调用失败"+err.Error())
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
	userID, err := domain.GetUserIdFromContext(c)
	if err != nil {
		domain.Error(c, ErrAuthFailedCode, err.Error())
		return
	}

	var req tag.RemoveTagsFromItemsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		domain.Error(c, ErrInvalidReq, "请求参数错误")
		return
	}
	req.UserId = userID

	resp, err := rpc.TagClient.RemoveTagsFromItems(context.Background(), &req)
	if err != nil {
		domain.Error(c, ErrNoPermission, "没有操作权限")
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
	userID, err := domain.GetUserIdFromContext(c)
	if err != nil {
		domain.Error(c, ErrAuthFailedCode, err.Error())
		return
	}

	resp, err := rpc.TagClient.GetItemTags(context.Background(), &tag.GetItemTagsRequest{
		ItemId: itemID,
		UserId: userID,
	})
	if err != nil {
		domain.Error(c, ErrItemNotFound, "内容项不存在")
		return
	}

	var tags []domain.TagResponse
	for _, tag := range resp.GetTags().Tags {
		tags = append(tags, domain.TagResponse{
			ID:          tag.Id,
			UserID:      tag.UserId,
			Name:        tag.Name,
			Description: tag.Description,
			Color:       tag.Color,
			ItemCount:   tag.ItemCount,
			CreatedAt:   tag.CreatedAt.AsTime(),
			UpdatedAt:   tag.UpdatedAt.AsTime(),
		})
	}

	// 使用domain中定义的ItemTagsResponse结构体
	responseData := domain.ItemTagsResponse{
		ItemID: itemID,
		Tags:   tags,
	}
	domain.Success(c, responseData)
}

// GetRelatedTags 获取相关标签
func (a *TagAPI) GetRelatedTags(c *gin.Context) {
	tagID := c.Param("tag_id")
	userID, err := domain.GetUserIdFromContext(c)
	if err != nil {
		domain.Error(c, ErrAuthFailedCode, err.Error())
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
	limit32 := int32(limit)

	resp, err := rpc.TagClient.GetRelatedTags(context.Background(), &tag.GetRelatedTagsRequest{
		TagId:  tagID,
		UserId: userID,
		Limit:  &limit32,
	})
	if err != nil {
		domain.Error(c, ErrTagNotFound, "标签不存在")
		return
	}

	var relatedTags []domain.TagRelatedResponse

	for _, relatedTag := range resp.GetRelatedTags().RelatedTags {
		tag := relatedTag.GetTag()
		relatedTags = append(relatedTags, domain.TagRelatedResponse{
			Tag: domain.TagResponse{
				ID:          tag.Id,
				UserID:      tag.UserId,
				Name:        tag.Name,
				Description: tag.Description,
				Color:       tag.Color,
				ItemCount:   tag.ItemCount,
				CreatedAt:   tag.CreatedAt.AsTime(),
				UpdatedAt:   tag.UpdatedAt.AsTime(),
			},
			Correlation:  relatedTag.GetCorrelation(),
			CoOccurrence: relatedTag.GetCoOccurrence(),
		})
	}

	// 使用domain中定义的GetRelatedTagsFullResponse结构体
	responseData := domain.GetRelatedTagsFullResponse{
		TagID:       tagID,
		RelatedTags: relatedTags,
	}
	domain.Success(c, responseData)
}
