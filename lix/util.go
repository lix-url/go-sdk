package lix

import (
	"fmt"
	"strings"
)

func buildQuery(limit *int, fromID *int) string {
	var params []string
	if limit != nil {
		params = append(params, fmt.Sprintf("limit=%d", *limit))
	}
	if fromID != nil {
		params = append(params, fmt.Sprintf("from_id=%d", *fromID))
	}
	return strings.Join(params, "&")
}
