package main

import (
	"github.com/opentracing/opentracing-go"
	"github.com/zillow/howwegoatzillow/libs/config"
	"github.com/zillow/howwegoatzillow/libs/db"
	zhttp "github.com/zillow/howwegoatzillow/libs/http"
	"github.com/zillow/howwegoatzillow/libs/kafka"
	"github.com/zillow/howwegoatzillow/libs/logger"
	"github.com/zillow/howwegoatzillow/libs/server"
	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
)

func NewServer(
	config server.Config,
	logger logger.Logger,
	tracer opentracing.Tracer,
) *server.Server {
	s := server.Server{}
	s.Configure(
		server.WithServerLogger(logger),
		server.WithServerTracer(tracer),
		server.WithServerConfig(config),
		server.WithServerRouter(httptrace.NewServeMux()),
	)
	return &s
}

//NOTE NONE OF THESE CONFIGURATIONS ARE CORRECTLY POPULATED HERE.
//JUST HERE TO SHOW MOST OF INITIAL SETUP CAN BE CONFIGURATION DRIVEN

func NewServerConfig(ac *config.AppConfig) server.Config {
	return server.Config{}
}

func NewHttpServiceConfig(ac *config.AppConfig) zhttp.Config {
	return zhttp.Config{}
}

func NewKafkaConfig(ac *config.AppConfig) kafka.Config {
	return kafka.Config{}
}
func NewDbConfig(ac *config.AppConfig) db.Config {
	return db.Config{}
}
func NewTracer() opentracing.Tracer {
	return opentracing.GlobalTracer() //Create your own tracer with your addr, host, serviceName, etc.
}
