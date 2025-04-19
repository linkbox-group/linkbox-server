//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/linkbox-group/linkbox-server/user/internal/core"
	"github.com/linkbox-group/linkbox-server/user/internal/delivery"
	"github.com/linkbox-group/linkbox-server/user/internal/repository"
	"github.com/linkbox-group/linkbox-server/user/internal/service"
)

func NewUserHandler() *delivery.UserDelivery {
	panic(wire.Build(
		core.ProviderSet,
		repository.ProviderSet,
		service.ProviderSet,
		delivery.ProviderSet,
	))
}
