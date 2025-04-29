package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/linkbox-group/linkbox-server/common/ecode"
	"github.com/linkbox-group/linkbox-server/gateway/internal/domain"
	"github.com/linkbox-group/linkbox-server/gateway/internal/infra/rpc"
	"github.com/linkbox-group/linkbox-server/rpc-gen/ai"
	"github.com/linkbox-group/linkbox-server/rpc-gen/common/pagination"
	"io"
)

type ChatAPI struct{}

func NewChatAPI() *ChatAPI {
	return &ChatAPI{}
}

// SendMessage 发送消息
func (ChatAPI) SendMessage(c *gin.Context) {
	userID, err := domain.GetUserIdFromContext(c)
	if err != nil {
		domain.ErrorMsg(c, ecode.ErrAuthFailed, err.Error())
		return
	}

	//发送消息请求
	itemID := c.Query("item_id")
	if itemID == "" {
		domain.Error(c, ecode.ErrInvalidParam, "item_id is required")
		return
	}
	content := c.Query("content")
	if content == "" {
		domain.Error(c, ecode.ErrInvalidParam, "content is required")
		return
	}
	req := &ai.SendMessageRequest{}

	req.Content = content
	req.ItemId = itemID
	req.UserId = userID
	stream, err := rpc.AiClient.SendMessage(context.Background(), req)
	if err != nil {
		domain.ErrorMsg(c, ecode.ErrRpcServiceError, err.Error())
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")

	c.Stream(func(w io.Writer) bool {
		resp, err := stream.Recv(c)
		if err == io.EOF {
			return true
		}
		if err != nil {
			return false
		}
		c.SSEvent("message", resp)
		return true
	})
}

// ListMessages 获取消息列表
func (ChatAPI) ListMessages(c *gin.Context) {
	userID, err := domain.GetUserIdFromContext(c)
	if err != nil {
		fmt.Println(err)
		domain.ErrorMsg(c, ecode.ErrAuthFailed, err.Error())
		return
	}
	req := &ai.ListMessagesRequest{}
	err = c.ShouldBind(req)
	if err != nil {
		domain.ErrorMsg(c, ecode.ErrInvalidParam, err.Error())
	}
	if req.Pagination == nil {
		req.Pagination = &pagination.PaginationRequest{
			PageSize: 10,
			Page:     1,
		}
	}
	if req.Pagination.Page <= 0 {
		req.Pagination.Page = 1
	}
	if req.Pagination.PageSize <= 0 {
		req.Pagination.PageSize = 10
	}
	req.UserId = userID

	resp, err := rpc.AiClient.ListMessages(context.Background(), req)
	if err != nil {
		domain.ErrorMsg(c, ecode.ErrRpcServiceError, err.Error())
		return
	}
	messages := make([]*domain.MessageResp, 0)
	for _, m := range resp.Messages {
		message := &domain.MessageResp{
			ID:         m.Id,
			UserId:     m.UserId,
			Content:    m.Content,
			SenderType: m.SenderType.String(),
			SendTime:   m.SendTime.AsTime(),
		}
		messages = append(messages, message)
	}
	vo := domain.ListMessagesResp{
		Messages: messages,
		Pagination: domain.Pagination{
			Total:      resp.GetPagination().TotalItems,
			Page:       resp.GetPagination().Page,
			PageSize:   resp.GetPagination().PageSize,
			TotalPages: resp.GetPagination().TotalPages,
		},
	}
	domain.Success(c, vo)
}
func (ChatAPI) DeleteMessage(c *gin.Context) {
	userID, err := domain.GetUserIdFromContext(c)
	if err != nil {
		domain.ErrorMsg(c, ecode.ErrAuthFailed, err.Error())
		return
	}
	req := &ai.DeleteMessageRequest{}
	err = c.ShouldBind(req)
	if err != nil {
		domain.ErrorMsg(c, ecode.ErrInvalidParam, err.Error())
		return
	}
	req.UserId = userID
	res, err := rpc.AiClient.DeleteMessage(context.Background(), req)
	if err != nil {
		domain.ErrorMsg(c, ecode.ErrRpcServiceError, err.Error())
		return

	}
	domain.Success(c, res.Success)

}
func (ChatAPI) SuggestTags(c *gin.Context) {
	userID, err := domain.GetUserIdFromContext(c)
	if err != nil {
		domain.ErrorMsg(c, ecode.ErrAuthFailed, err.Error())
		return
	}
	req := &ai.SuggestTagsRequest{}
	err = c.ShouldBind(req)
	if err != nil {
		domain.ErrorMsg(c, ecode.ErrInvalidParam, err.Error())
		return
	}
	req.UserId = userID
	tags, err := rpc.AiClient.SuggestTags(context.Background(), req)
	if err != nil {
		domain.ErrorMsg(c, ecode.ErrRpcServiceError, err.Error())
		return
	}
	domain.Success(c, tags.Tags)
}
