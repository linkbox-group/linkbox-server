package mysql_repository

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/google/wire"
	"github.com/linkbox-group/linkbox-server/item/internal/acl"
	"github.com/linkbox-group/linkbox-server/rpc-gen/model"
)

var ProviderSet = wire.NewSet(wire.Bind(new(acl.UserRepositoryItf), new(*EsRepository)), NewEsRepository)

type EsRepository struct {
	es *elasticsearch.TypedClient
}

func NewEsRepository(es *elasticsearch.TypedClient) *EsRepository {
	return &EsRepository{es: es}
}

func (r *EsRepository) SearchItems(ctx context.Context, UserID string, query string, itemType model.ItemType) (items []model.Item, count int, err error) {
	
	panic("implement me")
}
