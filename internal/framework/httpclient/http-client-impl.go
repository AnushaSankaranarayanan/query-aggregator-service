package httpclient

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
)

// SendRequest - Invokes the HTTP Request with retries. Returns an error if response is error/not expected
func (c *HttpClient) SendRequest(url, method string, header map[string][]string) ([]byte, int, error) {
	log := logrus.StandardLogger()
	var b []byte
	var err error
	var req *retryablehttp.Request
	var resp *http.Response

	req, err = retryablehttp.NewRequest(method, url, nil)

	if header != nil {
		req.Header = header
	}

	if err != nil {
		return b, http.StatusInternalServerError, fmt.Errorf("error creating http request to URL %s.Error Details: %v", url, err)
	}

	resp, err = c.httpClient.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	log.Tracef("URL called: %s", url)
	log.Tracef("HTTP Response: %v", resp)

	if resp != nil {
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(resp.Body)
			return nil, resp.StatusCode, fmt.Errorf("status code unexpected, statusCode %d. Response body: %v", resp.StatusCode, string(body))
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return b, http.StatusInternalServerError, fmt.Errorf("error reading Response body: %s. Error: %v", string(body), err)
		}
		return body, resp.StatusCode, nil
	}
	return nil, http.StatusInternalServerError, fmt.Errorf("error returned from URL %s . Error message: %v", url, err)
}
