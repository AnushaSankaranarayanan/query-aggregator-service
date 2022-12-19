//go:build fake

package httpclient

import (
	"fmt"
	"net/http"
	"testing"
)

func TestHttpClient(t *testing.T) {

	tests := []struct {
		testName           string
		flag               string
		errorExpected      string
		URL                string
		statusCodeExpected int
		headers            map[string][]string
	}{
		{
			"Do: should pass",
			"",
			"",
			"https://fakeurl.com",
			http.StatusOK,
			map[string][]string{"TestHeader": {"TestHeaderValue"}},
		},
		{
			"Do: should fail (forced http error)",
			"forced-http-error",
			"forced http error",
			"https://fakeurl.com",
			http.StatusInternalServerError,
			nil,
		},
		{
			"Do: should fail (forced status code 500)",
			"forced-not-ok-response",
			"status code unexpected, statusCode 500. Response body: ",
			"https://fakeurl.com",
			http.StatusInternalServerError,
			nil,
		},
		{
			"Do: should fail (forced nil response)",
			"forced-nil-response",
			"error returned from URL https://fakeurl.com . Error message: <nil>",
			"https://fakeurl.com",
			http.StatusInternalServerError,
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			httpClient := NewFakeHttpClient(test.flag, []byte{})

			_, statusCode, err := httpClient.SendRequest(test.URL, http.MethodGet, test.headers)

			fmt.Println(statusCode)

			if err != nil && err.Error() != test.errorExpected {
				t.Errorf("Test %s returned with incorrect error - got (%s) wanted (%s)", test.testName, err.Error(), test.errorExpected)
			}

			if statusCode != test.statusCodeExpected {
				t.Errorf("Test %s returned with incorrect status code - got (%d) wanted (%d)", test.testName, statusCode, test.statusCodeExpected)
			}
		})
	}
}

func TestNewFakeHttpClient(t *testing.T) {
	httpClient := NewFakeHttpClient("", []byte{})

	if httpClient == nil {
		t.Error("NewFakeHttpClient should not return nil")
	}
}
