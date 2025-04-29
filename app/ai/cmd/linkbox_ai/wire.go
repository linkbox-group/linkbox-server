//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/linkbox-group/linkbox-server/ai/internal/core"
	"github.com/linkbox-group/linkbox-server/ai/internal/delivery"
	"github.com/linkbox-group/linkbox-server/ai/internal/repository"
	"github.com/linkbox-group/linkbox-server/ai/internal/service"
)

func NewAiHandler() *delivery.AiDelivery {
	wire.Build(delivery.ProviderSet, service.ProviderSet, repository.ProviderSet, core.ProviderSet)
	return &delivery.AiDelivery{}

}
