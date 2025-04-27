package repository

import (
	"context"
	"github.com/linkbox-group/linkbox-server/model"
)

// GetOrganizationItems 获取组织下的所有项目
func (r *OrganizationRepository) GetOrganizationItems(ctx context.Context, orgId string) ([]*model.Item, error) {
	var items []*model.Item
	err := r.db.Joins("JOIN organization_item ON organization_item.item_id = item.id").
		Where("organization_item.organization_id = ?", orgId).
		Order("organization_item.sort_order ASC").
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

// UpdateOrganizationItemSort 更新组织项目排序
func (r *OrganizationRepository) UpdateOrganizationItemSort(ctx context.Context, orgId string, itemId string, sortOrder int) error {
	return r.db.Model(&model.Item{}).
		Where("organization_id = ? AND item_id = ?", orgId, itemId).
		Update("sort_order", sortOrder).Error
}
