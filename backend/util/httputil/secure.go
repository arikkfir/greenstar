package httputil

import (
	"fmt"
	"net/url"
)

func IsSecure(rawURL string) (bool, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return false, fmt.Errorf("failed to parse URL: %w", err)
	}
	return parsed.Scheme == "https", nil
}
