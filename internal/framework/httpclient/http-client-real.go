//go:build real

package httpclient

import (
	"crypto/tls"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

type HttpClient struct {
	Name       string
	httpClient *retryablehttp.Client
}

// NewHttpClient - initializes the http connections
func NewHttpClient() *HttpClient {
	l := logrus.StandardLogger()

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	minWaitTime, _ := time.ParseDuration(os.Getenv("HTTP_RETRY_MIN_WAIT"))
	l.Info(fmt.Sprintf("HTTP RetryWaitMin set to %.0f seconds.", minWaitTime.Seconds()))
	maxWaitTime, _ := time.ParseDuration(os.Getenv("HTTP_RETRY_MAX_WAIT"))
	l.Info(fmt.Sprintf("HTTP RetryWaitMin set to %.0f seconds.", maxWaitTime.Seconds()))
	retryMax, _ := strconv.Atoi(os.Getenv("HTTP_MAX_RETRIES"))
	l.Info(fmt.Sprintf("HTTP Maximum retries set to %d.", retryMax))

	client := retryablehttp.NewClient()
	client.HTTPClient = &http.Client{Transport: tr}
	client.RetryWaitMin = minWaitTime
	client.RetryWaitMax = maxWaitTime
	client.RetryMax = retryMax
	return &HttpClient{Name: "HttpClient", httpClient: client}

}
