package router

import (
	"github.com/gin-gonic/gin"
	"github.com/linkbox-group/linkbox-server/gateway/internal/api"
	"github.com/linkbox-group/linkbox-server/gateway/internal/middleware"
)

type Group struct {
	*gin.RouterGroup
}

func InitRouter() *gin.Engine {
	router := gin.Default()
	router.Use(middleware.CORS())
	routerGroup := Group{router.Group("/api")}
	// 用户服务
	routerGroup.SetUserRouter()
	// 标签服务
	routerGroup.RegisterTagRoutes()
	// 内容服务
	routerGroup.RegisterItemRoutes()
	// 组织服务
	routerGroup.RegisterOrganizationRoutes()
	routerGroup.RegisterAiRoutes()

	return router
}

// RegisterTagRoutes registers tag-related routes
func (r *Group) RegisterTagRoutes() {
	r.Use(middleware.JWT())
	var tagAPI api.TagAPI
	tagGroup := r.Group("/tags")

	{
		tagGroup.POST("", tagAPI.CreateTag)
		tagGroup.GET("/:id", tagAPI.GetTag)
		tagGroup.PUT("/:id", tagAPI.UpdateTag)
		tagGroup.DELETE("/:id", tagAPI.DeleteTag)
		tagGroup.POST("/items", tagAPI.AddTagsToItems)
		tagGroup.DELETE("/items", tagAPI.RemoveTagsFromItems)
	}

	userTagGroup := r.Group("/users/tags")
	{
		userTagGroup.GET("", tagAPI.GetUserTags)
	}

	itemTagGroup := r.Group("/items/tags/:items_id")
	{
		itemTagGroup.GET("", tagAPI.GetItemTags)
	}
}

// RegisterItemRoutes registers content-related routes
func (r *Group) RegisterItemRoutes() {
	r.Use(middleware.JWT())
	var contentAPI api.ItemAPI
	contentGroup := r.Group("/items")
	{
		contentGroup.POST("", contentAPI.CreateItem)
		contentGroup.GET("/:id", contentAPI.GetItem)
		contentGroup.PUT("/:id", contentAPI.UpdateItem)
		contentGroup.DELETE("/:id", contentAPI.DeleteItem)
		contentGroup.POST("/search", contentAPI.SearchItems)
	}

	contentTagGroup := r.Group("/items/tags")
	{
		contentTagGroup.POST("", contentAPI.GetItemsByTags)
	}
	organizationGroup := r.Group("/items/organization")
	{
		organizationGroup.POST("", contentAPI.GetItemsByOrganization)
	}
}

// RegisterOrganizationRoutes registers organization-related routes
func (r *Group) RegisterOrganizationRoutes() {
	r.Use(middleware.JWT())
	var orgAPI api.OrganizationAPI
	orgGroup := r.Group("/organization")
	{
		// 基础组织操作
		orgGroup.POST("", orgAPI.CreateOrganization)
		orgGroup.GET("/:id", orgAPI.GetOrganization)
		orgGroup.PUT("", orgAPI.UpdateOrganization)
		orgGroup.DELETE("/:id", orgAPI.DeleteOrganization)

		// 组织树相关
		orgGroup.GET("/tree", orgAPI.GetOrganizationTree)
		orgGroup.GET("/children", orgAPI.GetOrganizationChildren)
		orgGroup.PATCH("/move", orgAPI.MoveOrganization)

		// 组织内容项操作
		orgGroup.POST("/items", orgAPI.AddItemsToOrganization)
		orgGroup.DELETE("/items", orgAPI.RemoveItemsFromOrganization)

	}

	// 获取用户组织列表
	orgGroup.GET("", orgAPI.GetUserOrganizations)
}
