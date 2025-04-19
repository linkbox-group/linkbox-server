package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/linkbox-group/linkbox-server/model"
)

var (
	ErrOrganizationNotFound       = errors.New("organization not found")
	ErrUpdateOrganizationFailed   = errors.New("update organization failed")
	ErrDeleteOrganizationFailed   = errors.New("delete organization failed")
	ErrGetUserOrganizationsFailed = errors.New("get user organizations failed")
	ErrMoveOrganizationFailed     = errors.New("move organization failed")
)

// GetOrganizationService 获取组织信息
func (s *OrganizationService) GetOrganizationService(ctx context.Context, id string, userId string) (*model.Organization, error) {
	organization, err := s.repo.GetOrganization(ctx, id, userId)
	if err != nil {
		return nil, fmt.Errorf("%w:%w", ErrOrganizationNotFound, err)
	}
	return organization, nil
}

// UpdateOrganizationService 更新组织信息
func (s *OrganizationService) UpdateOrganizationService(ctx context.Context, organization *model.Organization) (*model.Organization, error) {
	updatedOrg, err := s.repo.UpdateOrganization(ctx, organization)
	if err != nil {
		return nil, fmt.Errorf("%w:%w", ErrUpdateOrganizationFailed, err)
	}
	return updatedOrg, nil
}

// DeleteOrganizationService 删除组织
func (s *OrganizationService) DeleteOrganizationService(ctx context.Context, id string, userId string, cascade bool) error {
	err := s.repo.DeleteOrganization(ctx, id, userId, cascade)
	if err != nil {
		return fmt.Errorf("%w:%w", ErrDeleteOrganizationFailed, err)
	}
	return nil
}

// GetUserOrganizationsService 获取用户的所有组织
func (s *OrganizationService) GetUserOrganizationsService(ctx context.Context, userId string) ([]*model.Organization, error) {
	organizations, err := s.repo.GetUserOrganizations(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("%w:%w", ErrGetUserOrganizationsFailed, err)
	}
	return organizations, nil
}

// MoveOrganizationService 移动组织
func (s *OrganizationService) MoveOrganizationService(ctx context.Context, id string, userId string, newParentCode string) error {
	err := s.repo.MoveOrganization(ctx, id, userId, newParentCode)
	if err != nil {
		return fmt.Errorf("%w:%w", ErrMoveOrganizationFailed, err)
	}
	return nil
}
