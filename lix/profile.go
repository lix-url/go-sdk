package lix

import "encoding/json"

// ProfileResource provides methods for retrieving profile information.
type ProfileResource struct {
	api *apiClient
}

// Me returns profile information for the authenticated client.
func (r *ProfileResource) Me() (*Profile, error) {
	data, err := r.api.get("me")
	if err != nil {
		return nil, err
	}
	var profile Profile
	if err := json.Unmarshal(data, &profile); err != nil {
		return nil, err
	}
	return &profile, nil
}
