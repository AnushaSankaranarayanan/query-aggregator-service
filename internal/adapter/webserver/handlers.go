package webserver

import (
	"fmt"
	"github.com/anushasankaranarayanan/query-aggregator-service/internal/consts"
	"github.com/anushasankaranarayanan/query-aggregator-service/internal/entity"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
)

var l = logrus.StandardLogger()

// QueryHandler - queries the given URLs and returns an aggregated response sorted by sortKey parameter limited by the limit parameter
func (s *Server) QueryHandler(c *gin.Context) {

	sortKey := c.Query(consts.PSortKey)
	if !sortKeyValid(sortKey) {
		handleError(c, http.StatusBadRequest, fmt.Sprintf("QueryHandler non recognizable Parameter: %s. Received: %s Expected: %s or %s", consts.PSortKey, sortKey, consts.RelevanceScore, consts.Views))
		return
	}

	limit := c.Query(consts.PLimit)
	if !limitValid(limit) {
		handleError(c, http.StatusBadRequest, fmt.Sprintf("QueryHandler non recognizable Parameter: %s. Received: %s Expected: %s", consts.PLimit, limit, "2 to 199"))
		return
	}

	numLimit, _ := strconv.Atoi(limit)
	qryResponse, err := s.Services.QueryAggregator.AggregateResults(sortKey, numLimit)
	if err != nil {
		handleError(c, http.StatusInternalServerError, fmt.Sprintf("QueryHandler internal error. Refer to logs for further details %s", err.Error()))
		return
	}
	resp := entity.Response{
		Code:                 http.StatusOK,
		Status:               http.StatusText(http.StatusOK),
		Message:              "Success",
		QueryServiceResponse: qryResponse,
		Count:                len(qryResponse.Data),
	}
	c.JSON(http.StatusOK, resp)
}

func sortKeyValid(sortKey string) bool {
	if strings.ToLower(sortKey) == consts.RelevanceScore || strings.ToLower(sortKey) == consts.Views {
		return true
	}
	return false
}

func limitValid(limit string) bool {
	num, err := strconv.Atoi(limit)
	if err != nil {
		return false
	}
	if num <= 1 || num >= 200 {
		return false
	}
	return true
}

func handleError(c *gin.Context, code int, msg string) {
	l.Error(msg)
	resp := entity.Response{Code: code, Status: "ERROR", Message: msg}
	c.JSON(code, resp)
	return
}
