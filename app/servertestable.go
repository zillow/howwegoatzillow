package main

import (
	"github.com/zillow/howegoatzillow/libs/server"
	mock_db "github.com/zillow/howegoatzillow/mocks/db"
	mock_kafka "github.com/zillow/howegoatzillow/mocks/kafka"
)

type ServerTestable struct {
	Server     *server.Server
	DBProvider *mock_db.MockProvider
	KProvider  *mock_kafka.MockClient
	KWriter    *mock_kafka.MockWriter
}
