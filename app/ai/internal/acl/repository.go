package acl

import (
	"context"
	"github.com/linkbox-group/linkbox-server/model"
)

type AiRepositoryItf interface {
	StoreMessageRecord(ctx context.Context, m *model.Chat) (err error)
	ListMessages(ctx context.Context, userID string, pageNum, pageSize int) ([]*model.Chat, int, error)
	DeleteMessageRecord(ctx context.Context, userID string, IDs []string) error
}
