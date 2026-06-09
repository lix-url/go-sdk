package lix

import (
	"encoding/json"
	"testing"
)

const groupJSON = `{"data":{"id":1503,"alias":"demo","url":"https://lix.li/g/demo","name":"Seller group","is_rotate":false,"description":"Marketing group","created_datetime":"2026-05-21T22:08:37+03:00","deactivated_datetime":null}}`
const groupJSONID10 = `{"data":{"id":10,"alias":"demo","url":"https://lix.li/g/demo","name":"Seller group","is_rotate":false,"description":"Marketing group","created_datetime":"2026-05-21T22:08:37+03:00","deactivated_datetime":null}}`
const groupsListJSON = `{"data":[{"id":1503,"alias":"demo","url":"https://lix.li/g/demo","name":"Seller group","is_rotate":false,"description":"Marketing group","created_datetime":"2026-05-21T22:08:37+03:00","deactivated_datetime":null}],"meta":{"total":1,"limit":20,"next_url":null}}`
const updateGroupJSON = `{"data":{"id":1503,"alias":"demo","url":"https://lix.li/g/demo","name":"Seller group","is_rotate":false,"description":"Updated description","created_datetime":"2026-05-21T22:08:37+03:00","deactivated_datetime":null}}`

func TestGroupsCreate(t *testing.T) {
	t.Run("maps DTO", func(t *testing.T) {
		client, mock := newTestClient(t, "lix_test_123")
		mock.AddResponse(200, groupJSON)

		group, err := client.Groups().Create("Seller group", strPtr("Marketing group"), false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		assertInt(t, "ID", group.ID, 1503)
		assertStr(t, "Name", group.Name, "Seller group")
		assertStr(t, "Alias", group.Alias, "demo")
		assertStr(t, "URL", group.URL, "https://lix.li/g/demo")
		assertBool(t, "IsRotate", group.IsRotate, false)
		if group.Description == nil || *group.Description != "Marketing group" {
			t.Errorf("Description: got %v, want %q", group.Description, "Marketing group")
		}
		assertStr(t, "CreatedDatetime", group.CreatedDatetime, "2026-05-21T22:08:37+03:00")
		assertNilString(t, "DeactivatedDatetime", group.DeactivatedDatetime)
	})

	t.Run("sends correct request", func(t *testing.T) {
		client, mock := newTestClient(t, "lix_test_some_key1")
		mock.AddResponse(200, groupJSON)

		_, err := client.Groups().Create("Seller group", strPtr("Marketing group"), true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		req := assertOneRequest(t, mock.Requests)
		assertStr(t, "URL", req.URL, "https://lix.li/api/1.0/groups")
		assertStr(t, "Method", req.Method, "POST")
		assertStr(t, "X-Api-Key", req.Headers.Get("X-Api-Key"), "lix_test_some_key1")
		assertStr(t, "User-Agent", req.Headers.Get("User-Agent"), "lix-go-sdk/0.1.0")

		var body map[string]interface{}
		if err := json.Unmarshal(req.Body, &body); err != nil {
			t.Fatalf("body parse error: %v", err)
		}
		assertBodyStr(t, body, "name", "Seller group")
		assertBodyStr(t, body, "description", "Marketing group")
		assertBodyBool(t, body, "is_rotate", true)
	})
}

func TestGroupsGet(t *testing.T) {
	t.Run("maps DTO", func(t *testing.T) {
		client, mock := newTestClient(t, "lix_test_some_key1")
		mock.AddResponse(200, groupJSONID10)

		group, err := client.Groups().Get(10)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		assertInt(t, "ID", group.ID, 10)
		assertStr(t, "Name", group.Name, "Seller group")
		assertStr(t, "Alias", group.Alias, "demo")
		assertStr(t, "URL", group.URL, "https://lix.li/g/demo")
		assertBool(t, "IsRotate", group.IsRotate, false)
		assertNilString(t, "DeactivatedDatetime", group.DeactivatedDatetime)
		assertStr(t, "CreatedDatetime", group.CreatedDatetime, "2026-05-21T22:08:37+03:00")
	})

	t.Run("sends correct request", func(t *testing.T) {
		client, mock := newTestClient(t, "lix_test_some_key1")
		mock.AddResponse(200, groupJSONID10)

		_, _ = client.Groups().Get(10)

		req := assertOneRequest(t, mock.Requests)
		assertStr(t, "URL", req.URL, "https://lix.li/api/1.0/groups/10")
		assertStr(t, "Method", req.Method, "GET")
		assertStr(t, "X-Api-Key", req.Headers.Get("X-Api-Key"), "lix_test_some_key1")
		assertStr(t, "User-Agent", req.Headers.Get("User-Agent"), "lix-go-sdk/0.1.0")
		if len(req.Body) != 0 {
			t.Errorf("expected empty body, got %q", req.Body)
		}
	})
}

func TestGroupsList(t *testing.T) {
	t.Run("maps DTO", func(t *testing.T) {
		client, mock := newTestClient(t, "lix_test_some_key1")
		mock.AddResponse(200, groupsListJSON)

		result, err := client.Groups().List(nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		assertInt(t, "len(Groups)", len(result.Groups), 1)
		assertInt(t, "Groups[0].ID", result.Groups[0].ID, 1503)
		assertStr(t, "Groups[0].Name", result.Groups[0].Name, "Seller group")
		assertStr(t, "Groups[0].URL", result.Groups[0].URL, "https://lix.li/g/demo")
		assertInt(t, "Meta.Total", result.Meta.Total, 1)
		assertInt(t, "Meta.Limit", result.Meta.Limit, 20)
		if result.Meta.NextURL != nil {
			t.Errorf("Meta.NextURL: expected nil, got %q", *result.Meta.NextURL)
		}
	})

	t.Run("no pagination params", func(t *testing.T) {
		client, mock := newTestClient(t, "lix_test_some_key1")
		mock.AddResponse(200, groupsListJSON)
		_, _ = client.Groups().List(nil, nil)
		req := assertOneRequest(t, mock.Requests)
		assertStr(t, "URL", req.URL, "https://lix.li/api/1.0/groups")
		assertStr(t, "Method", req.Method, "GET")
	})

	t.Run("with pagination params", func(t *testing.T) {
		client, mock := newTestClient(t, "lix_test_some_key1")
		mock.AddResponse(200, groupsListJSON)
		_, _ = client.Groups().List(intPtr(10), intPtr(1000))
		req := assertOneRequest(t, mock.Requests)
		assertStr(t, "URL", req.URL, "https://lix.li/api/1.0/groups?limit=10&from_id=1000")
	})
}

func TestGroupsUpdate(t *testing.T) {
	t.Run("maps DTO", func(t *testing.T) {
		client, mock := newTestClient(t, "lix_test_some_key1")
		mock.AddResponse(200, updateGroupJSON)

		group, err := client.Groups().Update(10, nil, strPtr("Updated description"), false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		assertInt(t, "ID", group.ID, 1503)
		assertStr(t, "Name", group.Name, "Seller group")
		if group.Description == nil || *group.Description != "Updated description" {
			t.Errorf("Description: got %v, want %q", group.Description, "Updated description")
		}
	})

	t.Run("sends correct request", func(t *testing.T) {
		client, mock := newTestClient(t, "lix_test_some_key1")
		mock.AddResponse(200, updateGroupJSON)

		_, _ = client.Groups().Update(10, nil, strPtr("Updated description"), false)

		req := assertOneRequest(t, mock.Requests)
		assertStr(t, "URL", req.URL, "https://lix.li/api/1.0/groups/10")
		assertStr(t, "Method", req.Method, "PATCH")
		assertStr(t, "X-Api-Key", req.Headers.Get("X-Api-Key"), "lix_test_some_key1")
		assertStr(t, "User-Agent", req.Headers.Get("User-Agent"), "lix-go-sdk/0.1.0")

		var body map[string]interface{}
		if err := json.Unmarshal(req.Body, &body); err != nil {
			t.Fatalf("body parse error: %v", err)
		}
		assertBodyNil(t, body, "name")
		assertBodyStr(t, body, "description", "Updated description")
		assertBodyBool(t, body, "is_rotate", false)
	})
}

func TestGroupsDelete(t *testing.T) {
	t.Run("sends correct request", func(t *testing.T) {
		client, mock := newTestClient(t, "lix_test_some_key1")
		mock.AddResponse(200, "{}")

		err := client.Groups().Delete(10)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		req := assertOneRequest(t, mock.Requests)
		assertStr(t, "URL", req.URL, "https://lix.li/api/1.0/groups/10")
		assertStr(t, "Method", req.Method, "DELETE")
		assertStr(t, "X-Api-Key", req.Headers.Get("X-Api-Key"), "lix_test_some_key1")
		assertStr(t, "User-Agent", req.Headers.Get("User-Agent"), "lix-go-sdk/0.1.0")
		if len(req.Body) != 0 {
			t.Errorf("expected empty body, got %q", req.Body)
		}
	})
}
