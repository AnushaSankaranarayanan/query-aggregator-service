//go:build fake

package httpclient

import (
	"bytes"
	"errors"
	"github.com/hashicorp/go-retryablehttp"
	"io"
	"net/http"
)

type HttpClient struct {
	Name       string
	httpClient *FakeHttpClient
}

type FakeHttpClient struct {
	Flag            string
	ResponsePayload []byte
}

func NewFakeHttpClient(flag string, responsePayload []byte) *HttpClient {
	return &HttpClient{Name: "FakeHttpClient", httpClient: &FakeHttpClient{Flag: flag, ResponsePayload: responsePayload}}
}

// Do - fake implementation. Will be injected using appropriate build tags
func (c *FakeHttpClient) Do(_ *retryablehttp.Request) (*http.Response, error) {
	r := io.NopCloser(bytes.NewReader(c.ResponsePayload))

	if c.Flag == "forced-http-error" {
		return nil, errors.New("forced http error")
	}

	if c.Flag == "forced-not-ok-response" {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       r,
		}, nil
	}
	if c.Flag == "forced-nil-response" {
		return nil, nil
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       r,
	}, nil
}
