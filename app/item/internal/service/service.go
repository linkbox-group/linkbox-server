package service

import (
	"context"

	"github.com/google/wire"
	"github.com/linkbox-group/linkbox-server/rpc-gen/common/pagination"

	"github.com/linkbox-group/linkbox-server/item/internal/acl"
	"github.com/linkbox-group/linkbox-server/model"
)

var ProviderSet = wire.NewSet(wire.Bind(new(acl.UserServiceItf), new(*Service)), NewItemService)

var _ acl.UserServiceItf = &Service{}

type Service struct {
	Repo acl.UserRepositoryItf
}

func NewItemService(r acl.UserRepositoryItf) *Service {
	return &Service{
		Repo: r,
	}
}
func (s *Service) CreateItem(ctx context.Context, item *model.Item) error {

	return s.Repo.CreateItem(ctx, item)
}

func (s *Service) GetItem(ctx context.Context, item *model.Item) error {
	return s.Repo.GetItem(ctx, item)
}

func (s *Service) UpdateItem(ctx context.Context, item *model.Item) error {
	return s.Repo.UpdateItem(ctx, item)
}

func (s *Service) DeleteItem(ctx context.Context, item *model.Item) error {
	return s.Repo.DeleteItem(ctx, item)
}
func (s *Service) GetItemsByTags(ctx context.Context, userID string, tagNames []string, pagination *pagination.PaginationRequest) ([]model.Item, int, error) {
	return s.Repo.GetItemsByTags(ctx, userID, tagNames, pagination)
}
func (s *Service) GetItemsByOrganization(ctx context.Context, userID string, organizationID string, pageNum int, pageSize int) ([]model.Item, int, error) {
	return s.Repo.GetItemsByOrganization(ctx, userID, organizationID, pageNum, pageSize)

}
func (s *Service) SearchItems(ctx context.Context, userID string, query string, pageNum int, pageSize int) ([]model.Item, int, error) {
	return s.Repo.SearchItemsByTitle(ctx, userID, query, pageNum, pageSize)
}
func (s *Service) RecoverItemsBatch(ctx context.Context, userID string, ids []string) (err error) {
	return s.Repo.RecoverItemsBatch(ctx, userID, ids)
}
func (s *Service) DeleteItemsBatch(ctx context.Context, userID string, ids []string) (err error) {
	return s.Repo.DeleteItemsBatch(ctx, userID, ids)
}
