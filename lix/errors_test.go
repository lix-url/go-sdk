package lix

import (
	"errors"
	"testing"
)

const validationErrJSON = `{"error":"invalid_parameters","parameter_errors":{"name":{"code":"required","message":"field required"}},"error_message":null}`

func TestValidationError(t *testing.T) {
	client, mock := newTestClient(t, "lix_test_some_key1")
	mock.AddResponse(400, validationErrJSON)

	_, err := client.Groups().Create("test", nil, false)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	nameErr, ok := ve.Data["name"].(map[string]interface{})
	if !ok {
		t.Fatalf("Data[\"name\"]: expected map, got %T", ve.Data["name"])
	}
	if nameErr["code"] != "required" {
		t.Errorf("code: got %v, want %q", nameErr["code"], "required")
	}
}

type errorCase struct {
	name       string
	statusCode int
	errType    interface{}
	call       func(*Client) error
}

func TestAPIErrors(t *testing.T) {
	cases := []struct {
		label      string
		statusCode int
		target     interface{ Error() string }
		call       func(*Client) error
	}{
		// Groups
		{"groups.Get 401", 401, &UnauthorizedError{}, func(c *Client) error { _, e := c.Groups().Get(1); return e }},
		{"groups.Get 404", 404, &NotFoundException{}, func(c *Client) error { _, e := c.Groups().Get(1); return e }},
		{"groups.Get 429", 429, &RateLimitError{}, func(c *Client) error { _, e := c.Groups().Get(1); return e }},
		{"groups.Get 500", 500, &ServerError{}, func(c *Client) error { _, e := c.Groups().Get(1); return e }},

		{"groups.List 401", 401, &UnauthorizedError{}, func(c *Client) error { _, e := c.Groups().List(nil, nil); return e }},
		{"groups.List 404", 404, &NotFoundException{}, func(c *Client) error { _, e := c.Groups().List(nil, nil); return e }},
		{"groups.List 429", 429, &RateLimitError{}, func(c *Client) error { _, e := c.Groups().List(nil, nil); return e }},
		{"groups.List 500", 500, &ServerError{}, func(c *Client) error { _, e := c.Groups().List(nil, nil); return e }},

		{"groups.Delete 401", 401, &UnauthorizedError{}, func(c *Client) error { return c.Groups().Delete(1) }},
		{"groups.Delete 404", 404, &NotFoundException{}, func(c *Client) error { return c.Groups().Delete(1) }},
		{"groups.Delete 429", 429, &RateLimitError{}, func(c *Client) error { return c.Groups().Delete(1) }},
		{"groups.Delete 500", 500, &ServerError{}, func(c *Client) error { return c.Groups().Delete(1) }},

		{"groups.Create 401", 401, &UnauthorizedError{}, func(c *Client) error { _, e := c.Groups().Create("t", nil, false); return e }},
		{"groups.Create 404", 404, &NotFoundException{}, func(c *Client) error { _, e := c.Groups().Create("t", nil, false); return e }},
		{"groups.Create 429", 429, &RateLimitError{}, func(c *Client) error { _, e := c.Groups().Create("t", nil, false); return e }},
		{"groups.Create 500", 500, &ServerError{}, func(c *Client) error { _, e := c.Groups().Create("t", nil, false); return e }},

		{"groups.Update 401", 401, &UnauthorizedError{}, func(c *Client) error { _, e := c.Groups().Update(1, nil, nil, false); return e }},
		{"groups.Update 404", 404, &NotFoundException{}, func(c *Client) error { _, e := c.Groups().Update(1, nil, nil, false); return e }},
		{"groups.Update 429", 429, &RateLimitError{}, func(c *Client) error { _, e := c.Groups().Update(1, nil, nil, false); return e }},
		{"groups.Update 500", 500, &ServerError{}, func(c *Client) error { _, e := c.Groups().Update(1, nil, nil, false); return e }},

		// Links
		{"links.Get 401", 401, &UnauthorizedError{}, func(c *Client) error { _, e := c.Links().Get(1); return e }},
		{"links.Get 404", 404, &NotFoundException{}, func(c *Client) error { _, e := c.Links().Get(1); return e }},
		{"links.Get 429", 429, &RateLimitError{}, func(c *Client) error { _, e := c.Links().Get(1); return e }},
		{"links.Get 500", 500, &ServerError{}, func(c *Client) error { _, e := c.Links().Get(1); return e }},

		{"links.List 401", 401, &UnauthorizedError{}, func(c *Client) error { _, e := c.Links().List(nil, nil); return e }},
		{"links.List 404", 404, &NotFoundException{}, func(c *Client) error { _, e := c.Links().List(nil, nil); return e }},
		{"links.List 429", 429, &RateLimitError{}, func(c *Client) error { _, e := c.Links().List(nil, nil); return e }},
		{"links.List 500", 500, &ServerError{}, func(c *Client) error { _, e := c.Links().List(nil, nil); return e }},

		{"links.Delete 401", 401, &UnauthorizedError{}, func(c *Client) error { return c.Links().Delete(1) }},
		{"links.Delete 404", 404, &NotFoundException{}, func(c *Client) error { return c.Links().Delete(1) }},
		{"links.Delete 429", 429, &RateLimitError{}, func(c *Client) error { return c.Links().Delete(1) }},
		{"links.Delete 500", 500, &ServerError{}, func(c *Client) error { return c.Links().Delete(1) }},

		{"links.Create 401", 401, &UnauthorizedError{}, func(c *Client) error { _, e := c.Links().Create("https://x.com", nil); return e }},
		{"links.Create 404", 404, &NotFoundException{}, func(c *Client) error { _, e := c.Links().Create("https://x.com", nil); return e }},
		{"links.Create 429", 429, &RateLimitError{}, func(c *Client) error { _, e := c.Links().Create("https://x.com", nil); return e }},
		{"links.Create 500", 500, &ServerError{}, func(c *Client) error { _, e := c.Links().Create("https://x.com", nil); return e }},

		{"links.Update 401", 401, &UnauthorizedError{}, func(c *Client) error { _, e := c.Links().Update(1, nil); return e }},
		{"links.Update 404", 404, &NotFoundException{}, func(c *Client) error { _, e := c.Links().Update(1, nil); return e }},
		{"links.Update 429", 429, &RateLimitError{}, func(c *Client) error { _, e := c.Links().Update(1, nil); return e }},
		{"links.Update 500", 500, &ServerError{}, func(c *Client) error { _, e := c.Links().Update(1, nil); return e }},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.label, func(t *testing.T) {
			client, mock := newTestClient(t, "lix_test_some_key1")
			mock.AddResponse(tc.statusCode, "")
			err := tc.call(client)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !errors.As(err, &tc.target) {
				t.Errorf("expected %T, got %T: %v", tc.target, err, err)
			}
		})
	}
}
