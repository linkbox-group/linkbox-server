package service

import (
	"context"
	"github.com/linkbox-group/linkbox-server/model"
)

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
