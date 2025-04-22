package acl

import (
	"context"

	"github.com/linkbox-group/linkbox-server/model"
)

type OrganizationServiceItf interface {
	CreateOrganizationService(ctx context.Context, organization *model.Organization) (err error)
	GetOrganizationService(ctx context.Context, id string, userId string) (*model.Organization, error)
	UpdateOrganizationService(ctx context.Context, organization *model.Organization) (*model.Organization, error)
	DeleteOrganizationService(ctx context.Context, id string, userId string, cascade bool) error
	GetUserOrganizationsService(ctx context.Context, userId string) ([]*model.Organization, error)
	MoveOrganizationService(ctx context.Context, id string, userId string, newParentCode string) error

	// Organization Item Operations
	AddItemsToOrganizationService(ctx context.Context, orgId string, userId string, itemIds []string) error
	RemoveItemsFromOrganizationService(ctx context.Context, orgId string, userId string, itemIds []string) error
	GetOrganizationItemsService(ctx context.Context, orgId string, userId string) ([]*model.Item, error)
	ReorderOrganizationItemsService(ctx context.Context, orgId string, userId string, itemOrders map[string]int) error
}
