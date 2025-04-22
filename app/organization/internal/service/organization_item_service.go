package service

import (
	"context"
	"github.com/linkbox-group/linkbox-server/model"
)

// AddItemsToOrganizationService 添加项目到组织
func (s *OrganizationService) AddItemsToOrganizationService(ctx context.Context, orgId string, userId string, itemIds []string) error {
	// 检查组织是否存在
	org, err := s.repo.GetOrganization(ctx, orgId, userId)
	if err != nil {
		return err
	}

	// 创建组织项目关联
	for _, itemId := range itemIds {
		orgItem := &model.OrganizationItem{
			OrganizationID: orgId,
			ItemID:         itemId,
			SortOrder:      0, // 默认排序
		}
		err = s.repo.CreateOrganizationItem(ctx, orgItem)
		if err != nil {
			return err
		}
	}

	// 更新组织项目数
	org.ItemsCount += uint32(len(itemIds))
	_, err = s.repo.UpdateOrganization(ctx, org)
	if err != nil {
		return err
	}

	return nil
}

// RemoveItemsFromOrganizationService 从组织移除项目
func (s *OrganizationService) RemoveItemsFromOrganizationService(ctx context.Context, orgId string, userId string, itemIds []string) error {
	// 检查组织是否存在
	org, err := s.repo.GetOrganization(ctx, orgId, userId)
	if err != nil {
		return err
	}

	// 删除组织项目关联
	for _, itemId := range itemIds {
		err = s.repo.DeleteOrganizationItem(ctx, orgId, itemId)
		if err != nil {
			return err
		}
	}

	// 更新组织项目数
	org.ItemsCount -= uint32(len(itemIds))
	_, err = s.repo.UpdateOrganization(ctx, org)
	if err != nil {
		return err
	}

	return nil
}

// GetOrganizationItemsService 获取组织下的所有项目
func (s *OrganizationService) GetOrganizationItemsService(ctx context.Context, orgId string, userId string) ([]*model.Item, error) {
	// 检查组织是否存在
	_, err := s.repo.GetOrganization(ctx, orgId, userId)
	if err != nil {
		return nil, err
	}

	// 获取组织项目
	return s.repo.GetOrganizationItems(ctx, orgId)
}

// ReorderOrganizationItemsService 重新排序组织项目
func (s *OrganizationService) ReorderOrganizationItemsService(ctx context.Context, orgId string, userId string, itemOrders map[string]int) error {
	// 检查组织是否存在
	_, err := s.repo.GetOrganization(ctx, orgId, userId)
	if err != nil {
		return err
	}

	// 更新项目排序
	for itemId, sortOrder := range itemOrders {
		err = s.repo.UpdateOrganizationItemSort(ctx, orgId, itemId, sortOrder)
		if err != nil {
			return err
		}
	}

	return nil
}
