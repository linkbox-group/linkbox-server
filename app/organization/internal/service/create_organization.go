package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/linkbox-group/linkbox-server/model"
)

var (
	ErrDbCreateTagFailed = errors.New("create tag failed")
)

func (s *OrganizationService) CreateOrganizationService(ctx context.Context, tag *model.Organization) (err error) {
	err = s.repo.CreateOrganization(ctx, tag)
	if err != nil {
		return fmt.Errorf("%w:%w", ErrDbCreateTagFailed, err)
	}
	return nil
}
