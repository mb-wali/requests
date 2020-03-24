package model

import "time"

// RequestSubmission represents a new request being submitted to this service.
type RequestSubmission struct {
	// The request type
	RequestType string `json:"request_type" validate:"required"`

	// Arbitrary JSON object describing the request details
	Details interface{} `json:"details" validate:"required"`
}

// RequestSummary represents a brief overview of an administrative request.
type RequestSummary struct {
	// The request ID
	ID string `json:"id"`

	// The username of the requesting user
	RequestingUser string `json:"requesting_user"`

	// The request type
	RequestType string `json:"request_type"`

	// Arbitrary JSON object describing the request details
	Details interface{} `json:"details"`
}

// RequestUpdate represents a request status update.
type RequestUpdate struct {
	// The request update ID
	ID string `json:"id"`

	// The request status code
	StatusCode string `json:"status"`

	// The username of the updating user
	UpdatingUser string `json:"updating_user"`

	// The timestamp corresponding to when the request was updated
	CreatedDate time.Time `json:"created_date"`

	// The message that was entered when the update was created
	Message string `json:"message"`
}

// RequestDetails represents the details of a request.
type RequestDetails struct {
	// The request ID
	ID string `json:"id"`

	// The username of the requesting user
	RequestingUser string `json:"requesting_user"`

	// The request type
	RequestType string `json:"request_type"`

	// Arbitrary JSON object describing the request details
	Details interface{} `json:"details"`

	// The status updates for this request.
	Updates []*RequestUpdate `json:"updates"`
}
