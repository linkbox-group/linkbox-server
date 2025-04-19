package service

import (
	"context"
	"github.com/google/wire"

	"github.com/linkbox-group/linkbox-server/content/internal/acl"
	"github.com/linkbox-group/linkbox-server/model"
)

var ProviderSet = wire.NewSet(wire.Bind(new(acl.UserServiceItf), new(*Service)), NewContentService)

var _ acl.UserServiceItf = &Service{}

type Service struct {
	Repo acl.UserRepositoryItf
}

func NewContentService(r acl.UserRepositoryItf) *Service {
	return &Service{
		Repo: r,
	}
}
func (s *Service) CreateItem(ctx context.Context, item *model.Item) error {
	return s.Repo.CreateItem(ctx, item)
}
