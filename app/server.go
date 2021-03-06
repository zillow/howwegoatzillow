package main

import (
	"net/http"

	"github.com/zillow/howwegoatzillow/libs/db"
	zhttp "github.com/zillow/howwegoatzillow/libs/http"
	"github.com/zillow/howwegoatzillow/libs/kafka"
	"github.com/zillow/howwegoatzillow/libs/server"
)

// This is a very crude representation of what each application looks like.
// Everything from here and underneath is this application domain and should be well tested.
func NewServer(service MyService) *server.Server {
	handleRequest := func(w http.ResponseWriter, r *http.Request) {
		httpClient := service.HTTPClientProvider.GetWrappedClient(service.HTTPConfig)
		_, _ = httpClient.Get("http://hello.com/")

		sqlx, _ := service.DBProvider.Get(r.Context(), service.DBConfig)
		_, _ = sqlx.ExecContext(r.Context(), "select * from table")

		kw, _ := service.KafkaClient.Writer(r.Context(), service.KafkaConfig)
		_, _ = kw.Write(r.Context(), "apple", []byte("message"))
		w.WriteHeader(http.StatusNoContent)
	}

	s := service.ServerFactory.Create()
	s.Router.HandleFunc("/", handleRequest)
	return s
}

type MyService struct {
	ServerFactory server.Factory

	HTTPConfig         zhttp.Config
	HTTPClientProvider zhttp.Provider

	DBConfig   db.Config
	DBProvider db.Provider

	KafkaConfig kafka.Config
	KafkaClient kafka.Client
}
