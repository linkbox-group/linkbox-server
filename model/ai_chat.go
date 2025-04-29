package model

import (
	"github.com/linkbox-group/linkbox-server/rpc-gen/ai"
	"time"
)

type Chat struct {
	BaseModel
	UserID     string        `gorm:"column:user_id;type:varchar(36);not null" json:"user_id"`
	Content    string        `gorm:"type:text;not null;column:content" json:"content"`
	SendTime   time.Time     `gorm:"column:send_time" json:"send_time"`
	SenderType ai.SenderType `gorm:"type:varchar(20);not null;column:sender_type" json:"sender_type"`
}

func (Chat) TableName() string {
	return "chat"
}
