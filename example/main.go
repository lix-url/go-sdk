package main

import (
	"fmt"
	"log"

	"github.com/lix-url/go-sdk/lix"
)

func main() {
	// Insert your API key here
	client := lix.NewClient("YOUR_API_KEY_HERE")

	// --- Profile ------------------------------------------------------------------
	fmt.Println("=== Profile ===")
	profile, err := client.Profile().Me()
	if err != nil {
		log.Fatalf("profile error: %v", err)
	}
	fmt.Printf("Client: %s | %s\n", profile.Client.Name, profile.Client.Email)
	fmt.Printf("User:   %s | %s\n", profile.User.Name, profile.User.Email)
	fmt.Printf("Plan:   %s (until %s)\n", profile.Plan.Name, profile.Plan.EndDatetime)
	if profile.Usages.APILinks.Remaining != nil {
		fmt.Printf("API links remaining: %d / %d\n", *profile.Usages.APILinks.Remaining, *profile.Usages.APILinks.Limit)
	}
	fmt.Println()

	// --- Create a short link -------------------------------------------------------
	fmt.Println("=== Create link ===")
	opts := lix.NewCreateLinkOptions()
	opts.Title = strPtr("My test link")

	created, err := client.Links().Create("https://example.com", opts)
	if err != nil {
		log.Fatalf("create link error: %v", err)
	}
	fmt.Printf("Created: %s → %s\n", created.Link.ShortURL, created.Link.URL)
	fmt.Printf("ID: %d | Alias: %s\n", created.Link.ID, created.Link.Alias)
	if created.Usage.Limit != nil {
		fmt.Printf("API links used: %d / %d\n", created.Usage.Used, *created.Usage.Limit)
	}
	fmt.Println()

	linkID := created.Link.ID

	// --- Get link by ID -----------------------------------------------------------
	fmt.Println("=== Get link by ID ===")
	link, err := client.Links().Get(linkID)
	if err != nil {
		log.Fatalf("get link error: %v", err)
	}
	title := "(none)"
	if link.Title != nil {
		title = *link.Title
	}
	fmt.Printf("Link: %s | Title: %s\n", link.ShortURL, title)
	fmt.Println()

	// --- Update link --------------------------------------------------------------
	fmt.Println("=== Update link ===")
	updateOpts := lix.NewUpdateLinkOptions()
	updateOpts.Title = strPtr("Updated title")

	updated, err := client.Links().Update(linkID, updateOpts)
	if err != nil {
		log.Fatalf("update link error: %v", err)
	}
	updatedTitle := "(none)"
	if updated.Link.Title != nil {
		updatedTitle = *updated.Link.Title
	}
	fmt.Printf("Updated: %s | Title: %s\n", updated.Link.ShortURL, updatedTitle)
	fmt.Println()

	// --- List links ---------------------------------------------------------------
	fmt.Println("=== List links (first 3) ===")
	limit := 3
	linksPage, err := client.Links().List(&limit, nil)
	if err != nil {
		log.Fatalf("list links error: %v", err)
	}
	for _, l := range linksPage.Links {
		fmt.Printf("  - %s → %s\n", l.ShortURL, l.URL)
	}
	fmt.Printf("Total links: %d\n", linksPage.Meta.Total)
	fmt.Println()

	// --- Create a group -----------------------------------------------------------
	fmt.Println("=== Create group ===")
	group, err := client.Groups().Create("Test group", strPtr("Created via Go SDK"), false)
	if err != nil {
		log.Fatalf("create group error: %v", err)
	}
	fmt.Printf("Group created: %s | ID: %d | URL: %s\n", group.Name, group.ID, group.URL)
	fmt.Println()

	// --- List groups --------------------------------------------------------------
	fmt.Println("=== List groups ===")
	groupLimit := 5
	groupsPage, err := client.Groups().List(&groupLimit, nil)
	if err != nil {
		log.Fatalf("list groups error: %v", err)
	}
	for _, g := range groupsPage.Groups {
		fmt.Printf("  - %s | %s\n", g.Name, g.URL)
	}
	fmt.Printf("Total groups: %d\n", groupsPage.Meta.Total)
	fmt.Println()

	// --- Delete test link ---------------------------------------------------------
	fmt.Println("=== Delete test link ===")
	if err := client.Links().Delete(linkID); err != nil {
		log.Fatalf("delete link error: %v", err)
	}
	fmt.Printf("Link %d deleted.\n", linkID)
}

func strPtr(s string) *string { return &s }
