package router

import (
	"github.com/gin-gonic/gin"
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

	return router
}
