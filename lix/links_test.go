package lix

import (
	"encoding/json"
	"testing"
)

const linkWithGroupJSON = `{"data":{"id":79697,"alias":"demo","short_url":"https://lix.li/demo","url":"https://example.com/very/long/page","is_public":true,"title":"Some Title","created_datetime":"2026-05-27T22:16:22+03:00","active_before_datetime":null,"deleted_datetime":null,"group":{"id":1503,"alias":"demo","url":"https://lix.li/g/demo","name":"Seller group","is_rotate":false,"description":"Marketing group","created_datetime":"2026-05-21T22:08:37+03:00","deactivated_datetime":null},"tags":["sale","promo"],"meta":{"title":"Awesome sale!","og:title":"Awesome sale!!!!","description":"Woooooo, its wonderful!","og:description":"Woooowww, its wonderful!","keywords":"sale, promo"}},"usage":{"limit":500,"used":3,"remaining":497}}`
const linkNoGroupJSON = `{"data":{"id":79697,"alias":"demo2","short_url":"https://lix.li/demo2","url":"https://example.com","is_public":true,"title":null,"created_datetime":"2026-05-27T22:16:22+03:00","active_before_datetime":null,"deleted_datetime":null,"group":null,"tags":[],"meta":[]},"usage":{"limit":500,"used":3,"remaining":497}}`
const getLinkJSON = `{"data":{"id":79697,"alias":"demo2","short_url":"https://lix.li/demo2","url":"https://example.com","is_public":true,"title":null,"created_datetime":"2026-05-27T22:16:22+03:00","active_before_datetime":null,"deleted_datetime":null,"group":null,"tags":["sale","promo"],"meta":[]}}`
const linksListJSON = `{"data":[{"id":79618,"alias":"a2Ag4","short_url":"https://lix.li/a2Ag4","url":"https://example.com","is_public":true,"title":"Promo","created_datetime":"2026-05-21T23:35:02+03:00","active_before_datetime":"2029-05-21T21:25:40+03:00","deleted_datetime":null,"group":{"id":1005,"alias":"2222","url":"https://lix.li/g/2222","name":"Test","is_rotate":false,"description":"df","created_datetime":"2026-05-18T01:55:50+03:00","deactivated_datetime":null},"tags":["promo","sale"],"meta":{"title":"Promo"}},{"id":79615,"alias":"oSCZ0mP","short_url":"https://lix.li/oSCZ0mP","url":"https://console.cloud.google.com","is_public":true,"title":null,"created_datetime":"2026-05-18T03:31:12+03:00","active_before_datetime":null,"deleted_datetime":null,"group":null,"tags":[],"meta":{}}],"meta":{"total":5,"limit":2,"next_url":"https://lix.li/api/1.0/links?from_id=79615&limit=2"}}`
const updateLinkJSON = `{"data":{"id":79697,"alias":"demo2","short_url":"https://lix.li/demo2","url":"https://example.com","is_public":true,"title":"Updated title","created_datetime":"2026-05-27T22:16:22+03:00","active_before_datetime":null,"deleted_datetime":null,"group":null,"tags":[],"meta":[]},"usage":{"limit":500,"used":3,"remaining":497}}`

