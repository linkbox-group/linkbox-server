package delivery

import (
	"github.com/google/wire"
	"github.com/linkbox-group/linkbox-server/content/internal/acl"
)

var ProviderSet = wire.NewSet(NewContentDelivery)

type ContentDelivery struct {
	s acl.UserServiceItf
}

func NewContentDelivery(s acl.UserServiceItf) *ContentDelivery {
	return &ContentDelivery{s: s}
}
