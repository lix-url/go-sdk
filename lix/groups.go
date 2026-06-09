package lix

import (
	"encoding/json"
	"fmt"
)

// GroupsResource provides methods for managing link groups.
type GroupsResource struct {
	api *apiClient
}

// Create creates a new link group.
func (r *GroupsResource) Create(name string, description *string, isRotate bool) (*Group, error) {
	data, err := r.api.post("groups", map[string]interface{}{
		"name":        name,
		"description": description,
		"is_rotate":   isRotate,
	})
	if err != nil {
		return nil, err
	}
	var resp struct {
		Data Group `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Update updates an existing group.
func (r *GroupsResource) Update(groupID int, name *string, description *string, isRotate bool) (*Group, error) {
	data, err := r.api.patch(fmt.Sprintf("groups/%d", groupID), map[string]interface{}{
		"name":        name,
		"description": description,
		"is_rotate":   isRotate,
	})
	if err != nil {
		return nil, err
	}
	var resp struct {
		Data Group `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Get retrieves a group by ID.
func (r *GroupsResource) Get(id int) (*Group, error) {
	data, err := r.api.get(fmt.Sprintf("groups/%d", id))
	if err != nil {
		return nil, err
	}
	var resp struct {
		Data Group `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Delete deletes a group by ID.
func (r *GroupsResource) Delete(id int) error {
	_, err := r.api.delete(fmt.Sprintf("groups/%d", id))
	return err
}

// List returns a paginated list of groups. Pass nil for limit/fromID to use API defaults.
func (r *GroupsResource) List(limit *int, fromID *int) (*GroupsResult, error) {
	endpoint := "groups"
	if q := buildQuery(limit, fromID); q != "" {
		endpoint += "?" + q
	}
	data, err := r.api.get(endpoint)
	if err != nil {
		return nil, err
	}
	var result GroupsResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
