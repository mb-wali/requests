package model

// RequestStatusCode describes a single request status code.
type RequestStatusCode struct {
	// The request status code identifier
	ID string `json:"id"`

	// The name of the request status code
	Name string `json:"name"`

	// The displayable name of the request status code
	DisplayName string `json:"display_name"`

	// The email template used for request status code
	EmailTemplate string `json:"email_template"`
}

// RequestStatusCodeListing describes a listing of existing request status codes.
type RequestStatusCodeListing struct {
	// The request status codes
	RequestStatusCodes []*RequestStatusCode `json:"request_status_codes"`
}