func TestLinksCreate(t *testing.T) {
	t.Run("maps DTO with group", func(t *testing.T) {
		client, mock := newTestClient(t, "lix_test_123")
		mock.AddResponse(201, linkWithGroupJSON)

		result, err := client.Links().Create("https://example.com/very/long/page", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		link := result.Link

		assertInt(t, "ID", link.ID, 79697)
		assertStr(t, "Alias", link.Alias, "demo")
		assertStr(t, "ShortURL", link.ShortURL, "https://lix.li/demo")
		assertStr(t, "URL", link.URL, "https://example.com/very/long/page")
		assertBool(t, "IsPublic", link.IsPublic, true)
		assertStr(t, "CreatedDatetime", link.CreatedDatetime, "2026-05-27T22:16:22+03:00")
		assertNilString(t, "ActiveBeforeDatetime", link.ActiveBeforeDatetime)
		assertNilString(t, "DeletedDatetime", link.DeletedDatetime)
		if link.Title == nil || *link.Title != "Some Title" {
			t.Errorf("Title: got %v, want %q", link.Title, "Some Title")
		}
		if len(link.Tags) != 2 || link.Tags[0] != "sale" || link.Tags[1] != "promo" {
			t.Errorf("Tags: got %v", link.Tags)
		}
		if link.Group == nil {
			t.Fatal("Group: expected non-nil")
		}
		assertInt(t, "Group.ID", link.Group.ID, 1503)
		assertStr(t, "Group.Alias", link.Group.Alias, "demo")
		assertStr(t, "Group.URL", link.Group.URL, "https://lix.li/g/demo")
		assertBool(t, "Group.IsRotate", link.Group.IsRotate, false)
		assertNilString(t, "Group.DeactivatedDatetime", link.Group.DeactivatedDatetime)
		assertStr(t, "Group.CreatedDatetime", link.Group.CreatedDatetime, "2026-05-21T22:08:37+03:00")

		assertInt(t, "Usage.Limit", *result.Usage.Limit, 500)
		assertInt(t, "Usage.Used", result.Usage.Used, 3)
		assertInt(t, "Usage.Remaining", *result.Usage.Remaining, 497)
	})

	t.Run("nil group and title when absent", func(t *testing.T) {
		client, mock := newTestClient(t, "lix_test_123")
		mock.AddResponse(201, linkNoGroupJSON)

		result, err := client.Links().Create("https://example.com", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Link.Group != nil {
			t.Errorf("Group: expected nil, got %v", result.Link.Group)
		}
		if result.Link.Title != nil {
			t.Errorf("Title: expected nil, got %q", *result.Link.Title)
		}
	})

	t.Run("sends correct request", func(t *testing.T) {
		client, mock := newTestClient(t, "lix_test_some_key1")
		mock.AddResponse(201, linkNoGroupJSON)

		opts := &CreateLinkOptions{
			Alias:                strPtr("demo"),
			Title:                strPtr("Some Title"),
			GroupID:              intPtr(1000),
			Tags:                 []string{"sale", "promo"},
			TrackingPixelIDs:     []int{1110, 1023},
			ActiveBeforeDatetime: strPtr("2029-05-21T21:25:40+03:00"),
			Password:             strPtr("12345"),
			IsPublic:             true,
		}
		_, err := client.Links().Create("https://example.com/very/long/page", opts)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		req := assertOneRequest(t, mock.Requests)
		assertStr(t, "URL", req.URL, "https://lix.li/api/1.0/links")
		assertStr(t, "Method", req.Method, "POST")
		assertStr(t, "X-Api-Key", req.Headers.Get("X-Api-Key"), "lix_test_some_key1")
		assertStr(t, "User-Agent", req.Headers.Get("User-Agent"), "lix-go-sdk/0.1.0")

		var body map[string]interface{}
		if err := json.Unmarshal(req.Body, &body); err != nil {
			t.Fatalf("body parse error: %v", err)
		}
		assertBodyStr(t, body, "url", "https://example.com/very/long/page")
		assertBodyStr(t, body, "alias", "demo")
		assertBodyStr(t, body, "title", "Some Title")
		assertBodyFloat(t, body, "group_id", 1000)
		assertBodyBool(t, body, "is_public", true)
		assertBodyStr(t, body, "active_before_datetime", "2029-05-21T21:25:40+03:00")
		assertBodyStr(t, body, "password", "12345")

		tags, _ := body["tags"].([]interface{})
		if len(tags) != 2 || tags[0] != "sale" || tags[1] != "promo" {
			t.Errorf("tags: got %v", tags)
		}
		pixels, _ := body["tracking_pixel_ids"].([]interface{})
		if len(pixels) != 2 || pixels[0] != float64(1110) || pixels[1] != float64(1023) {
			t.Errorf("tracking_pixel_ids: got %v", pixels)
		}
	})
}

func TestLinksGet(t *testing.T) {
	t.Run("maps DTO", func(t *testing.T) {
		client, mock := newTestClient(t, "lix_test_some_key1")
		mock.AddResponse(200, getLinkJSON)

		link, err := client.Links().Get(79697)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		assertInt(t, "ID", link.ID, 79697)
		assertStr(t, "Alias", link.Alias, "demo2")
		assertStr(t, "ShortURL", link.ShortURL, "https://lix.li/demo2")
		assertStr(t, "URL", link.URL, "https://example.com")
		assertBool(t, "IsPublic", link.IsPublic, true)
		if link.Group != nil {
			t.Errorf("Group: expected nil")
		}
		if len(link.Tags) != 2 || link.Tags[0] != "sale" || link.Tags[1] != "promo" {
			t.Errorf("Tags: got %v", link.Tags)
		}
	})

	t.Run("sends correct request", func(t *testing.T) {
		client, mock := newTestClient(t, "lix_test_some_key1")
		mock.AddResponse(200, getLinkJSON)

		_, _ = client.Links().Get(123)

		req := assertOneRequest(t, mock.Requests)
		assertStr(t, "URL", req.URL, "https://lix.li/api/1.0/links/123")
		assertStr(t, "Method", req.Method, "GET")
		assertStr(t, "X-Api-Key", req.Headers.Get("X-Api-Key"), "lix_test_some_key1")
		assertStr(t, "User-Agent", req.Headers.Get("User-Agent"), "lix-go-sdk/0.1.0")
		if len(req.Body) != 0 {
			t.Errorf("expected empty body")
		}
	})
}

func TestLinksList(t *testing.T) {
	t.Run("maps DTO", func(t *testing.T) {
		client, mock := newTestClient(t, "lix_test_some_key1")
		mock.AddResponse(200, linksListJSON)

		result, err := client.Links().List(nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		assertInt(t, "len(Links)", len(result.Links), 2)
		assertStr(t, "Links[0].ShortURL", result.Links[0].ShortURL, "https://lix.li/a2Ag4")
		assertStr(t, "Links[1].ShortURL", result.Links[1].ShortURL, "https://lix.li/oSCZ0mP")
		if result.Links[0].Group == nil {
			t.Error("Links[0].Group: expected non-nil")
		}
		if result.Links[1].Group != nil {
			t.Error("Links[1].Group: expected nil")
		}
		assertInt(t, "Meta.Total", result.Meta.Total, 5)
		assertInt(t, "Meta.Limit", result.Meta.Limit, 2)
		if result.Meta.NextURL == nil || *result.Meta.NextURL != "https://lix.li/api/1.0/links?from_id=79615&limit=2" {
			t.Errorf("Meta.NextURL: got %v", result.Meta.NextURL)
		}
	})

	t.Run("no pagination", func(t *testing.T) {
		client, mock := newTestClient(t, "lix_test_some_key1")
		mock.AddResponse(200, linksListJSON)
		_, _ = client.Links().List(nil, nil)
		req := assertOneRequest(t, mock.Requests)
		assertStr(t, "URL", req.URL, "https://lix.li/api/1.0/links")
		assertStr(t, "Method", req.Method, "GET")
	})

	t.Run("with pagination", func(t *testing.T) {
		client, mock := newTestClient(t, "lix_test_some_key1")
		mock.AddResponse(200, linksListJSON)
		_, _ = client.Links().List(intPtr(100), intPtr(500))
		req := assertOneRequest(t, mock.Requests)
		assertStr(t, "URL", req.URL, "https://lix.li/api/1.0/links?limit=100&from_id=500")
	})
}

func TestLinksUpdate(t *testing.T) {
	t.Run("maps DTO", func(t *testing.T) {
		client, mock := newTestClient(t, "lix_test_some_key1")
		mock.AddResponse(200, updateLinkJSON)

		opts := &UpdateLinkOptions{Title: strPtr("Updated title"), IsPublic: true}
		result, err := client.Links().Update(123, opts)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		assertInt(t, "ID", result.Link.ID, 79697)
		if result.Link.Title == nil || *result.Link.Title != "Updated title" {
			t.Errorf("Title: got %v", result.Link.Title)
		}
		assertStr(t, "ShortURL", result.Link.ShortURL, "https://lix.li/demo2")
		assertInt(t, "Usage.Limit", *result.Usage.Limit, 500)
		assertInt(t, "Usage.Used", result.Usage.Used, 3)
		assertInt(t, "Usage.Remaining", *result.Usage.Remaining, 497)
	})

	t.Run("sends correct request", func(t *testing.T) {
		client, mock := newTestClient(t, "lix_test_some_key1")
		mock.AddResponse(200, updateLinkJSON)

		opts := &UpdateLinkOptions{Title: strPtr("Updated title"), IsPublic: true}
		_, _ = client.Links().Update(123, opts)

		req := assertOneRequest(t, mock.Requests)
		assertStr(t, "URL", req.URL, "https://lix.li/api/1.0/links/123")
		assertStr(t, "Method", req.Method, "PATCH")
		assertStr(t, "X-Api-Key", req.Headers.Get("X-Api-Key"), "lix_test_some_key1")
		assertStr(t, "User-Agent", req.Headers.Get("User-Agent"), "lix-go-sdk/0.1.0")

		var body map[string]interface{}
		if err := json.Unmarshal(req.Body, &body); err != nil {
			t.Fatalf("body parse error: %v", err)
		}
		assertBodyStr(t, body, "title", "Updated title")
		assertBodyNil(t, body, "url")
		assertBodyNil(t, body, "group_id")
		assertBodyBool(t, body, "is_public", true)
	})
}

func TestLinksDelete(t *testing.T) {
	t.Run("sends correct request", func(t *testing.T) {
		client, mock := newTestClient(t, "lix_test_some_key1")
		mock.AddResponse(200, "{}")

		err := client.Links().Delete(123)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		req := assertOneRequest(t, mock.Requests)
		assertStr(t, "URL", req.URL, "https://lix.li/api/1.0/links/123")
		assertStr(t, "Method", req.Method, "DELETE")
		assertStr(t, "X-Api-Key", req.Headers.Get("X-Api-Key"), "lix_test_some_key1")
		assertStr(t, "User-Agent", req.Headers.Get("User-Agent"), "lix-go-sdk/0.1.0")
		if len(req.Body) != 0 {
			t.Errorf("expected empty body")
		}
	})
}
