package router

import (
	"github.com/linkbox-group/linkbox-server/gateway/internal/api"
	"github.com/linkbox-group/linkbox-server/gateway/internal/middleware"
)

func (r *Group) RegisterTrashRoutes() {
	var TrashApi api.TrashApi
	chatGroup := r.Group("/trash")
	r.Use(middleware.JWT())
	{
		//获取回收站数据
		chatGroup.GET("/list", TrashApi.ListTrashItem)
		// 恢复回收站内容
		chatGroup.POST("/recovery", TrashApi.RecoveryTrashItem)
		//删除回收站内容
		chatGroup.DELETE("/:item_id", TrashApi.DeleteTrashItem)
	}
}
