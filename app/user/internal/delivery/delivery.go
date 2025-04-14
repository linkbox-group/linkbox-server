package delivery

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewUserDelivery)

type UserDelivery struct {
}

func NewUserDelivery() *UserDelivery {
	return &UserDelivery{}
}
