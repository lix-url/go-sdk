package lix

import (
	"encoding/json"
	"fmt"
)

// CreateLinkOptions holds optional parameters for link creation.
type CreateLinkOptions struct {
	Alias                *string
	Title                *string
	GroupID              *int
	Tags                 []string
	Meta                 interface{}
	UTM                  interface{}
	TrackingPixelIDs     []int
	ActiveBeforeDatetime *string
	Password             *string
	IsPublic             bool
}

// NewCreateLinkOptions returns CreateLinkOptions with sensible defaults (IsPublic: true).
func NewCreateLinkOptions() *CreateLinkOptions {
	return &CreateLinkOptions{IsPublic: true}
}

// UpdateLinkOptions holds optional parameters for link updates.
type UpdateLinkOptions struct {
	URL                  *string
	Title                *string
	GroupID              *int
	Tags                 []string
	Meta                 interface{}
	UTM                  interface{}
	TrackingPixelIDs     []int
	ActiveBeforeDatetime *string
	Password             *string
	IsPublic             bool
}

// NewUpdateLinkOptions returns UpdateLinkOptions with sensible defaults (IsPublic: true).
func NewUpdateLinkOptions() *UpdateLinkOptions {
	return &UpdateLinkOptions{IsPublic: true}
}

// LinksResource provides methods for managing short links.
type LinksResource struct {
	api *apiClient
}

// Create creates a new short link. opts may be nil (defaults to IsPublic: true).
func (r *LinksResource) Create(url string, opts *CreateLinkOptions) (*LinkShortenResult, error) {
	if opts == nil {
		opts = NewCreateLinkOptions()
	}
	data, err := r.api.post("links", buildCreateBody(url, opts))
	if err != nil {
		return nil, err
	}
	return parseLinkShortenResult(data)
}

// Update updates an existing short link.
func (r *LinksResource) Update(id int, opts *UpdateLinkOptions) (*LinkShortenResult, error) {
	if opts == nil {
		opts = NewUpdateLinkOptions()
	}
	data, err := r.api.patch(fmt.Sprintf("links/%d", id), buildUpdateBody(opts))
	if err != nil {
		return nil, err
	}
	return parseLinkShortenResult(data)
}

// Get retrieves a link by ID.
func (r *LinksResource) Get(id int) (*Link, error) {
	data, err := r.api.get(fmt.Sprintf("links/%d", id))
	if err != nil {
		return nil, err
	}
	var resp struct {
		Data Link `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Delete deletes a link by ID.
func (r *LinksResource) Delete(id int) error {
	_, err := r.api.delete(fmt.Sprintf("links/%d", id))
	return err
}

// List returns a paginated list of links. Pass nil for limit/fromID to use API defaults.
func (r *LinksResource) List(limit *int, fromID *int) (*LinksResult, error) {
	endpoint := "links"
	if q := buildQuery(limit, fromID); q != "" {
		endpoint += "?" + q
	}
	data, err := r.api.get(endpoint)
	if err != nil {
		return nil, err
	}
	var result LinksResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func buildCreateBody(url string, opts *CreateLinkOptions) map[string]interface{} {
	tags := opts.Tags
	if tags == nil {
		tags = []string{}
	}
	pixels := opts.TrackingPixelIDs
	if pixels == nil {
		pixels = []int{}
	}
	meta := opts.Meta
	if meta == nil {
		meta = []interface{}{}
	}
	utm := opts.UTM
	if utm == nil {
		utm = []interface{}{}
	}
	return map[string]interface{}{
		"group_id":               opts.GroupID,
		"url":                    url,
		"alias":                  opts.Alias,
		"password":               opts.Password,
		"title":                  opts.Title,
		"tags":                   tags,
		"is_public":              opts.IsPublic,
		"tracking_pixel_ids":     pixels,
		"meta":                   meta,
		"utm":                    utm,
		"active_before_datetime": opts.ActiveBeforeDatetime,
	}
}

func buildUpdateBody(opts *UpdateLinkOptions) map[string]interface{} {
	tags := opts.Tags
	if tags == nil {
		tags = []string{}
	}
	pixels := opts.TrackingPixelIDs
	if pixels == nil {
		pixels = []int{}
	}
	meta := opts.Meta
	if meta == nil {
		meta = []interface{}{}
	}
	utm := opts.UTM
	if utm == nil {
		utm = []interface{}{}
	}
	return map[string]interface{}{
		"group_id":               opts.GroupID,
		"url":                    opts.URL,
		"password":               opts.Password,
		"title":                  opts.Title,
		"tags":                   tags,
		"is_public":              opts.IsPublic,
		"tracking_pixel_ids":     pixels,
		"meta":                   meta,
		"utm":                    utm,
		"active_before_datetime": opts.ActiveBeforeDatetime,
	}
}

func parseLinkShortenResult(data []byte) (*LinkShortenResult, error) {
	var resp struct {
		Data  Link      `json:"data"`
		Usage UsageItem `json:"usage"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return &LinkShortenResult{Link: resp.Data, Usage: resp.Usage}, nil
}
