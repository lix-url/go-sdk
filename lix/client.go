package lix

import "net/http"

// Client is the main entry point for the Lix SDK.
type Client struct {
	profile *ProfileResource
	groups  *GroupsResource
	links   *LinksResource
}

// NewClient creates a new Lix API client with the given API key.
// An optional HTTPClient may be supplied as the second argument (useful for testing).
func NewClient(apiKey string, httpClient ...HTTPClient) *Client {
	var hc HTTPClient
	if len(httpClient) > 0 && httpClient[0] != nil {
		hc = httpClient[0]
	} else {
		hc = &http.Client{}
	}
	api := newAPIClient(apiKey, hc)
	return &Client{
		profile: &ProfileResource{api},
		groups:  &GroupsResource{api},
		links:   &LinksResource{api},
	}
}

// Profile returns the profile resource.
func (c *Client) Profile() *ProfileResource { return c.profile }

// Groups returns the groups resource.
func (c *Client) Groups() *GroupsResource { return c.groups }

// Links returns the links resource.
func (c *Client) Links() *LinksResource { return c.links }
