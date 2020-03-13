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
// Returns general information about the API.
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
// Returns the list of registered request types.
// responses:
//    200: requestTypeListing

// Request type listing response.
// swagger:response requestTypeListing
type requestTypeListingWrapper struct {
	// in:body
	Body model.RequestTypeListing
}
