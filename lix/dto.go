package lix

// UsageItem holds a single usage counter.
type UsageItem struct {
	Limit     *int `json:"limit"`
	Used      int  `json:"used"`
	Remaining *int `json:"remaining"`
}

// Usages holds usage counters for different quota types.
type Usages struct {
	Links     UsageItem `json:"links"`
	APILinks  UsageItem `json:"api_links"`
	MassLinks UsageItem `json:"mass_links"`
}

// ResponseMeta holds pagination metadata.
type ResponseMeta struct {
	Total   int     `json:"total"`
	Limit   int     `json:"limit"`
	NextURL *string `json:"next_url"`
}

// ClientInfo holds information about the authenticated API client.
type ClientInfo struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	CreatedDatetime string `json:"created_datetime"`
}

// User holds user account information.
type User struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	CreatedDatetime string `json:"created_datetime"`
}

// Plan holds subscription plan information.
type Plan struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	StartDatetime string `json:"start_datetime"`
	EndDatetime   string `json:"end_datetime"`
}

// Profile holds the full profile returned by /me.
type Profile struct {
	Client ClientInfo `json:"client"`
	User   User       `json:"user"`
	Plan   Plan       `json:"plan"`
	Usages Usages     `json:"usage"`
}

// Group holds link group information.
type Group struct {
	ID                  int     `json:"id"`
	Alias               string  `json:"alias"`
	URL                 string  `json:"url"`
	Name                string  `json:"name"`
	IsRotate            bool    `json:"is_rotate"`
	Description         *string `json:"description"`
	CreatedDatetime     string  `json:"created_datetime"`
	DeactivatedDatetime *string `json:"deactivated_datetime"`
}

// GroupsResult holds a paginated list of groups.
type GroupsResult struct {
	Groups []Group      `json:"data"`
	Meta   ResponseMeta `json:"meta"`
}

// Link holds short link information.
type Link struct {
	ID                   int         `json:"id"`
	Alias                string      `json:"alias"`
	URL                  string      `json:"url"`
	ShortURL             string      `json:"short_url"`
	Title                *string     `json:"title"`
	Group                *Group      `json:"group"`
	Tags                 []string    `json:"tags"`
	Meta                 interface{} `json:"meta"`
	IsPublic             bool        `json:"is_public"`
	CreatedDatetime      string      `json:"created_datetime"`
	ActiveBeforeDatetime *string     `json:"active_before_datetime"`
	DeletedDatetime      *string     `json:"deleted_datetime"`
}

// LinksResult holds a paginated list of links.
type LinksResult struct {
	Links []Link       `json:"data"`
	Meta  ResponseMeta `json:"meta"`
}

// LinkShortenResult holds the result of creating or updating a link.
type LinkShortenResult struct {
	Link  Link
	Usage UsageItem
}
