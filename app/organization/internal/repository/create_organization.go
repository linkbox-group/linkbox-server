package repository

import (
	"context"
	"github.com/linkbox-group/linkbox-server/model"
)

func (r *OrganizationRepository) CreateOrganization(ctx context.Context, organization *model.Organization) (err error) {
	err = r.db.Create(organization).Error

	return err

}
