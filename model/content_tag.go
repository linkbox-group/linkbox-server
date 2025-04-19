package model

// 项目标签关联模型
type ContentTag struct {
	BaseModel
	ItemID string `gorm:"type:varchar(36);not null;uniqueIndex:idx_item_tag;comment:项目ID" json:"itemId"`
	TagID  string `gorm:"type:varchar(36);not null;uniqueIndex:idx_item_tag;index:idx_tag_id;comment:标签ID" json:"tagId"`

	// 关联
	Item Item `gorm:"foreignKey:ItemID" json:"-"`
	Tag  Tag  `gorm:"foreignKey:TagID" json:"-"`
}

func (ContentTag) TableName() string {
	return "content_tag"
}
