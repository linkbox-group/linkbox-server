package delivery

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewContentDelivery)

type ContentDelivery struct {
}

func NewContentDelivery() *ContentDelivery {
	return &ContentDelivery{}
}
