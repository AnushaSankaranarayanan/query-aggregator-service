//go:build real

package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/anushasankaranarayanan/query-aggregator-service/internal/adapter/webserver"
	"github.com/anushasankaranarayanan/query-aggregator-service/internal/framework/httpclient"
	"github.com/anushasankaranarayanan/query-aggregator-service/internal/service"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// main- needs no explanation :-) Initializes the connectors and services and starts the server
func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Panic!! Service Startup error: %s\\n", err.Error())
		os.Exit(1)
	}
}

func run() error {
	log := logrus.StandardLogger()
	if err := godotenv.Load(".env"); err != nil {
		log.Debugf(".env file not detected.... falling through to Kubernetes ✿✿")
	}

	httpClient := httpclient.NewHttpClient()
	if httpClient == nil {
		log.Error("cannot set up http client")
		return errors.New("cannot set up http client")
	}

	queryService := service.NewQueryAggregator(httpClient)
	if queryService == nil {
		log.Error("cannot set up services")
		return errors.New("cannot set up services")
	}

	services := webserver.Services{
		QueryAggregator: queryService,
	}

	server := webserver.NewServer(services)
	if server == nil {
		log.Error("cannot set up server")
		return errors.New("cannot set up server")
	}

	if err := server.Run(); err != nil {
		return err
	}

	return nil
}
