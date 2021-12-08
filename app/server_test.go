package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/zillow/howwegoatzillow/libs/kafka"
	"gopkg.in/h2non/gock.v1"
)

func Test_Server_UsingMocks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s, f := InitializeServerTestable(ctrl)
	defer f()

	defer gock.Off()
	// Intercept http client call
	gock.New("http://hello.com").Get("/").Reply(200)

	db, mock, _ := sqlmock.New()
	defer db.Close()
	// Intercept DB call
	mock.ExpectExec("select * from table").WillReturnResult(sqlmock.NewResult(1, 1))
	s.DBProvider.EXPECT().Get(gomock.Any(), gomock.Any()).Times(1).Return(sqlx.NewDb(db, "sqlmock"), nil)

	s.KProvider.EXPECT().Writer(gomock.Any(), gomock.Any()).Times(1).Return(s.KWriter, nil)
	// Intercept Kafka write
	s.KWriter.EXPECT().Write(gomock.Any(), "apple", []byte("message")).Times(1).Return(kafka.Response{}, nil)

	// Make request
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	s.Server.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusNoContent {
		t.Error("no content expected")
	}
}
