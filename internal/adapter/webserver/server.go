package webserver

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"

	"github.com/anushasankaranarayanan/query-aggregator-service/internal/service"
)

type Server struct {
	Services Services
}

type Services struct {
	QueryAggregator service.QueryAggregator
}

func NewServer(services Services) *Server {
	return &Server{
		Services: services,
	}
}

// Run -initializes the routes and starts the server
func (s *Server) Run() error {
	l := logrus.StandardLogger()
	var err error

	srv := &http.Server{Addr: ":" + os.Getenv("SERVER_PORT")}

	err = s.Routes()
	if err != nil {
		return err
	}

	fmt.Println(os.Getenv("NAME"))
	l.Info("Starting server on port " + os.Getenv("SERVER_PORT"))
	if err = srv.ListenAndServe(); err != nil {
		l.Error("Httpserver: ListenAndServe() error: " + err.Error())
	}

	return err
}
