package iplantgroups

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// Subject represents a single subject returned by iplant-groups.
type Subject struct {
	Email       *string `json:"email"`
	Name        *string `json:"name"`
	LastName    *string `json:"last_name"`
	Description *string `json:"description"`
	ID          *string `json:"id"`
	Institution *string `json:"institution"`
	FirstName   *string `json:"first_name"`
	SourceID    *string `json:"source_id"`
}

// Client describes a single instance of this client library.
type Client struct {
	baseURL       string
	deGrouperUser string
}

// NewClient creates a new instance of this client library.
func NewClient(baseURL, deGrouperUser string) *Client {
	return &Client{
		baseURL:       baseURL,
		deGrouperUser: deGrouperUser,
	}
}

// buildURL builds the URL to use for the given path components.
func (c *Client) buildURL(pathComponents ...string) (string, error) {
	var err error

	// Parse the base URL.
	parsedBaseURL, err := url.Parse(c.baseURL)
	if err != nil {
		return "", err
	}

	// Build the relative path from the path components.
	relativePath := ""
	for _, pathComponent := range pathComponents {
		relativePath = relativePath + "/" + url.PathEscape(pathComponent)
	}

	// Append the relative path to the existing URL path.
	parsedBaseURL.Path = strings.TrimRight(parsedBaseURL.Path, "/") + relativePath

	// Add the user query argument.
	query := parsedBaseURL.Query()
	query.Set("user", c.deGrouperUser)
	parsedBaseURL.RawQuery = query.Encode()

	// Return the updated URL.
	return parsedBaseURL.String(), nil
}

// GetUserInfo looks up information for a single user by calling iplant-groups.
func (c *Client) GetUserInfo(username string) (*Subject, error) {
	errorMessage := fmt.Sprintf("unable to look up userinformation for %s", username)
	var err error

	// Build the request URL.
	requestURL, err := c.buildURL("subjects", username)
	if err != nil {
		return nil, errors.Wrap(err, errorMessage)
	}
	fmt.Printf("Request URL: %s\n", requestURL)

	// Submit the request.
	resp, err := http.Get(requestURL)
	if err != nil {
		return nil, errors.Wrap(err, errorMessage)
	}
	defer resp.Body.Close()

	// Check the HTTP status code.
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrap(err, errorMessage)
		}
		return nil, fmt.Errorf("%s: %s", errorMessage, respBody)
	}

	// Parse the response body.
	var subject Subject
	err = json.NewDecoder(resp.Body).Decode(&subject)
	if err != nil {
		return nil, errors.Wrap(err, errorMessage)
	}

	return &subject, nil
}
