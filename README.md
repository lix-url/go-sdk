# Lix.li Go SDK

Official Go SDK for the [Lix.li](https://lix.li) URL shortening and link analytics API.

## Requirements

- Go 1.22+
- No external dependencies (stdlib only)

## Installation

```bash
go get github.com/lix-url/go-sdk/lix
```

## Quick Start

```go
package main

import (
	"fmt"
	"log"
	"github.com/lix-url/go-sdk/lix"
)

func main() {
	client := lix.NewClient("lix_live_xxx")
	result, err := client.Links().Create(
		"https://example.com",
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.Link.ShortURL)
}
```

---

## Links

### Create a short link

```go
// Simple — defaults IsPublic: true
result, err := client.Links().Create("https://example.com", nil)
fmt.Println(result.Link.ShortURL)
fmt.Println(result.Usage.Remaining)
```

### Create with a custom alias

```go
opts := lix.NewCreateLinkOptions()
opts.Alias = strPtr("my-link")

result, err := client.Links().Create("https://example.com", opts)
fmt.Println(result.Link.ShortURL) // https://lix.li/my-link
```

### Create with all options

```go
opts := lix.NewCreateLinkOptions()
opts.Alias = strPtr("my-alias")
opts.Title = strPtr("My Page Title")
opts.GroupID = intPtr(42)
opts.Tags = []string{"sale", "promo"}
opts.Meta = map[string]string{"title": "Sale!", "description": "..."}
opts.UTM = map[string]string{"utm_source": "google", "utm_medium": "email"}
opts.TrackingPixelIDs = []int{1001, 1002}
opts.ActiveBeforeDatetime = strPtr("2029-12-31T23:59:59+00:00")
opts.Password = strPtr("secret123")
opts.IsPublic = true

result, err := client.Links().Create("https://example.com", opts)
```

### Update a link

```go
opts := lix.NewUpdateLinkOptions()
opts.Title = strPtr("New Title")

result, err := client.Links().Update(79697, opts)
fmt.Println(*result.Link.Title)
```

### Get a link by ID

```go
link, err := client.Links().Get(79697)
fmt.Println(link.ShortURL)
fmt.Println(link.URL)
if link.Group != nil {
    fmt.Println(link.Group.Name)
}
```

### List links

```go
// All links (API default pagination)
page, err := client.Links().List(nil, nil)

// With pagination
limit := 20
fromID := 79500
page, err := client.Links().List(&limit, &fromID)

for _, link := range page.Links {
    fmt.Println(link.ShortURL, link.URL)
}
fmt.Println(page.Meta.Total)
if page.Meta.NextURL != nil {
    fmt.Println(*page.Meta.NextURL)
}
```

### Delete a link

```go
err := client.Links().Delete(79697)
```

---

## Groups

### Create a group

```go
group, err := client.Groups().Create("Marketing", nil, false)
// With description and rotation:
group, err := client.Groups().Create("Landing Pages", strPtr("Rotating pages"), true)
```

### Update a group

```go
group, err := client.Groups().Update(10, nil, strPtr("Updated description"), false)
```

### Get a group by ID

```go
group, err := client.Groups().Get(10)
fmt.Println(group.Name)
fmt.Println(group.Alias)
fmt.Println(group.URL)
```

### List groups

```go
page, err := client.Groups().List(nil, nil)

limit := 10
fromID := 1000
page, err := client.Groups().List(&limit, &fromID)

for _, g := range page.Groups {
    fmt.Println(g.Name)
}
fmt.Println(page.Meta.Total)
```

### Delete a group

```go
err := client.Groups().Delete(10)
```

---

## Profile

```go
profile, err := client.Profile().Me()

fmt.Println(profile.Client.Name)
fmt.Println(profile.User.Email)
fmt.Println(profile.Plan.Name)
if profile.Usages.APILinks.Remaining != nil {
    fmt.Println(*profile.Usages.APILinks.Remaining)
}
```

---

## Error Handling

```go
import "errors"

result, err := client.Links().Create("https://example.com", nil)
if err != nil {
    var ve *lix.ValidationError
    var notFound *lix.NotFoundException
    var unauth *lix.UnauthorizedError
    var rateLimit *lix.RateLimitError
    var server *lix.ServerError

    switch {
    case errors.As(err, &ve):
        fmt.Println("Validation errors:", ve.Data)
    case errors.As(err, &notFound):
        fmt.Println("Not found")
    case errors.As(err, &unauth):
        fmt.Println("Invalid API key")
    case errors.As(err, &rateLimit):
        fmt.Println("Rate limit exceeded")
    case errors.As(err, &server):
        fmt.Println("Server error")
    default:
        log.Fatal(err)
    }
}
```

---

## Custom HTTP Client

Inject your own `HTTPClient` for timeouts, proxies, or testing:

```go
import "net/http"

httpClient := &http.Client{Timeout: 10 * time.Second}
client := lix.NewClient("your_api_key", httpClient)
```

---

## Running Tests

```bash
go test ./lix/
go test ./lix/ -v   # verbose
```

## Project Structure

```
lix/
  client.go       — main entry point (NewClient)
  dto.go          — all DTO structs (Group, Link, Profile, etc.)
  errors.go       — typed error types
  http_client.go  — HTTPClient interface + apiClient
  util.go         — query string builder
  groups.go       — GroupsResource
  links.go        — LinksResource + CreateLinkOptions / UpdateLinkOptions
  profile.go      — ProfileResource
example/
  main.go         — runnable usage example
go.mod
```

## Pointer Helpers

Since many parameters are optional (`*string`, `*int`), a simple helper in your code makes it convenient:

```go
func strPtr(s string) *string { return &s }
func intPtr(i int) *int       { return &i }
```


## Documentation

* API Documentation: https://lix.li/api
* OpenAPI Specification: https://github.com/lix-url/openapi

## Other SDKs

- PHP SDK: https://github.com/lix-url/php-sdk
- JavaScript SDK: https://github.com/lix-url/js-sdk
- Python SDK: https://github.com/lix-url/python-sdk

## Support

Need help with the API or SDK?

- Support Center: https://lix.li/support

## License

MIT