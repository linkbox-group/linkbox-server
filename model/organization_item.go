package model // 组织项目关联模型
type OrganizationItem struct {
	BaseModel
	OrganizationID string `gorm:"type:varchar(36);not null;uniqueIndex:idx_organization_item;comment:组织ID" json:"organizationId"`
	ItemID         string `gorm:"type:varchar(36);not null;uniqueIndex:idx_organization_item;index:idx_item_id;comment:项目ID" json:"itemId"`
	SortOrder      int    `gorm:"type:int;not null;default:0;comment:排序顺序" json:"sortOrder"`

	// 关联
	Organization Organization `gorm:"foreignKey:OrganizationID" json:"-"`
	Item         Item         `gorm:"foreignKey:ItemID" json:"-"`
}

func (OrganizationItem) TableName() string {
	return "organization_item"
}
