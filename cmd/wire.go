//go:build wireinject
// +build wireinject

package cmd

import (
	"context"
	"qqlx/base/app"
	"qqlx/base/handler"
	"qqlx/base/middleware"
	"qqlx/base/server"
	"qqlx/base/validator"
	"qqlx/controller"
	"qqlx/router"
	"qqlx/service"
	"qqlx/store"

	"github.com/google/wire"
)

func InitApplication(ctx context.Context) (*app.Application, func(), error) {
	//panic(wire.Build(
	//	server.NewHttpServer,
	//	store.ProviderStore,
	//	service.ProviderService,
	//	validator.ProviderValidator,
	//	handler.ProviderHandler,
	//	controller.ProviderContr,
	//	middleware.ProviderMiddleware,
	//	router.ProviderRouter,
	//	app.NewApplication,
	//))
	wire.Build(
		server.NewHttpServer,
		store.ProviderStore,
		service.ProviderService,
		validator.ProviderValidator,
		handler.ProviderHandler,
		controller.ProviderContr,
		middleware.ProviderMiddleware,
		router.ProviderRouter,
		app.NewApplication,
	)
	return nil, nil, nil
}
