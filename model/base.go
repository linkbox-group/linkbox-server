package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// 基础模型，包含通用字段
type BaseModel struct {
	ID        string    `gorm:"type:varchar(36);primaryKey;comment:UUID" json:"id"`
	CreatedAt time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;onUpdate:CURRENT_TIMESTAMP" json:"updatedAt"`
}

// 为模型生成UUID
func (base *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if base.ID == "" {
		base.ID = uuid.New().String()
	}
	return nil
}
