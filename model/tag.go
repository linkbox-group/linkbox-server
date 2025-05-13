package model

import (
	"github.com/linkbox-group/linkbox-server/rpc-gen/tag"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// 标签模型
type Tag struct {
	BaseModel
	UserID   string `gorm:"type:varchar(36);not null;uniqueIndex:idx_user_name;comment:用户ID" json:"userId"`
	Name     string `gorm:"type:varchar(100);not null;uniqueIndex:idx_user_name;comment:标签名称" json:"name"`
	Color    string `gorm:"type:varchar(20);comment:标签颜色" json:"color,omitempty"`
	Icon     string
	UseCount uint32 `gorm:"type:int unsigned;not null;default:0;comment:使用次数" json:"useCount"`
	// 关联
	User  User   `gorm:"foreignKey:UserID" json:"-"`
	Items []Item `gorm:"many2many:item_tag;" json:"items,omitempty"`
}

func (Tag) TableName() string {
	return "tag"
}
func (t Tag) ConvertTo() *tag.Tag {
	return &tag.Tag{
		Id:        t.ID,
		UserId:    t.UserID,
		Name:      t.Name,
		Color:     t.Color,
		Icon:      t.Icon,
		ItemCount: int32(t.UseCount),
		CreatedAt: timestamppb.New(t.CreatedAt),
		UpdatedAt: timestamppb.New(t.UpdatedAt),
	}
}
