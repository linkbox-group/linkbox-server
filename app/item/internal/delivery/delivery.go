package delivery

import (
	"github.com/google/wire"
	"github.com/linkbox-group/linkbox-server/item/internal/acl"
)

var ProviderSet = wire.NewSet(NewItemDelivery)

type ItemDelivery struct {
	s acl.UserServiceItf
}

func NewItemDelivery(s acl.UserServiceItf) *ItemDelivery {
	return &ItemDelivery{s: s}
}
