package main

import (
	"github.com/opentracing/opentracing-go"
	"github.com/zillow/howegoatzillow/libs/config"
	"github.com/zillow/howegoatzillow/libs/db"
	zhttp "github.com/zillow/howegoatzillow/libs/http"
	"github.com/zillow/howegoatzillow/libs/kafka"
	"github.com/zillow/howegoatzillow/libs/logger"
	"github.com/zillow/howegoatzillow/libs/server"
	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
)

func NewServerFactory(
	config server.Config,
	logger logger.Logger,
	tracer opentracing.Tracer,
) server.Factory {
	return server.NewFactory(
		server.WithLogger(logger),
		server.WithRouter(func() server.Handler {
			return httptrace.NewServeMux()
		}))
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
