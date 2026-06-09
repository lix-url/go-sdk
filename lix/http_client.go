package lix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// HTTPClient is the interface for the HTTP transport layer. Implement it to
// inject a custom or mock HTTP client.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

const (
	defaultBaseURL   = "https://lix.li/api/1.0"
	defaultUserAgent = "lix-go-sdk/0.1.0"
)

type apiClient struct {
	apiKey     string
	httpClient HTTPClient
}

func newAPIClient(apiKey string, httpClient HTTPClient) *apiClient {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	return &apiClient{apiKey: apiKey, httpClient: httpClient}
}

func (c *apiClient) sendRequest(method, endpoint string, body interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", defaultBaseURL, endpoint)

	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Api-Key", c.apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", defaultUserAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &HttpClientError{LixError{err.Error()}}
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case 400:
		var errResp struct {
			ParameterErrors map[string]interface{} `json:"parameter_errors"`
		}
		_ = json.Unmarshal(data, &errResp)
		return nil, &ValidationError{LixError{"validation error"}, errResp.ParameterErrors}
	case 401:
		return nil, &UnauthorizedError{LixError{"unauthorized"}}
	case 404:
		return nil, &NotFoundException{LixError{"not found"}}
	case 429:
		return nil, &RateLimitError{LixError{"rate limit exceeded"}}
	case 500:
		return nil, &ServerError{LixError{"server error"}}
	}

	return data, nil
}

func (c *apiClient) get(endpoint string) ([]byte, error) {
	return c.sendRequest(http.MethodGet, endpoint, nil)
}

func (c *apiClient) delete(endpoint string) ([]byte, error) {
	return c.sendRequest(http.MethodDelete, endpoint, nil)
}

func (c *apiClient) patch(endpoint string, body interface{}) ([]byte, error) {
	return c.sendRequest(http.MethodPatch, endpoint, body)
}

func (c *apiClient) post(endpoint string, body interface{}) ([]byte, error) {
	return c.sendRequest(http.MethodPost, endpoint, body)
}
