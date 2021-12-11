//go:build wireinject
// +build wireinject

package main

import (
	"github.com/golang/mock/gomock"
	"github.com/google/wire"
	"github.com/zillow/howwegoatzillow/libs/config"
	"github.com/zillow/howwegoatzillow/libs/db"
	zhttp "github.com/zillow/howwegoatzillow/libs/http"
	"github.com/zillow/howwegoatzillow/libs/kafka"
	"github.com/zillow/howwegoatzillow/libs/logger"
	mock_db "github.com/zillow/howwegoatzillow/mocks/db"
	mock_kafka "github.com/zillow/howwegoatzillow/mocks/kafka"
)

func InitializeServer() (*MyServer, func()) {
	wire.Build(
		ZCommonSet,
		wire.Struct(new(MyService), "*"),
		NewMyServer,
	)
	return &MyServer{}, nil
}

func InitializeServerTestable(ctrl *gomock.Controller) (*ServerTestable, func()) {
	wire.Build(
		ZCommonMockSet,
		wire.Struct(new(MyService), "*"),
		NewMyServer,
		wire.Struct(new(ServerTestable), "*"),
	)
	return &ServerTestable{}, nil
}

// This is in a separate common package
var ZCommonSet = wire.NewSet(
	NewServerConfig,
	NewServer,
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
	NewServer,
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
