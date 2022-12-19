//go:build fake

package service

import (
	"github.com/anushasankaranarayanan/query-aggregator-service/internal/consts"
	"reflect"
	"testing"

	"github.com/anushasankaranarayanan/query-aggregator-service/internal/entity"
	"github.com/anushasankaranarayanan/query-aggregator-service/internal/framework/httpclient"
)

var (
	responseSortedByViews = entity.QueryServiceResponse{Data: []entity.QueryServiceData{
		{"www.test.com/1", 1000, 0.2},
	}}

	responseSortedByScore = entity.QueryServiceResponse{Data: []entity.QueryServiceData{
		{"www.test.com/2", 2000, 0.1},
	}}

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

func TestServices(t *testing.T) {
	tests := []struct {
		name, httpResponse, httpFlag, sortKey string
		limit                                 int
		expectedResponse                      entity.QueryServiceResponse
	}{
		{
			name:             "AggregateResults - sorted by relevanceScore",
			httpResponse:     fakeResponseJson,
			httpFlag:         "",
			sortKey:          consts.RelevanceScore,
			limit:            1,
			expectedResponse: responseSortedByScore,
		},
		{
			name:             "AggregateResults - sorted by views",
			httpResponse:     fakeResponseJson,
			httpFlag:         "",
			sortKey:          consts.Views,
			limit:            1,
			expectedResponse: responseSortedByViews,
		},
		{
			name:             "AggregateResults - forced http error",
			httpResponse:     fakeResponseJson,
			httpFlag:         "forced-http-error",
			sortKey:          consts.Views,
			limit:            1,
			expectedResponse: entity.QueryServiceResponse{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			httpclient := httpclient.NewFakeHttpClient(test.httpFlag, []byte(test.httpResponse))
			svc := NewQueryAggregator(httpclient)
			actualResponse := svc.AggregateResults(test.sortKey, test.limit)

			if !reflect.DeepEqual(actualResponse, test.expectedResponse) {
				t.Errorf("Test  %s returned with an unexpected response - got (%v) wanted (%v)", test.name, actualResponse, test.expectedResponse)
			}

		})
	}
}
