package service

import (
	"context"
	"github.com/google/wire"
	"github.com/linkbox-group/linkbox-server/organization/internal/acl"
)

var ProviderSet = wire.NewSet(wire.Bind(new(acl.OrganizationServiceItf), new(*OrganizationService)), NewOrganizationService)
var _ acl.OrganizationServiceItf = &OrganizationService{}

type OrganizationService struct {
	repo acl.OrganizationRepositoryItf
}

func (s *OrganizationService) GetDefaultOrgID(ctx context.Context, code string, userID string) (id string, err error) {
	return s.repo.GetDefaultOrgID(ctx, code, userID)
}

func NewOrganizationService(repo acl.OrganizationRepositoryItf) *OrganizationService {
	return &OrganizationService{repo: repo}
}
