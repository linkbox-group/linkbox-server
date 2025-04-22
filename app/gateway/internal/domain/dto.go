package domain

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 定义通用的分页请求结构
type PageRequest struct {
	Page     int `form:"page" json:"page" binding:"required,min=1"`
	PageSize int `form:"page_size" json:"page_size" binding:"required,min=1,max=100"`
}

// 定义通用的分页响应结构
type PageResponse struct {
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Items    any   `json:"items"`
}

// 创建分页响应
func NewPageResponse(total int64, page, pageSize int, items any) *PageResponse {
	return &PageResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Items:    items,
	}
}

var EmptyData = struct{}{}

type Resp struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
	Data any    `json:"data"`
}

// Success 返回成功的响应
func Success(c *gin.Context, data any) {
	resp := Resp{
		Msg:  "success",
		Code: 20000,
		Data: data,
	}
	c.JSON(http.StatusOK, resp)
	c.Abort()
}

// Error 返回错误的响应
func Error(c *gin.Context, code int, msg string) {
	resp := Resp{
		Msg:  msg,
		Code: code,
		Data: EmptyData,
	}
	c.JSON(http.StatusInternalServerError, resp)
	c.Abort()
}

// ErrorMsg 返回错误的响应
func ErrorMsg(c *gin.Context, code int, msg string) {
	resp := Resp{
		Msg:  msg,
		Code: code,
		Data: EmptyData,
	}
	c.JSON(http.StatusInternalServerError, resp)
	c.Abort()
}

// GetUserIdFromContext 从上下文中获取 userId
func GetUserIdFromContext(ctx *gin.Context) (string, error) {
	userIdAny, exists := ctx.Get("userId")
	if !exists {
		return "", fmt.Errorf("用户未认证")
	}

	// 类型断言，确保 userId 是 int32
	userId, ok := userIdAny.(string)
	if !ok {
		return "", fmt.Errorf("用户ID类型错误")
	}

	return userId, nil
}
