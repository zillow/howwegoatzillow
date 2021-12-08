package main

import (
	"github.com/zillow/howwegoatzillow/libs/server"
	mock_db "github.com/zillow/howwegoatzillow/mocks/db"
	mock_kafka "github.com/zillow/howwegoatzillow/mocks/kafka"
)

type ServerTestable struct {
	Server     *server.Server
	DBProvider *mock_db.MockProvider
	KProvider  *mock_kafka.MockClient
	KWriter    *mock_kafka.MockWriter
}
