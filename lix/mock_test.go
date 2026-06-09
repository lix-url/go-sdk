package lix

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

type capturedRequest struct {
	Method  string
	URL     string
	Headers http.Header
	Body    []byte
}

type mockResponse struct {
	statusCode int
	body       string
}

type mockHTTPClient struct {
	Requests  []capturedRequest
	responses []mockResponse
	inited    bool
}

func (m *mockHTTPClient) AddResponse(statusCode int, body string) {
	m.inited = true
	m.responses = append(m.responses, mockResponse{statusCode, body})
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	var body []byte
	if req.Body != nil {
		body, _ = io.ReadAll(req.Body)
	}
	m.Requests = append(m.Requests, capturedRequest{
		Method:  req.Method,
		URL:     req.URL.String(),
		Headers: req.Header,
		Body:    body,
	})
	if m.inited && len(m.responses) == 0 {
		return nil, fmt.Errorf("no responses in mock chain")
	}
	if m.inited {
		r := m.responses[0]
		m.responses = m.responses[1:]
		return &http.Response{
			StatusCode: r.statusCode,
			Body:       io.NopCloser(strings.NewReader(r.body)),
		}, nil
	}
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader("")),
	}, nil
}

func newTestClient(t *testing.T, apiKey string) (*Client, *mockHTTPClient) {
	t.Helper()
	mock := &mockHTTPClient{}
	client := NewClient(apiKey, mock)
	return client, mock
}

func strPtr(s string) *string { return &s }
func intPtr(i int) *int       { return &i }

func assertStr(t *testing.T, label, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %q, want %q", label, got, want)
	}
}

func assertInt(t *testing.T, label string, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %d, want %d", label, got, want)
	}
}

func assertBool(t *testing.T, label string, got, want bool) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %v, want %v", label, got, want)
	}
}

func assertNilString(t *testing.T, label string, got *string) {
	t.Helper()
	if got != nil {
		t.Errorf("%s: expected nil, got %q", label, *got)
	}
}

func assertNilInt(t *testing.T, label string, got *int) {
	t.Helper()
	if got != nil {
		t.Errorf("%s: expected nil, got %d", label, *got)
	}
}

func bodyField(t *testing.T, body map[string]interface{}, key string) interface{} {
	t.Helper()
	v, ok := body[key]
	if !ok {
		t.Errorf("body missing field %q", key)
	}
	return v
}

func assertBodyStr(t *testing.T, body map[string]interface{}, key, want string) {
	t.Helper()
	v := bodyField(t, body, key)
	if v != want {
		t.Errorf("body[%q]: got %v, want %q", key, v, want)
	}
}

func assertBodyFloat(t *testing.T, body map[string]interface{}, key string, want float64) {
	t.Helper()
	v := bodyField(t, body, key)
	if v != want {
		t.Errorf("body[%q]: got %v, want %v", key, v, want)
	}
}

func assertBodyBool(t *testing.T, body map[string]interface{}, key string, want bool) {
	t.Helper()
	v := bodyField(t, body, key)
	if v != want {
		t.Errorf("body[%q]: got %v, want %v", key, v, want)
	}
}

func assertBodyNil(t *testing.T, body map[string]interface{}, key string) {
	t.Helper()
	v := bodyField(t, body, key)
	if v != nil {
		t.Errorf("body[%q]: got %v, want nil", key, v)
	}
}

func assertOneRequest(t *testing.T, reqs []capturedRequest) capturedRequest {
	t.Helper()
	if len(reqs) != 1 {
		t.Fatalf("expected 1 request, got %d", len(reqs))
	}
	return reqs[0]
}
