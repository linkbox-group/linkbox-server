package acl

import (
	"context"
	"github.com/linkbox-group/linkbox-server/model"
	"github.com/linkbox-group/linkbox-server/rpc-gen/ai"
)

type AiServiceItf interface {
	SendMessage(ctx context.Context, req *ai.SendMessageRequest, stream ai.AIService_SendMessageServer) (err error)
	ListMessages(ctx context.Context, req *ai.ListMessagesRequest) ([]*model.Chat, int, error)
	DeleteMessage(ctx context.Context, userID string, IDs []string) (success bool, err error)
}
