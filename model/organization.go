package model

import (
	"time"

	"github.com/linkbox-group/linkbox-server/model/treemodel"
	"github.com/linkbox-group/linkbox-server/rpc-gen/organization"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Organization struct {
	BaseModel
	treemodel.TreeModel
	UserID        string    `gorm:"type:varchar(36);not null;index:idx_user_code;comment:用户ID" json:"userId"`
	Description   *string   `gorm:"type:varchar(500);comment:描述" json:"description,omitempty"`
	IsDefault     bool      `gorm:"type:tinyint(1);not null;default:0;comment:是否默认组织" json:"isDefault"`
	IsShared      bool      `gorm:"type:tinyint(1);not null;default:0;comment:是否共享" json:"isShared"`
	ShareCode     *string   `gorm:"type:varchar(32);index:idx_share_code;comment:分享码" json:"shareCode,omitempty"`
	ShareExpireAt time.Time `gorm:"type:datetime;comment:分享过期时间" json:"shareExpireAt,omitempty"`
	SortOrder     int       `gorm:"type:int;not null;default:0;comment:排序顺序" json:"sortOrder"`
	ItemsCount    uint32    `gorm:"type:int unsigned;not null;default:0;comment:组织项目数" json:"itemsCount"`

	// 关联
	User     User           `gorm:"foreignKey:UserID" json:"-"`
	Children []Organization `gorm:"foreignKey:ParentCode;references:Code" json:"children,omitempty"`
	Items    []Item         `gorm:"many2many:organization_item;" json:"items,omitempty"`
}

func (Organization) TableName() string {
	return "organization"
}

func (o Organization) Convert() *organization.Organization {
	return &organization.Organization{
		Id:            o.ID,
		Code:          o.Code,
		ParentCode:    o.ParentCode,
		ParentCodes:   o.ParentCodes,
		TreeLeaf:      o.TreeLeaf,
		TreeLevel:     int32(o.TreeLevel),
		TreeNames:     o.TreeNames,
		Name:          o.Name,
		UserId:        o.UserID,
		Description:   *o.Description,
		IsDefault:     o.IsDefault,
		IsShared:      o.IsShared,
		ShareCode:     *o.ShareCode,
		ShareExpireAt: timestamppb.New(o.ShareExpireAt),
		SortOrder:     int32(o.SortOrder),
		ItemsCount:    o.ItemsCount,
		CreatedAt:     timestamppb.New(o.CreatedAt),
		UpdatedAt:     timestamppb.New(o.UpdatedAt),
	}

}
