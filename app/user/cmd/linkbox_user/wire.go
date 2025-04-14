//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/linkbox-group/linkbox-server/user/internal/delivery"
)

func NewUserHandler() *delivery.UserDelivery {
	wire.Build(delivery.ProviderSet)
	return &delivery.UserDelivery{}
}
