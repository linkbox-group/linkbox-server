package model

import "time"

// 功能使用统计模型
type FeatureUsage struct {
	BaseModel
	UserID      string     `gorm:"type:varchar(36);not null;uniqueIndex:idx_user_feature;comment:用户ID" json:"userId"`
	FeatureName string     `gorm:"type:varchar(100);not null;uniqueIndex:idx_user_feature;comment:功能名称" json:"featureName"`
	UsageCount  uint32     `gorm:"type:int unsigned;not null;default:0;comment:使用次数" json:"usageCount"`
	LastUsedAt  *time.Time `gorm:"type:datetime;comment:最后使用时间" json:"lastUsedAt,omitempty"`

	// 关联
	User User `gorm:"foreignKey:UserID" json:"-"`
}

func (FeatureUsage) TableName() string {
	return "feature_usage"
}
