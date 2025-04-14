//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/linkbox-group/linkbox-server/tag/internal/core"
	"github.com/linkbox-group/linkbox-server/tag/internal/delivery"
	"github.com/linkbox-group/linkbox-server/tag/internal/repository"
	"github.com/linkbox-group/linkbox-server/tag/internal/service"
)

func NewTagHandler() *delivery.TagDelivery {
	wire.Build(delivery.ProviderSet, service.ProviderSet, repository.ProviderSet, core.ProviderSet)
	return &delivery.TagDelivery{}

}
