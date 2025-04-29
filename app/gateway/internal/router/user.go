package router

import (
	"github.com/linkbox-group/linkbox-server/gateway/internal/api"
	"github.com/linkbox-group/linkbox-server/gateway/internal/middleware"
)

func (r *Group) SetUserRouter() {
	group := r.Group("/user")
	var userApi api.UserApi
	group.POST("/send_code", userApi.SendCode)
	group.POST("/login", userApi.Login)
	group.POST("/register", userApi.Register)
	group.POST("/refresh-token", userApi.RefreshToken)

	groupAuthed := group.Use(middleware.JWT())
	groupAuthed.PUT("/info", userApi.UpdateUserInfo)
	groupAuthed.PUT("/password", userApi.ChangePassword)
	groupAuthed.GET("/info", userApi.GetUserInfo)
	group.DELETE("", userApi.DeleteUser)
}
