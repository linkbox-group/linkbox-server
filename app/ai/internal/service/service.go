package service

import (
	"context"
	"fmt"
	"github.com/google/wire"
	"github.com/linkbox-group/linkbox-server/ai/internal/acl"
	"github.com/linkbox-group/linkbox-server/ai/pkg/llm"
	"github.com/linkbox-group/linkbox-server/ai/pkg/log"
	"github.com/linkbox-group/linkbox-server/model"
	"github.com/linkbox-group/linkbox-server/rpc-gen/ai"
	"io"
	"time"
)

var ProviderSet = wire.NewSet(wire.Bind(new(acl.AiServiceItf), new(*AiService)), NewAiService)
var _ acl.AiServiceItf = &AiService{}

type AiService struct {
	repo acl.AiRepositoryItf
}

func NewAiService(repo acl.AiRepositoryItf) *AiService {
	return &AiService{repo: repo}
}
func (s *AiService) DeleteMessage(ctx context.Context, userID string, IDs []string) (success bool, err error) {
	err = s.repo.DeleteMessageRecord(ctx, userID, IDs)
	if err != nil {
		log.Log().Error(err.Error())
		return false, err
	}
	return true, nil

}
func (s *AiService) ListMessages(ctx context.Context, req *ai.ListMessagesRequest) ([]*model.Chat, int, error) {

	chats, total, err := s.repo.ListMessages(ctx, req.UserId, int(req.Pagination.Page), int(req.Pagination.PageSize))
	if err != nil {
		log.Log().Error(err.Error())
		return nil, 0, fmt.Errorf("failed to get messages from cache: %v", err)
	}

	return chats, total, nil
}
func (s *AiService) SendMessage(ctx context.Context, req *ai.SendMessageRequest, stream ai.AIService_SendMessageServer) (err error) {
	userMessage := &model.Chat{
		UserID:     req.UserId,
		Content:    req.GetContent(),
		SenderType: ai.SenderType_SENDER_USER,
		SendTime:   time.Now(),
	}
	err = s.repo.StoreMessageRecord(ctx, userMessage)
	if err != nil {
		log.Log().Error(err.Error())
		return err
	}
	//itemRes, err := rpc.ItemClient.GetItem(ctx, &item.GetItemRequest{
	//	Id:     req.ItemId,
	//	UserId: req.UserId,
	//})
	//if err != nil {
	//	log.Log().Error(err.Error())
	//	return err
	//}
	//itemData := itemRes.GetItem()
	//if itemData.Type == itemmodel.ItemType_LINK {
	//
	//}
	//knowledge := itemData.Note
	sr, err := llm.Chat(ctx, req.GetContent(), "")
	if err != nil {
		log.Log().Error(err.Error())
		return
	}
	defer sr.Close()
	i := 0
	resp := &ai.SendMessageResponse{
		Message: &ai.Message{},
	}
	wholeContent := ""
	for {
		msg, err := sr.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Log().Error(err.Error())
			return err
		}
		resp.Message.Content = msg.Content
		wholeContent += msg.Content
		err = stream.Send(ctx, resp)
		if err != nil {
			log.Log().Error(err.Error())
			return err
		}
		i++
	}

	messageModel := model.Chat{
		UserID:     req.UserId,
		Content:    wholeContent,
		SenderType: ai.SenderType_SENDER_AI,
		SendTime:   time.Now(),
	}
	err = s.repo.StoreMessageRecord(ctx, &messageModel)
	if err != nil {
		log.Log().Error(err.Error())
		return err

	}
	return nil
}
