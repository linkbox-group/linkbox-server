package delivery

import (
	"github.com/google/wire"
	"github.com/linkbox-group/linkbox-server/tag/internal/acl"
)

var ProviderSet = wire.NewSet(NewTagDelivery)

type TagDelivery struct {
	service acl.TagServiceItf
}

func NewTagDelivery(service acl.TagServiceItf) *TagDelivery {
	return &TagDelivery{
		service: service,
	}
}
