//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/linkbox-group/linkbox-server/organization/internal/core"
	"github.com/linkbox-group/linkbox-server/organization/internal/delivery"
	"github.com/linkbox-group/linkbox-server/organization/internal/repository"
	"github.com/linkbox-group/linkbox-server/organization/internal/service"
)

func NewOrganizationHandler() *delivery.OrganizationDelivery {
	wire.Build(delivery.ProviderSet, service.ProviderSet, repository.ProviderSet, core.ProviderSet)
	return &delivery.OrganizationDelivery{}

}
