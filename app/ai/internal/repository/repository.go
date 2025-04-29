package repository

import (
	"context"
	"fmt"
	"github.com/google/wire"
	"github.com/linkbox-group/linkbox-server/ai/internal/acl"
	"github.com/linkbox-group/linkbox-server/ai/pkg/log"
	"github.com/linkbox-group/linkbox-server/model"
	"gorm.io/gorm"
)

var ProviderSet = wire.NewSet(wire.Bind(new(acl.AiRepositoryItf), new(*AiRepository)), NewAiRepository)

var _ acl.AiRepositoryItf = &AiRepository{}

type AiRepository struct {
	db *gorm.DB
}

func NewAiRepository(db *gorm.DB) *AiRepository {
	return &AiRepository{
		db: db,
	}
}

func (r *AiRepository) StoreMessageRecord(ctx context.Context, m *model.Chat) (err error) {
	err = r.db.Transaction(func(tx *gorm.DB) error {
		if err := r.db.Create(m).Error; err != nil {
			log.Log().Error(err.Error())
			return err
		}
		if err != nil {
			log.Log().Error(err.Error())
			return err
		}
		return nil
	})
	if err != nil {
		log.Log().Error(err.Error())
		return err
	}
	return nil
}
func (r *AiRepository) ListMessages(ctx context.Context, userID string, pageNum, pageSize int) ([]*model.Chat, int, error) {
	var messages []*model.Chat
	// 参数验证
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize < 1 {
		pageSize = 10 // 默认页大小
	}
	// 计算分页参数
	limit, offset := pageSize, (pageNum-1)*pageSize
	query := r.db.WithContext(ctx).Model(&model.Chat{}).Where("user_id = ?", userID)

	var totalCount int64
	countQuery := query.Session(&gorm.Session{}) // 创建查询的副本用于计数
	if err := countQuery.Count(&totalCount).Error; err != nil {
		return nil, 0, fmt.Errorf("counting organization items failed: %w", err)
	}
	total := int(totalCount)

	// 如果没有结果，提前返回
	if total == 0 {
		return []*model.Chat{}, 0, nil
	}

	// 获取分页数据
	if err := query.
		Order("send_time DESC").
		Offset(offset).
		Limit(limit).
		Find(&messages).Error; err != nil {
		log.Log().Error(err.Error())
		return nil, 0, fmt.Errorf("fetching chats failed: %w", err)
	}

	return messages, total, nil
}
func (r *AiRepository) DeleteMessageRecord(ctx context.Context, userID string, IDs []string) error {
	return r.db.Model(&model.Chat{}).Where("user_id = ? AND id in ?", userID, IDs).Delete(&model.Chat{}).Error

}
