package util

import (
	"net/url"
	"strings"
)

// BuildURL builds a complete URL from a base URL and a slice of path components.
func BuildURL(baseURL string, pathComponents []string) (*url.URL, error) {
	var err error

	// Parse the base URL.
	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	// Build the relative path from the path components.
	relativePath := ""
	for _, pathComponent := range pathComponents {
		relativePath = relativePath + "/" + url.PathEscape(pathComponent)
	}

	// Append the relative path to the existing URL path.
	parsedBaseURL.Path = strings.TrimRight(parsedBaseURL.Path, "/") + relativePath

	return parsedBaseURL, nil
}
