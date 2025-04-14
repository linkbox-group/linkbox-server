//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/linkbox-group/linkbox-server/tag/internal/delivery"
)

func NewTagHandler() *delivery.TagDelivery {
	wire.Build(delivery.ProviderSet)
	return &delivery.TagDelivery{}

}
