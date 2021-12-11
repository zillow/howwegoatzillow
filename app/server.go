package main

import (
	"net/http"

	"github.com/zillow/howwegoatzillow/libs/db"
	zhttp "github.com/zillow/howwegoatzillow/libs/http"
	"github.com/zillow/howwegoatzillow/libs/kafka"
	"github.com/zillow/howwegoatzillow/libs/server"
)

type MyServer struct{ *server.Server }

// This is a very crude representation of what each application looks like.
// Everything from here and underneath is this application domain and should be well tested.
func NewMyServer(server *server.Server, service *MyService) *MyServer {
	handleRequest := func(w http.ResponseWriter, r *http.Request) {
		httpClient := service.HTTPClientProvider.GetWrappedClient(service.HTTPConfig)
		_, _ = httpClient.Get("http://hello.com/")

		sqlx, _ := service.DBProvider.Get(r.Context(), service.DBConfig)
		_, _ = sqlx.ExecContext(r.Context(), "select * from table")

		kw, _ := service.KafkaClient.Writer(r.Context(), service.KafkaConfig)
		_, _ = kw.Write(r.Context(), "apple", []byte("message"))
		w.WriteHeader(http.StatusNoContent)
	}

	server.Router.HandleFunc("/", handleRequest)
	return &MyServer{server}
}

type MyService struct {
	HTTPConfig         zhttp.Config
	HTTPClientProvider zhttp.Provider

	DBConfig   db.Config
	DBProvider db.Provider

	KafkaConfig kafka.Config
	KafkaClient kafka.Client
}
