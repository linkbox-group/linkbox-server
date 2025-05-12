package service

import (
	"context"
	itemmodel "github.com/linkbox-group/linkbox-server/rpc-gen/model"

	"github.com/google/wire"
	"github.com/linkbox-group/linkbox-server/rpc-gen/common/pagination"

	"github.com/linkbox-group/linkbox-server/item/internal/acl"
	"github.com/linkbox-group/linkbox-server/model"
)

var ProviderSet = wire.NewSet(wire.Bind(new(acl.UserServiceItf), new(*Service)), NewItemService)

var _ acl.UserServiceItf = &Service{}

type Service struct {
	repo   acl.UserRepositoryItf
	esRepo acl.EsRepositoryItf
}

func NewItemService(r acl.UserRepositoryItf, esr acl.EsRepositoryItf) *Service {
	return &Service{
		repo:   r,
		esRepo: esr,
	}
}
func (s *Service) CreateItem(ctx context.Context, item *model.Item) error {

	return s.repo.CreateItem(ctx, item)
}

func (s *Service) GetItem(ctx context.Context, item *model.Item) error {
	return s.repo.GetItem(ctx, item)
}

func (s *Service) UpdateItem(ctx context.Context, item *model.Item) error {
	return s.repo.UpdateItem(ctx, item)
}

func (s *Service) DeleteItem(ctx context.Context, item *model.Item) error {
	return s.repo.DeleteItem(ctx, item)
}
func (s *Service) GetItemsByTags(ctx context.Context, userID string, tagNames []string, pagination *pagination.PaginationRequest) ([]model.Item, int, error) {
	return s.repo.GetItemsByTags(ctx, userID, tagNames, pagination)
}
func (s *Service) GetItemsByOrganization(ctx context.Context, userID string, organizationID string, pageNum int, pageSize int) ([]model.Item, int, error) {
	return s.repo.GetItemsByOrganization(ctx, userID, organizationID, pageNum, pageSize)

}
func (s *Service) SearchItems(ctx context.Context, userID string, query string, itemType itemmodel.ItemType, pageNum int, pageSize int) ([]model.Item, int, error) {
	return s.esRepo.SearchItems(ctx, userID, query, itemType, pageNum, pageSize)
}
func (s *Service) RecoverItemsBatch(ctx context.Context, userID string, ids []string) (err error) {
	return s.repo.RecoverItemsBatch(ctx, userID, ids)
}
func (s *Service) DeleteItemsBatch(ctx context.Context, userID string, ids []string) (err error) {
	return s.repo.DeleteItemsBatch(ctx, userID, ids)
}
func (s *Service) GetDeletedItems(ctx context.Context, userID string, pagination *pagination.PaginationRequest) ([]*model.Item, int, error) {
	return s.repo.GetDeletedItems(ctx, userID, pagination)
}
