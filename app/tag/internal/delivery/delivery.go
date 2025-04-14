package delivery

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewTagDelivery)

type TagDelivery struct {
}

func NewTagDelivery() *TagDelivery {
	return &TagDelivery{}
}
