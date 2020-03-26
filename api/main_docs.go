// Package api DE Administrative Requests API
//
// Documentation of the DE Administrative Requests API
//
//     Schemes: http
//     BasePath: /
//     Version: 1.0.0
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package api

import "github.com/cyverse-de/requests/model"

// swagger:route GET / misc getRoot
//
// General API Information
//
// Lists general information about the service API itself.
//
// responses:
//   200: rootResponse

// General information about the API.
// swagger:response rootResponse
type rootResponseWrapper struct {
	// in:body
	Body RootResponse
}

// Basic error response.
// swagger:response errorResponse
type errorResponseWrapper struct {
	// in:body
	Body ErrorResponse
}

// swagger:route GET /request-types request-types getRequestTypes
//
// List Request Types
//
// This endpoint lists all of the request types that have been registered, sorted by name.
//
// responses:
//    200: requestTypeListing

// Request type listing response.
// swagger:response requestTypeListing
type requestTypeListingWrapper struct {
	// in:body
	Body model.RequestTypeListing
}

// swagger:route POST /request-types/{name} request-types registerRequestType
//
// Register a Request Type
//
// This endpoint registers a new request type if a request type with the same name hasn't been registered already.
// If a request type with the same name has been registered already then the database is not modified and information
// about the existing request type is returned.
//
// responses:
//   200: requestType
//   400: errorResponse

// swagger:route GET /request-types/{name} request-types getRequestType
//
// Get a Request Type by Name
//
// This endpoint returns the request type with the given name if one exists.
//
// responses:
//   200: requestType
//   404: errorResponse

// Request type response.
// swagger:response requestType
type requestTypeWrapper struct {
	// in:body
	Body model.RequestType
}

// Parameters for registering a request type.
// swagger:parameters registerRequestType getRequestType
type registerRequestTypeParameters struct {
	// the name of the request type being registered
	//
	// in:path
	Name string
}

// swagger:route GET /request-status-codes request-status-codes getRequestStatusCodes
//
// List Request Status Codes
//
// This endpoint lists all of the request status codes that have been registered.
//
// responses:
//    200: requestStatusCodeListing

// Request status code listing response.
// swagger:response requestStatusCodeListing
type requestStatusCodeListingWrapper struct {
	// in:body
	Body model.RequestStatusCodeListing
}

// swagger:route POST /requests requests submitRequest
//
// Submit a Request
//
// This endpoint submits a new administrative request.
//
// Responses:
//   200: requestSummary

// Request summary information.
// swagger:response requestSummary
type requestSummaryWrapper struct {
	// in:body
	Body model.RequestSummary
}

// Parameters for the request submission endpoint.
// swagger:parameters submitRequest
type requestSubmissionParameters struct {
	// The request submission
	//
	// in:body
	Body model.RequestSubmission

	// The username of the authenticated user
	//
	// in:query
	// required:true
	User *string `json:"user"`
}

// swagger:route GET /requests requests listRequests
//
// List Requests
//
// This endpoint lists existing requests.
//
// Responses:
//   200: requestListing

// Request listing.
// swagger:response requestListing
type requestListingWrapper struct {
	// in:body
	Body model.RequestListing
}

// Parameters for the request listing enpdoint.
// swagger:parameters listRequests
type requestListingParameters struct {
	// Whether or not completed requests should be included in the listing
	//
	// in:query
	IncludeCompleted *bool `json:"include-completed"`
}

// swagger:route GET /requests/{id} requests getRequestInformation
//
// Get Request Information
//
// This endpoint returns information about the request with the given identifier.
//
// Responses:
//   200: requestDetails

// Request detail information.
// swagger:response requestDetails
type requestDetailsWrapper struct {
	// in:body
	Body model.RequestDetails
}

// Parameters for the request details endpoing.
// swagger:parameters getRequestInformation
type getRequestInformationParameters struct {
	// The request ID
	//
	// in:path
	ID *string
}
