//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/linkbox-group/linkbox-server/auth/internal/core"
	"github.com/linkbox-group/linkbox-server/auth/internal/delivery"
	"github.com/linkbox-group/linkbox-server/auth/internal/repository"
	"github.com/linkbox-group/linkbox-server/auth/internal/service"
)

func NewAuthHandler() *delivery.AuthDelivery {
	wire.Build(delivery.ProviderSet, repository.ProviderSet, service.ProviderSet, core.ProviderSet)
	return &delivery.AuthDelivery{}
}
