package service

import (
	"encoding/json"
	"github.com/anushasankaranarayanan/query-aggregator-service/internal/consts"
	"github.com/anushasankaranarayanan/query-aggregator-service/internal/entity"
	"github.com/sirupsen/logrus"
	"net/http"
	"sort"
	"strings"
)

var l = logrus.StandardLogger()

const (
	hApplicationJson string = "application/json"
	hContentType     string = "Content-Type"
)

type QueryAggregator interface {
	AggregateResults(sortKey string, limit int) entity.QueryServiceResponse
}

type HttpClient interface {
	SendRequest(url, method string, header map[string][]string) ([]byte, int, error)
}

type queryAggregator struct {
	httpClient HttpClient
}

func NewQueryAggregator(httpClient HttpClient) QueryAggregator {
	return &queryAggregator{httpClient: httpClient}
}

type workFlow struct {
	url   string
	msgCh chan entity.QueryServiceResponse
}

// AggregateResults - calls the http requests concurrently using goroutines and channels , sorts and limits aggregated data
func (svc queryAggregator) AggregateResults(sortKey string, limit int) entity.QueryServiceResponse {

	pipeline := dataPipeline()
	aggregatedResponse := entity.QueryServiceResponse{}

	for _, item := range pipeline {
		go callUrl(item.msgCh, item.url, svc.httpClient)
	}

	for _, item := range pipeline {
		for message := range item.msgCh {
			aggregatedResponse.Data = append(aggregatedResponse.Data, message.Data...)
		}
	}

	sortResponse(sortKey, &aggregatedResponse)
	return limitResponse(limit, aggregatedResponse)
}

func dataPipeline() []workFlow {
	return []workFlow{
		{
			url:   "https://raw.githubusercontent.com/assignment132/assignment/main/duckduckgo.json",
			msgCh: make(chan entity.QueryServiceResponse),
		},
		{
			url:   "https://raw.githubusercontent.com/assignment132/assignment/main/google.json",
			msgCh: make(chan entity.QueryServiceResponse),
		},
		{
			url:   "https://raw.githubusercontent.com/assignment132/assignment/main/wikipedia.json",
			msgCh: make(chan entity.QueryServiceResponse),
		},
	}
}

func limitResponse(limit int, aggregatedResponse entity.QueryServiceResponse) entity.QueryServiceResponse {
	chopLimit := limit
	if limit > len(aggregatedResponse.Data) {
		chopLimit = len(aggregatedResponse.Data)
		l.Infof("Limiting response data from %d to %d", limit, chopLimit)
	}
	l.Tracef("Limiting response data to %d", chopLimit)
	return entity.QueryServiceResponse{Data: aggregatedResponse.Data[:chopLimit]}
}

func sortResponse(sortKey string, aggregatedResponse *entity.QueryServiceResponse) {
	if strings.ToLower(sortKey) == consts.RelevanceScore {
		sort.Slice(aggregatedResponse.Data, func(i, j int) bool {
			return aggregatedResponse.Data[i].RelevanceScore < aggregatedResponse.Data[j].RelevanceScore
		})
	}
	if strings.ToLower(sortKey) == consts.Views {
		sort.Slice(aggregatedResponse.Data, func(i, j int) bool {
			return aggregatedResponse.Data[i].Views < aggregatedResponse.Data[j].Views
		})
	}
}

func callUrl(ch chan entity.QueryServiceResponse, url string, httpClient HttpClient) {
	defer close(ch)
	resp := entity.QueryServiceResponse{}
	bytes, err := request(url, httpClient)
	if err != nil {
		l.Errorf("error calling URL: %s. Details: %s", url, err.Error())
		ch <- resp
	}
	if err = json.Unmarshal(bytes, &resp); err != nil {
		l.Errorf("error unmarshalling http response from url %s.Details: %s ", url, err.Error())
		ch <- resp
	}
	l.Tracef("Received response %v from url %s", resp, url)
	ch <- resp
}

func request(url string, httpClient HttpClient) ([]byte, error) {
	headers := make(map[string][]string)
	headers[hContentType] = []string{hApplicationJson}

	body, _, err := httpClient.SendRequest(url, http.MethodGet, headers)

	if err != nil {
		return nil, err
	}
	return body, nil
}
