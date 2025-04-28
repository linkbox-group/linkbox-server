package repository

import (
	"context"
	"errors"
	"github.com/linkbox-group/linkbox-server/model/treemodel"
	"github.com/linkbox-group/linkbox-server/organization/pkg/log"

	"github.com/linkbox-group/linkbox-server/model"
	"gorm.io/gorm"
)

// GetOrganization 获取组织信息
func (r *OrganizationRepository) GetOrganization(ctx context.Context, id string, userId string) (*model.Organization, error) {
	var organization model.Organization
	err := r.db.Where("id = ? AND user_id = ?", id, userId).First(&organization).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("organization not found")
		}
		return nil, err
	}
	return &organization, nil
}

// UpdateOrganization 更新组织信息
func (r *OrganizationRepository) UpdateOrganization(ctx context.Context, organization *model.Organization) (*model.Organization, error) {
	// 先检查组织是否存在
	var existingOrg model.Organization
	err := r.db.Where("id = ? AND user_id = ?", organization.ID, organization.UserID).First(&existingOrg).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("organization not found")
		}
		return nil, err
	}

	// 更新组织信息
	err = r.db.Model(&existingOrg).Updates(map[string]interface{}{
		"name":        organization.Name,
		"description": organization.Description,
		"sort_order":  organization.SortOrder,
	}).Error
	if err != nil {
		return nil, err
	}

	// 重新获取更新后的组织信息
	return r.GetOrganization(ctx, organization.ID, organization.UserID)
}

// DeleteOrganization 删除组织
func (r *OrganizationRepository) DeleteOrganization(ctx context.Context, id string, userId string, cascade bool) error {
	// 开启事务
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 先检查组织是否存在
	var organization model.Organization
	err := tx.Where("id = ? AND user_id = ?", id, userId).First(&organization).Error
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("organization not found")
		}
		return err
	}

	if cascade {
		// 级联删除子组织
		var childOrgs []model.Organization
		err = tx.Where("parent_code = ?", organization.Code).Find(&childOrgs).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		for _, childOrg := range childOrgs {
			err = r.DeleteOrganization(ctx, childOrg.ID, userId, cascade)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	} else {
		// 检查是否有子组织
		var count int64
		err = tx.Model(&model.Organization{}).Where("parent_code = ?", organization.Code).Count(&count).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		if count > 0 {
			tx.Rollback()
			return errors.New("cannot delete organization with children, use cascade delete instead")
		}
	}
	// 删除组织
	err = tx.Delete(&organization).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 更新父节点的叶子状态
	if organization.ParentCode != "0" && organization.ParentCode != "" {
		err = r.treeService.UpdateTreeLeaf(model.Organization{
			TreeModel: treemodel.TreeModel{
				Code: organization.ParentCode,
			},
		})
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// GetUserOrganizations 获取用户的所有组织
func (r *OrganizationRepository) GetUserOrganizations(ctx context.Context, userId string) ([]*model.Organization, error) {
	var organizations []*model.Organization
	err := r.db.Where("user_id = ?", userId).Order("sort_order ASC").Find(&organizations).Error
	if err != nil {
		return nil, err
	}
	return organizations, nil
}

// MoveOrganization 移动组织到新的父组织下
func (r *OrganizationRepository) MoveOrganization(ctx context.Context, id string, userId string, newParentCode string) error {
	// 开启事务
	log.Log().Info(id)
	log.Log().Info(userId)
	oldOrganization, err := r.GetOrganization(ctx, id, userId)
	if err != nil {
		log.Log().Error(err.Error())
		return err
	}
	err = r.treeService.MoveNode(r.db, oldOrganization, newParentCode)
	if err != nil {
		log.Log().Error(err.Error())
		return err
	}
	return nil
}
