package model

// 项目模型
type Item struct {
	BaseModel
	UserID       string `gorm:"type:varchar(36);not null;index:idx_user_id;comment:用户ID" json:"userId"`
	ItemType     string `gorm:"type:varchar(20);not null;index:idx_type;comment:类型:text,image,link,bookmark" json:"itemType"`
	Title        string `gorm:"type:varchar(500);comment:标题" json:"title,omitempty"`
	Note         string `gorm:"type:text;comment:内容/文本" json:"note,omitempty"`
	URL          string `gorm:"type:varchar(2000);comment:链接地址" json:"url,omitempty"`
	ThumbnailURL string `gorm:"type:varchar(2000);comment:缩略图地址" json:"thumbnailUrl,omitempty"`
	// 关联
	User          User           `gorm:"foreignKey:UserID" json:"-"`
	Organizations []Organization `gorm:"many2many:organization_item;" json:"organizations,omitempty"`
	Tags          []Tag          `gorm:"many2many:item_tag;" json:"tags,omitempty"`
}

func (Item) TableName() string {
	return "item"
}
