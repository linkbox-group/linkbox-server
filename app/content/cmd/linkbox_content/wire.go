//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/linkbox-group/linkbox-server/content/internal/delivery"
)

func NewContentHandler() *delivery.ContentDelivery {
	wire.Build(delivery.ProviderSet)
	return &delivery.ContentDelivery{}

}
