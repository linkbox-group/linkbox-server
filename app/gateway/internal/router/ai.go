package router

import (
	"github.com/linkbox-group/linkbox-server/gateway/internal/api"
	"github.com/linkbox-group/linkbox-server/gateway/internal/middleware"
)

func (r *Group) RegisterAiRoutes() {
	var chatApi api.ChatAPI
	chatGroup := r.Group("/ai")
	r.Use(middleware.JWT())
	{
		// 发送消息
		chatGroup.GET("/chat", chatApi.SendMessage)
		// 获取消息列表
		chatGroup.GET("/chat/list", chatApi.ListMessages)
		//删除消息
		chatGroup.DELETE("/chat", chatApi.DeleteMessage)
		chatGroup.POST("/tags", chatApi.SuggestTags)
	}
}
