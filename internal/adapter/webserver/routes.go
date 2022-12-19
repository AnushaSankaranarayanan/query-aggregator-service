package webserver

import (
	"github.com/anushasankaranarayanan/query-aggregator-service/internal/adapter/webserver/probes"
	"github.com/anushasankaranarayanan/query-aggregator-service/internal/adapter/webserver/swagger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) Routes() error {
	r := gin.New()

	p := r.Group("/api/v1/probes")
	p.GET("/liveness", probes.Liveness)

	o := r.Group("/api/v1/openapi")
	o.GET("/", swagger.Build)
	o.GET("/:resource", swagger.Build)

	api := r.Group("/api/v1")
	api.Use(gin.Logger())

	api.GET("/query", s.QueryHandler)

	http.Handle("/", r)

	return nil
}
