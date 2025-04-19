package repository

import (
	"context"
	"errors"
	"github.com/linkbox-group/linkbox-server/model"
)

// CreateOrganizationItem 创建组织项目关联
func (r *OrganizationRepository) CreateOrganizationItem(ctx context.Context, orgItem *model.OrganizationItem) error {
	// 检查是否已存在关联
	var count int64
	err := r.db.Model(&model.OrganizationItem{}).
		Where("organization_id = ? AND item_id = ?", orgItem.OrganizationID, orgItem.ItemID).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("item already exists in organization")
	}

	return r.db.Create(orgItem).Error
}

// DeleteOrganizationItem 删除组织项目关联
func (r *OrganizationRepository) DeleteOrganizationItem(ctx context.Context, orgId string, itemId string) error {
	return r.db.Where("organization_id = ? AND item_id = ?", orgId, itemId).
		Delete(&model.OrganizationItem{}).Error
}

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
	return r.db.Model(&model.OrganizationItem{}).
		Where("organization_id = ? AND item_id = ?", orgId, itemId).
		Update("sort_order", sortOrder).Error
}
