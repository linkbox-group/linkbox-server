//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/linkbox-group/linkbox-server/item/internal/core"
	"github.com/linkbox-group/linkbox-server/item/internal/delivery"
	"github.com/linkbox-group/linkbox-server/item/internal/repository/mysql-repository"
	"github.com/linkbox-group/linkbox-server/item/internal/service"
)

func NewItemHandler() *delivery.ItemDelivery {
	wire.Build(delivery.ProviderSet, service.ProviderSet, mysql_repository.ProviderSet, core.ProviderSet)
	return &delivery.ItemDelivery{}

}
