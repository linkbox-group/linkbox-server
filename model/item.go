package model

import "time"

// 项目模型
type Item struct {
	BaseModel
	UserID          string     `gorm:"type:varchar(36);not null;index:idx_user_id;comment:用户ID" json:"userId"`
	ItemType        string     `gorm:"type:varchar(20);not null;index:idx_type;comment:类型:text,image,link,bookmark" json:"itemType"`
	Title           *string    `gorm:"type:varchar(500);comment:标题" json:"title,omitempty"`
	Content         *string    `gorm:"type:text;comment:内容/文本" json:"content,omitempty"`
	URL             *string    `gorm:"type:varchar(2000);comment:链接地址" json:"url,omitempty"`
	ImageURL        *string    `gorm:"type:varchar(2000);comment:图片地址" json:"imageUrl,omitempty"`
	ThumbnailURL    *string    `gorm:"type:varchar(2000);comment:缩略图地址" json:"thumbnailUrl,omitempty"`
	SourceDomain    *string    `gorm:"type:varchar(255);index:idx_domain;comment:来源网站域名" json:"sourceDomain,omitempty"`
	SourcePageTitle *string    `gorm:"type:varchar(500);comment:来源页面标题" json:"sourcePageTitle,omitempty"`
	IsFavorited     bool       `gorm:"type:tinyint(1);not null;default:0;comment:是否收藏" json:"isFavorited"`
	FavoritedAt     *time.Time `gorm:"type:datetime;comment:收藏时间" json:"favoritedAt,omitempty"`
	IsDeleted       bool       `gorm:"type:tinyint(1);not null;default:0;index:idx_user_deleted;comment:是否删除" json:"isDeleted"`
	DeletedAt       *time.Time `gorm:"type:datetime;index:idx_deleted_at;comment:删除时间" json:"deletedAt,omitempty"`

	// 关联
	User          User           `gorm:"foreignKey:UserID" json:"-"`
	Organizations []Organization `gorm:"many2many:organization_item;" json:"organizations,omitempty"`
	Tags          []Tag          `gorm:"many2many:item_tag;" json:"tags,omitempty"`
}

func (Item) TableName() string {
	return "item"
}
