//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/linkbox-group/linkbox-server/content/internal/core"
	"github.com/linkbox-group/linkbox-server/content/internal/delivery"
	"github.com/linkbox-group/linkbox-server/content/internal/repository"
	"github.com/linkbox-group/linkbox-server/content/internal/service"
)

func NewContentHandler() *delivery.ContentDelivery {
	wire.Build(delivery.ProviderSet, service.ProviderSet, repository.ProviderSet, core.ProviderSet)
	return &delivery.ContentDelivery{}

}
