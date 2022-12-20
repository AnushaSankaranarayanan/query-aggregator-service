package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"golang.org/x/sync/errgroup"

	"github.com/anushasankaranarayanan/query-aggregator-service/internal/consts"
	"github.com/anushasankaranarayanan/query-aggregator-service/internal/entity"
	"github.com/sirupsen/logrus"
)

var l = logrus.StandardLogger()

const (
	hApplicationJson string = "application/json"
	hContentType     string = "Content-Type"
)

type QueryAggregator interface {
	AggregateResults(sortKey string, limit int) (entity.QueryServiceResponse, error)
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

// AggregateResults - calls the http servers concurrently,sorts and limits the aggregated data.Error groups are used for error handling when calling goroutines concurrently
func (svc queryAggregator) AggregateResults(sortKey string, limit int) (entity.QueryServiceResponse, error) {

	urls := []string{
		"https://raw.githubusercontent.com/assignment132/assignment/main/duckduckgo.json",
		"https://raw.githubusercontent.com/assignment132/assignment/main/google.json",
		"https://raw.githubusercontent.com/assignment132/assignment/main/wikipedia.json",
	}

	aggregatedResponse := entity.QueryServiceResponse{}
	g, _ := errgroup.WithContext(context.Background())

	//invoke concurrently using goroutines
	for _, url := range urls {
		url := url //https://golang.org/doc/faq#closures_and_goroutines
		g.Go(func() error {
			result, err := callUrl(url, svc.httpClient)
			if err == nil {
				aggregatedResponse.Data = append(aggregatedResponse.Data, result.Data...)
			}
			return err
		})
	}

	//handle errors
	if err := g.Wait(); err != nil {
		l.Error(err.Error())
		return aggregatedResponse, err
	}

	//sort and limit the response
	sortResponse(sortKey, &aggregatedResponse)
	return limitResponse(limit, aggregatedResponse), nil
}

func callUrl(url string, httpClient HttpClient) (entity.QueryServiceResponse, error) {
	resp := entity.QueryServiceResponse{}
	bytes, err := request(url, httpClient)
	if err != nil {
		return resp, errors.New(fmt.Sprintf("error calling URL: %s. Details: %s", url, err.Error()))
	}
	if err = json.Unmarshal(bytes, &resp); err != nil {
		return resp, errors.New(fmt.Sprintf("error unmarshalling http response from url %s.Details: %s ", url, err.Error()))
	}
	l.Tracef("Received response %v from url %s", resp, url)
	return resp, nil
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

func limitResponse(limit int, aggregatedResponse entity.QueryServiceResponse) entity.QueryServiceResponse {
	chopLimit := limit
	if limit > len(aggregatedResponse.Data) {
		chopLimit = len(aggregatedResponse.Data)
		l.Infof("Limiting response data from %d to %d", limit, chopLimit)
	}
	l.Tracef("Limiting response data to %d", chopLimit)
	return entity.QueryServiceResponse{Data: aggregatedResponse.Data[:chopLimit]}
}
