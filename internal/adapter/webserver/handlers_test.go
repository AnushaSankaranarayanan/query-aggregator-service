//go:build fake

package webserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anushasankaranarayanan/query-aggregator-service/internal/framework/httpclient"
	"github.com/anushasankaranarayanan/query-aggregator-service/internal/service"
	"github.com/gin-gonic/gin"
)

var (
	fakeResponseJson = `{
		  "data": [
			{
			  "url": "www.test.com/1",
			  "views": 1000,
              "relevanceScore": 0.2
			},
			{
			  "url": "www.test.com/2",
			  "views": 2000,
			  "relevanceScore": 0.1
			}
		  ],
		  "count": 2
		}`
)

func TestHandlers(t *testing.T) {

	tests := []struct {
		name, httpResponse, url string
		expectedResponseCode    int
	}{
		{
			name:                 "QueryHandler should fail(invalid/empty sortKey)",
			httpResponse:         fakeResponseJson,
			url:                  "https://fakeurl.com?limit=2",
			expectedResponseCode: http.StatusBadRequest,
		},
		{
			name:                 "QueryHandler should fail(invalid limit)",
			httpResponse:         fakeResponseJson,
			url:                  "https://fakeurl.com?sortKey=relevanceScore&limit=bla",
			expectedResponseCode: http.StatusBadRequest,
		},
		{
			name:                 "QueryHandler should fail(limit out of range)",
			httpResponse:         fakeResponseJson,
			url:                  "https://fakeurl.com?sortKey=relevanceScore&limit=201",
			expectedResponseCode: http.StatusBadRequest,
		},
		{
			name:                 "QueryHandler should pass(sorted by relevanceScore)",
			httpResponse:         fakeResponseJson,
			url:                  "https://fakeurl.com?sortKey=relevanceScore&limit=2",
			expectedResponseCode: http.StatusOK,
		},
		{
			name:                 "QueryHandler should pass(sorted by views)",
			httpResponse:         fakeResponseJson,
			url:                  "https://fakeurl.com?sortKey=views&limit=2",
			expectedResponseCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			rr := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, test.url, nil)

			httpclient := httpclient.NewFakeHttpClient("", []byte(test.httpResponse))
			service := service.NewQueryAggregator(httpclient)
			server := NewServer(Services{QueryAggregator: service})

			c, _ := gin.CreateTestContext(rr)
			c.Request = req

			server.QueryHandler(c)

			if rr.Code != test.expectedResponseCode {
				t.Errorf("Test %s returned with incorrect status code - got (%d) wanted (%d)", test.name, rr.Code, test.expectedResponseCode)
			}

		})
	}
}
