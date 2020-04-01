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

// SendEmail sends an arbitrary email.
func (c *Client) sendEmail(requestBody *EmailRequestBody) error {
	errorMessage := "unable to send email"
	var err error

	// Serialize the request body.
	body, err := json.Marshal(requestBody)
	if err != nil {
		return errors.Wrap(err, errorMessage)
	}

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

// SendRequestSubmittedEmail sends an email corresponding to a request.
func (c *Client) SendRequestSubmittedEmail(emailAddress, templateName string, requestDetails interface{}) error {
	return c.sendEmail(&EmailRequestBody{
		To:       emailAddress,
		Template: templateName,
		Subject:  "New Administrative Request",
		Values:   requestDetails,
	})
}

// SendRequestUpdatedEmail sends an email corresponding to a request status update.
func (c *Client) SendRequestUpdatedEmail(emailAddress, templateName string, requestDetails interface{}) error {
	return c.sendEmail(&EmailRequestBody{
		To:       emailAddress,
		Template: templateName,
		Subject:  "Administrative Request Updated",
		Values:   requestDetails,
	})
}
