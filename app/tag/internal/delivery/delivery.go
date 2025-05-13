package delivery

import (
	"context"
	"github.com/google/wire"
	"github.com/linkbox-group/linkbox-server/rpc-gen/tag"
	"github.com/linkbox-group/linkbox-server/tag/internal/acl"
)

var ProviderSet = wire.NewSet(NewTagDelivery)

type TagDelivery struct {
	service acl.TagServiceItf
}

func (s *TagDelivery) RemoveTagsFromItems(ctx context.Context, req *tag.RemoveTagsFromItemsRequest) (res *tag.RemoveTagsFromItemsResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *TagDelivery) GetItemTags(ctx context.Context, req *tag.GetItemTagsRequest) (res *tag.GetItemTagsResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func NewTagDelivery(service acl.TagServiceItf) *TagDelivery {
	return &TagDelivery{
		service: service,
	}
}
