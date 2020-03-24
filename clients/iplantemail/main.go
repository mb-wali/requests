package iplantemail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// EmailRequestBody represents a request body sent to iplant-email.
type EmailRequestBody struct {
	To       string      `json:"to"`
	Template string      `json:"template"`
	Subject  string      `json:"subject"`
	Values   interface{} `json:"values"`
}

// Client describes a single instance of this client library.
type Client struct {
	baseURL string
}

// NewClient creates a new instance of this client library.
func NewClient(baseURL string) *Client {
	return &Client{baseURL: baseURL}
}

// SendRequestSubmittedEmail sends an email corresponding to a request.
func (c *Client) SendRequestSubmittedEmail(emailAddress, templateName string, requestDetails interface{}) error {
	errorMessage := "unable to send reqest notification email"
	var err error

	// Build the request body.
	body, err := json.Marshal(&EmailRequestBody{
		To:       emailAddress,
		Template: templateName,
		Subject:  "New Administrative Request",
		Values:   requestDetails,
	})

	// Submit the request.
	resp, err := http.Post(c.baseURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return errors.Wrap(err, errorMessage)
	}
	defer resp.Body.Close()

	// Check the HTTP Status code.
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, errorMessage)
		}
		return fmt.Errorf("%s: %s", errorMessage, respBody)
	}

	return nil
}
