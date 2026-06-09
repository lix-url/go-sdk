package lix

import (
	"testing"
)

const profileJSON = `{"client":{"id":1022,"name":"Test Client","email":"test@lix.li","created_datetime":"2022-04-24T17:38:42+03:00"},"user":{"name":"John Doe","email":"test_user@lix.li","created_datetime":"2023-04-14T17:38:42+03:00"},"plan":{"id":2,"name":"Pro","start_datetime":"2026-05-09T13:12:46+03:00","end_datetime":"2027-05-09T13:12:46+03:00"},"usage":{"links":{"limit":null,"used":1,"remaining":null},"api_links":{"limit":500,"used":100,"remaining":400},"mass_links":{"limit":100,"used":10,"remaining":90}}}`

func TestProfileMe(t *testing.T) {
	t.Run("maps DTO", func(t *testing.T) {
		client, mock := newTestClient(t, "lix_test_some_key1")
		mock.AddResponse(200, profileJSON)

		profile, err := client.Profile().Me()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		assertInt(t, "Client.ID", profile.Client.ID, 1022)
		assertStr(t, "Client.Name", profile.Client.Name, "Test Client")
		assertStr(t, "Client.Email", profile.Client.Email, "test@lix.li")
		assertStr(t, "Client.CreatedDatetime", profile.Client.CreatedDatetime, "2022-04-24T17:38:42+03:00")

		assertStr(t, "User.Name", profile.User.Name, "John Doe")
		assertStr(t, "User.Email", profile.User.Email, "test_user@lix.li")
		assertStr(t, "User.CreatedDatetime", profile.User.CreatedDatetime, "2023-04-14T17:38:42+03:00")

		assertInt(t, "Plan.ID", profile.Plan.ID, 2)
		assertStr(t, "Plan.Name", profile.Plan.Name, "Pro")
		assertStr(t, "Plan.StartDatetime", profile.Plan.StartDatetime, "2026-05-09T13:12:46+03:00")
		assertStr(t, "Plan.EndDatetime", profile.Plan.EndDatetime, "2027-05-09T13:12:46+03:00")

		assertNilInt(t, "Usages.Links.Limit", profile.Usages.Links.Limit)
		assertInt(t, "Usages.Links.Used", profile.Usages.Links.Used, 1)
		assertNilInt(t, "Usages.Links.Remaining", profile.Usages.Links.Remaining)

		assertInt(t, "Usages.APILinks.Limit", *profile.Usages.APILinks.Limit, 500)
		assertInt(t, "Usages.APILinks.Used", profile.Usages.APILinks.Used, 100)
		assertInt(t, "Usages.APILinks.Remaining", *profile.Usages.APILinks.Remaining, 400)

		assertInt(t, "Usages.MassLinks.Limit", *profile.Usages.MassLinks.Limit, 100)
		assertInt(t, "Usages.MassLinks.Used", profile.Usages.MassLinks.Used, 10)
		assertInt(t, "Usages.MassLinks.Remaining", *profile.Usages.MassLinks.Remaining, 90)
	})

	t.Run("sends correct request", func(t *testing.T) {
		client, mock := newTestClient(t, "lix_test_some_key1")
		mock.AddResponse(200, profileJSON)

		_, _ = client.Profile().Me()

		req := assertOneRequest(t, mock.Requests)
		assertStr(t, "URL", req.URL, "https://lix.li/api/1.0/me")
		assertStr(t, "Method", req.Method, "GET")
		assertStr(t, "X-Api-Key", req.Headers.Get("X-Api-Key"), "lix_test_some_key1")
		assertStr(t, "User-Agent", req.Headers.Get("User-Agent"), "lix-go-sdk/0.1.0")
		if len(req.Body) != 0 {
			t.Errorf("expected empty body")
		}
	})
}
