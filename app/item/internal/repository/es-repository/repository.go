package es_repository

import (
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/google/wire"
	"github.com/linkbox-group/linkbox-server/item/internal/acl"
	"github.com/linkbox-group/linkbox-server/item/pkg/log"
	"github.com/linkbox-group/linkbox-server/model"
	itemmodel "github.com/linkbox-group/linkbox-server/rpc-gen/model"
	"strconv"
)

var ProviderSet = wire.NewSet(wire.Bind(new(acl.EsRepositoryItf), new(*EsRepository)), NewEsRepository)
var _ acl.EsRepositoryItf = &EsRepository{}

type EsRepository struct {
	es *elasticsearch.TypedClient
}

const (
	INDEX         = "item"
	FIELD_NOTE    = "note"
	FIELD_TYPE    = "item_type"
	FIELD_TITLE   = "title"
	FIELD_USER_ID = "user_id"
)

func NewEsRepository(es *elasticsearch.TypedClient) *EsRepository {
	return &EsRepository{es: es}
}

func (r *EsRepository) SearchItems(ctx context.Context, UserID string, query string, itemType itemmodel.ItemType, pageNum int, pageSize int) (items []model.Item, count int, err error) {
	queryCondition := types.Query{}
	switch itemType {
	case itemmodel.ItemType_NOTE:
		{
			queryCondition = types.Query{
				Match: map[string]types.MatchQuery{
					FIELD_NOTE: {Query: query},
				},
			}
		}
	case itemmodel.ItemType_LINK:
		{
			queryCondition = types.Query{
				Match: map[string]types.MatchQuery{
					FIELD_TITLE: {Query: query},
				},
			}
		}
	}
	boolQuery := &types.BoolQuery{
		Must: []types.Query{
			{
				Match: map[string]types.MatchQuery{
					FIELD_USER_ID: {Query: UserID},
				},
			},
			queryCondition,
			{
				Match: map[string]types.MatchQuery{
					FIELD_TYPE: {Query: strconv.Itoa(int(itemType))},
				},
			},
		},
	}
	resp, err := r.es.Search().
		Index(INDEX).
		Query(
			&types.Query{
				Bool: boolQuery,
			},
		).
		From((pageNum - 1) * pageSize).
		Size(pageSize).
		Do(ctx)
	if err != nil {
		log.Log().Error("EsRepository.SearchItems err" + err.Error())
		return nil, 0, err
	}
	count = len(resp.Hits.Hits)
	for _, hit := range resp.Hits.Hits {
		data, err := hit.Source_.MarshalJSON()
		if err != nil {
			log.Log().Error("EsRepository.SearchItems err" + err.Error())
			count--
			continue

		}
		item := model.Item{}
		err = json.Unmarshal(data, &item)
		if err != nil {
			log.Log().Error("EsRepository.SearchItems err" + err.Error())
			count--
			continue
		}
		items = append(items, item)
	}
	return items, count, nil
}
