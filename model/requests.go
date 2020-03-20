package model

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
