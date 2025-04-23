package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/linkbox-group/linkbox-server/model"
	"github.com/linkbox-group/linkbox-server/model/treemodel"
	"github.com/linkbox-group/linkbox-server/organization/pkg/log"
)

var (
	ErrDbCreateTagFailed = errors.New("create tag failed")
)

func (s *OrganizationService) CreateOrganizationService(ctx context.Context, org *model.Organization) (err error) {
	if org.ParentCode != "" || org.ParentCode != string(treemodel.ROOT_ID) {
		log.Log().Debug("have parent")
		_, err = s.repo.GetOrganization(ctx, org.ParentCode, org.UserID)

	}
	if err != nil {
		log.Log().Error("未获取到父节点" + err.Error())
		return ErrDbCreateTagFailed

	}

	err = s.repo.CreateOrganization(ctx, org)
	if err != nil {
		return fmt.Errorf("%w:%w", ErrDbCreateTagFailed, err)
	}
	return nil
}
