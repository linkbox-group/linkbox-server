package delivery

import (
	"context"
	"github.com/google/wire"
	"github.com/linkbox-group/linkbox-server/organization/internal/acl"
	"github.com/linkbox-group/linkbox-server/rpc-gen/organization"
)

var ProviderSet = wire.NewSet(NewOrganizationDelivery)

type OrganizationDelivery struct {
	service acl.OrganizationServiceItf
}

func (d *OrganizationDelivery) GetDefaultOrgID(ctx context.Context, req *organization.GetDefaultOrgIDReq) (res *organization.GetDefaultOrgIDResp, err error) {
	id, err := d.service.GetDefaultOrgID(ctx, req.Code, req.UserId)
	if err != nil {
		return nil, err
	}
	res = &organization.GetDefaultOrgIDResp{
		Id: id,
	}
	return res, nil
}

func NewOrganizationDelivery(service acl.OrganizationServiceItf) *OrganizationDelivery {
	return &OrganizationDelivery{
		service: service,
	}
}
