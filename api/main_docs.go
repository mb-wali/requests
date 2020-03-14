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
