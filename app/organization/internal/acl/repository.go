package acl

import (
	"context"

	"github.com/linkbox-group/linkbox-server/model"
)

type OrganizationRepositoryItf interface {
	CreateOrganization(ctx context.Context, organization *model.Organization) (err error)
	GetOrganization(ctx context.Context, id string, userId string) (*model.Organization, error)
	UpdateOrganization(ctx context.Context, organization *model.Organization) (*model.Organization, error)
	DeleteOrganization(ctx context.Context, id string, userId string, cascade bool) error
	GetUserOrganizations(ctx context.Context, userId string) ([]*model.Organization, error)
	MoveOrganization(ctx context.Context, id string, userId string, newParentCode string) error

	GetOrganizationItems(ctx context.Context, orgId string) ([]*model.Item, error)
	UpdateOrganizationItemSort(ctx context.Context, orgId string, itemId string, sortOrder int) error
}
