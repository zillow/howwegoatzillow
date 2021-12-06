//go:build wireinject
// +build wireinject

package main

import (
	"github.com/golang/mock/gomock"
	"github.com/google/wire"
	"github.com/zillow/howegoatzillow/libs/config"
	"github.com/zillow/howegoatzillow/libs/db"
	zhttp "github.com/zillow/howegoatzillow/libs/http"
	"github.com/zillow/howegoatzillow/libs/kafka"
	"github.com/zillow/howegoatzillow/libs/logger"
	"github.com/zillow/howegoatzillow/libs/server"
	mock_db "github.com/zillow/howegoatzillow/mocks/db"
	mock_kafka "github.com/zillow/howegoatzillow/mocks/kafka"
)

func InitializeServer() (*server.Server, func()) {
	wire.Build(
		ZCommonSet,
		wire.Struct(new(MyService), "*"),
		NewServer,
	)
	return &server.Server{}, nil
}

func InitializeServerTestable(ctrl *gomock.Controller) (*ServerTestable, func()) {
	wire.Build(
		ZCommonMockSet,
		wire.Struct(new(MyService), "*"),
		NewServer,
		wire.Struct(new(ServerTestable), "*"),
	)
	return &ServerTestable{}, nil
}

// This is in a separate common package
var ZCommonSet = wire.NewSet(
	NewServerConfig,
	NewServerFactory,
	config.NewAppConfig,
	NewKafkaConfig,
	kafka.NewClient,
	wire.Bind(new(kafka.Logger), new(logger.Logger)),
	logger.NewLogger,
	NewTracer,
	NewDbConfig,
	db.NewProvider,
	NewHttpServiceConfig,
	zhttp.NewClientProvider,
	wire.Bind(new(zhttp.Logger), new(logger.Logger)),
	zhttp.NewLeveledLogger,
)

var ZCommonMockSet = wire.NewSet(
	NewServerConfig,
	NewServerFactory,
	config.NewAppConfig,
	NewKafkaConfig,
	logger.NewLogger,
	NewTracer,
	NewDbConfig,
	NewHttpServiceConfig,
	zhttp.NewClientProvider,
	wire.Bind(new(zhttp.Logger), new(logger.Logger)),
	zhttp.NewLeveledLogger,

	mock_kafka.NewMockClient,
	mock_kafka.NewMockWriter,
	wire.Bind(new(kafka.Client), new(*mock_kafka.MockClient)),

	mock_db.NewMockProvider,
	wire.Bind(new(db.Provider), new(*mock_db.MockProvider)),
)
