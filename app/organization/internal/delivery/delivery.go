package delivery

import (
	"github.com/google/wire"
	"github.com/linkbox-group/linkbox-server/organization/internal/acl"
)

var ProviderSet = wire.NewSet(NewOrganizationDelivery)

type OrganizationDelivery struct {
	service acl.OrganizationServiceItf
}

func NewOrganizationDelivery(service acl.OrganizationServiceItf) *OrganizationDelivery {
	return &OrganizationDelivery{
		service: service,
	}
}
